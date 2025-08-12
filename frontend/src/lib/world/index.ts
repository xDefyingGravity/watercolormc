/**
 * index.ts
 *
 * Description:
 *   frontend — index.ts module.
 *
 * Created: 7/9/25
 * Author: Will Ballantine
 *
 * @packageDocumentation
 * @copyright 2025-present Will Ballantine
 */
import { safeFetch } from '$lib/utils/fetch'
import { baseUrl } from '$lib/config'
import { debug, error } from '$lib/logging'

export interface World {
	name: string
	path: string
	seed: string
	type: string
}

export async function getWorld(serverId: string): Promise<World> {
	const response = await safeFetch<World | string>(baseUrl + '/api/servers/' + serverId + '/world')
	if (typeof response === 'string') {
		error('Failed to fetch world data: ' + response)
		throw new Error('Failed to fetch world data: ' + response)
	}
	if (!response) {
		error('No world data found for server ' + serverId)
		throw new Error('No world data found for server ' + serverId)
	}
	return response
}

export async function uploadWorld(serverId: string, file: File): Promise<void> {
	const formData = new FormData()
	formData.append('file', file)

	const response = await safeFetch<string>(baseUrl + '/api/servers/' + serverId + '/world/upload', {
		method: 'POST',
		body: formData
	})

	if (response != 'ok') {
		error('Failed to upload world: ' + response)
		throw new Error('Failed to upload world: ' + response)
	} else {
		debug('World uploaded successfully for server ' + serverId)
	}
}
