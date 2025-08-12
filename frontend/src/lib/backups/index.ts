/**
 * index.ts
 *
 * Description:
 *   frontend — index.ts module.
 *
 * Created: 7/11/25
 * Author: Will Ballantine
 *
 * @packageDocumentation
 * @copyright 2025-present Will Ballantine
 */
import { safeFetch } from '$lib/utils/fetch'
import { baseUrl } from '$lib/config'

export async function getBackups(serverId: string): Promise<string[]> {
	const response = await safeFetch<string[]>(`${baseUrl}/api/servers/${serverId}/backups`)

	if (response === null) return []
	if (response === undefined || !Array.isArray(response)) {
		throw new Error(`Error fetching backups for server: ${serverId}. Response: ${response}`)
	}

	return response
}

export async function backupServer(serverId: string): Promise<void> {
	const response = await safeFetch<void>(`${baseUrl}/api/servers/${serverId}/backup`, {
		method: 'POST'
	})

	if (response === undefined || response !== 'ok') {
		throw new Error(`Error creating backup for server: ${serverId}. Response: ${response}`)
	}

	return response
}

export async function restoreServerBackup(serverId: string, backupName: string): Promise<void> {
	const response = await safeFetch<void>(`${baseUrl}/api/servers/${serverId}/restore`, {
		method: 'POST',
		body: JSON.stringify({ backupName }),
		headers: {
			'Content-Type': 'application/json'
		}
	})

	if (response === undefined || response !== 'ok') {
		throw new Error(`Error creating backup for server: ${serverId}. Response: ${response}`)
	}

	return response
}

export async function deleteServerBackup(serverId: string, backupName: string): Promise<void> {
	const response = await safeFetch<void>(`${baseUrl}/api/servers/${serverId}/backups`, {
		method: 'DELETE',
		body: JSON.stringify({ backupName }),
		headers: {
			'Content-Type': 'application/json'
		}
	})

	if (response === undefined || response !== 'ok') {
		throw new Error(`Error deleting backup for server: ${serverId}. Response: ${response}`)
	}

	return response
}
