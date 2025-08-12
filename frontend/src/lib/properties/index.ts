/**
 * index.ts
 *
 * Description:
 *   frontend — index.ts module.
 *
 * Created: 7/10/25
 * Author: Will Ballantine
 *
 * @packageDocumentation
 * @copyright 2025-present Will Ballantine
 */
import { baseUrl } from '$lib/config'
import { safeFetch } from '$lib/utils/fetch'
import { error } from '$lib/logging'
import { servers } from '$lib/stores'

export async function getServerProperties(serverId: string): Promise<Record<string, string>> {
	const response = await safeFetch<Record<string, string> | string>(
		`${baseUrl}/api/servers/${serverId}/properties`
	)
	if (typeof response === 'string' || response === undefined) {
		error('Failed to fetch server properties: ' + response)
		throw new Error(`Failed to fetch server properties: ${response}`)
	}
	if (response === undefined) {
		error('No properties found for server ' + serverId)
		return {}
	}
	return response
}

export async function updateServerProperties(
	serverId: string,
	props: Record<string, string>
): Promise<boolean> {
	const res = await safeFetch<string>(`${baseUrl}/api/servers/${serverId}/properties`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(props)
	})

	if ((typeof res === 'string' && res !== 'OK') || res === undefined) {
		error('Failed to update server properties: ' + res)
		throw new Error(`Failed to update server properties: ${res}`)
	}

	for (const [key, value] of Object.entries(props)) {
		if (key === 'server-port') {
			servers.update((list) =>
				list.map((s) => (s.id === serverId ? { ...s, port: parseInt(value, 10) } : s))
			)
		} else if (key === 'server-ip') {
			servers.update((list) => list.map((s) => (s.id === serverId ? { ...s, host: value } : s)))
		}
	}

	return res !== undefined
}
