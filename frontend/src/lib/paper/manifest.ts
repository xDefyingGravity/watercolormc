/**
 * manifest.ts
 *
 * Description:
 *   frontend — manifest.ts module.
 *
 * Created: 7/12/25
 * Author: Will Ballantine
 *
 * @packageDocumentation
 * @copyright 2025-present Will Ballantine
 */
import { safeFetch } from '$lib/utils/fetch'
import { baseUrl } from '$lib/config'

export interface ManifestPlugin {
	id: string
	jar_name: string
}

export interface Manifest {
	plugins: ManifestPlugin[]
}

export async function getManifest(serverId: string): Promise<Manifest> {
	const response = await safeFetch<ManifestPlugin[] | string>(baseUrl + `/api/servers/${serverId}/plugins/manifest`)

	if (typeof response === 'string' || response === undefined) {
		throw new Error(response)
	}

	return {
		plugins: response
	}
}

export async function addPluginToManifest(serverId: string, plugin: ManifestPlugin): Promise<void> {
	const response = await safeFetch<string>(baseUrl + `/api/servers/${serverId}/plugins/manifest`, {
		method: 'POST',
		body: JSON.stringify(plugin),
		headers: {
			'Content-Type': 'application/json'
		}
	})

	if (response !== 'ok') {
		throw new Error(typeof response === 'string' ? response : 'unexpected response')
	}
}

export async function removePluginFromManifest(serverId: string, pluginIdOrJarName: string): Promise<void> {
	const encoded = encodeURIComponent(pluginIdOrJarName)
	const response = await safeFetch<string>(
		`${baseUrl}/api/servers/${serverId}/plugins/manifest/${encoded}`,
		{
			method: 'DELETE'
		}
	)

	if (response !== 'ok') {
		throw new Error(typeof response === 'string' ? response : 'unexpected response')
	}
}