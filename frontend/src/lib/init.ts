import * as logger from '$lib/logging'
import { safeFetch } from '$lib/utils/fetch'
import { apiRoutes } from '$lib/config'
import type { Server } from '$lib/types/server'
import { backendStatus, ipinfo, servers as serversStore } from '$lib/stores'
import { extractPaperVersion } from '$lib/paper/version'

export async function init() {
	logger.init()

	logger.info('Application initialization started')

	if (import.meta.env.DEV) {
		logger.info('Development mode is enabled')
	} else {
		logger.info('Production mode is enabled')
	}

	const result = await safeFetch<string>(apiRoutes.upStatusCheck.path)
	if (result === 'ok') {
		logger.info('API status: OK')
	} else {
		logger.error('API status check failed', { result })
		backendStatus.set('error')
	}

	logger.info('Loading server list from API')

	const servers = await safeFetch<Server[]>(apiRoutes.getAllServers.path)
	if (servers !== null && Array.isArray(servers)) {
		logger.info(`Loaded ${servers.length} servers from API`)
	} else {
		logger.error('Failed to load server list from API')
	}

	for (const server of servers ?? []) {
		console.log(server.status)
		server.status ??= 'offline'
		if (server.version.includes('papermc')) {
			server.version = extractPaperVersion(server.version)
			server.type = 'papermc'
		} else {
			server.type = 'vanilla'
		}
	}

	serversStore.set((servers as Server[]) || [])

	const ipinfoFromServer = await safeFetch<{ privateIp: string; publicIp: string }>(
		apiRoutes.getIpInfo.path
	)
	if (ipinfoFromServer && typeof ipinfoFromServer === 'object') {
		logger.info(`IP Info: ${ipinfoFromServer.privateIp}, ${ipinfoFromServer.publicIp}`)
		ipinfo.set(ipinfoFromServer)
	} else {
		logger.error('Failed to fetch IP information')
	}

	logger.info('Application initialization completed')
}
