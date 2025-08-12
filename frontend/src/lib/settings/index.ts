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
import { safeFetch } from '$lib/utils/fetch'
import { baseUrl } from '$lib/config'

export interface ServerSettings {
	Versions: {
		WatercolorVersion: string
		MinecraftVersion: string
	}
	JavaSettings: {
		Memory: {
			Min: number
			Max: number
		}
		JavaPath: string
		JvmArgs: string[]
	}
}

export interface GloblSettings {
	BasePath: string
}

export async function getServerSettings(serverId: string): Promise<ServerSettings> {
	const response = await safeFetch<ServerSettings | string>(
		baseUrl + `/api/servers/${serverId}/config`,
		{
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		}
	)

	if (typeof response === 'string' || response === undefined) {
		throw new Error(`Failed to fetch server settings: ${response}`)
	}

	return response as ServerSettings
}

export async function updateServerSettings(
	serverId: string,
	settings: ServerSettings
): Promise<boolean> {
	const response = await safeFetch<string>(baseUrl + `/api/servers/${serverId}/config`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(settings)
	})

	if ((typeof response === 'string' && response !== 'OK') || response === undefined) {
		throw new Error(`Failed to update server settings: ${response}`)
	}

	return true
}

export async function getGlobalSettings(): Promise<GloblSettings> {
	const response = await safeFetch<GloblSettings | string>(baseUrl + '/api/settings', {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		}
	})

	if (typeof response === 'string' || response === undefined) {
		throw new Error(`Failed to fetch global settings: ${response}`)
	}

	return response as GloblSettings
}

export async function updateGlobalSettings(settings: GloblSettings): Promise<boolean> {
	if (settings.BasePath.endsWith('/')) settings.BasePath = settings.BasePath.slice(0, -1)

	const response = await safeFetch<string>(baseUrl + '/api/settings', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(settings)
	})

	if ((typeof response === 'string' && response !== 'ok') || response === undefined) {
		throw new Error(`Failed to update global settings: ${response}`)
	}

	return true
}