import type { Server } from '$lib/types/server'
import { servers } from '$lib/stores'
import { get } from 'svelte/store'
import { safeFetch } from '$lib/utils/fetch'
import { apiRoutes } from '$lib/config'

function createIdFromName(name: string): string {
	return name.toLowerCase().replace(/\s+/g, '-')
}

function randomPort(): number {
	const min = 3000
	const max = 65535
	return Math.floor(Math.random() * (max - min + 1)) + min
}

export interface CreateServerResult {
	failed: boolean
	message: string
}

export async function createServer(
	name: string,
	description: string,
	version: string
): Promise<CreateServerResult> {
	const id = createIdFromName(name)

	if (get(servers).find((server) => server.name === name || server.id === id)) {
		return {
			failed: true,
			message: `Server with name "${name}" or ID "${id}" already exists.`
		}
	}

	let port = randomPort()
	while (get(servers).some((server) => server.port === port)) {
		port = randomPort()
	}

	const server = {
		id,
		name,
		description,
		version,
		port,
		host: '0.0.0.0',
		createdAt: new Date().toISOString()
	}

	const response = await safeFetch<string | Server>(apiRoutes.createServer.path, {
		method: apiRoutes.createServer.method,
		body: JSON.stringify(server),
		headers: {
			'Content-Type': 'application/json'
		}
	})

	if (response && typeof response === 'object' && 'id' in response) {
		servers.update((current) => [
			...current,
			{
				...server,
				status: 'offline'
			}
		])
		return {
			failed: false,
			message: `Server "${name}" created successfully.`
		}
	} else {
		return {
			failed: true,
			message: `Failed to create server: ${response ?? 'Unknown error'}`
		}
	}
}
