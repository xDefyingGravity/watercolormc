import { baseUrl } from '$lib/config'
import { toWebSocketUrl } from '$lib/utils/websocket'
import { debug } from '$lib/logging'

export interface ServerLog {
	message: string
	type: 'stdout' | 'stderr'
}

export interface ServerStatsMessage {
	online: boolean
	stats?: {
		CPUPercent: number
		MemoryMB: number
		ThreadCount: number
		Uptime: number
	}
}

export interface ServerSocketHandlers {
	onStdout: (message: string) => void
	onStderr: (message: string) => void
	onStats: (stats: ServerStatsMessage) => void
}

export function connectToServerSockets(
	serverId: string,
	handlers: ServerSocketHandlers
): {
	stdout: WebSocket
	stderr: WebSocket
	stats: WebSocket
	stdin: WebSocket
	close: () => void
} {
	const stdoutSocket = new WebSocket(toWebSocketUrl(baseUrl, 'server:stdout:' + serverId))
	const stderrSocket = new WebSocket(toWebSocketUrl(baseUrl, 'server:stderr:' + serverId))
	const statsSocket = new WebSocket(toWebSocketUrl(baseUrl, 'server:stats:' + serverId))
	const stdinSocket = new WebSocket(toWebSocketUrl(baseUrl, 'server:stdin:' + serverId))

	stdoutSocket.onopen = () => debug('Connected to stdout channel for server ' + serverId)
	stderrSocket.onopen = () => debug('Connected to stderr channel for server ' + serverId)
	statsSocket.onopen = () => debug('Connected to stats channel for server ' + serverId)
	stdinSocket.onopen = () => debug('Connected to stdin channel for server ' + serverId)

	stdoutSocket.onmessage = (event) => handlers.onStdout(event.data)
	stderrSocket.onmessage = (event) => handlers.onStderr(event.data)
	statsSocket.onmessage = (event) => handlers.onStats(JSON.parse(event.data))

	stdoutSocket.onclose = () => debug('Disconnected from stdout channel for server ' + serverId)
	stderrSocket.onclose = () => debug('Disconnected from stderr channel for server ' + serverId)
	statsSocket.onclose = () => debug('Disconnected from stats channel for server ' + serverId)
	stdinSocket.onclose = () => debug('Disconnected from stdin channel for server ' + serverId)

	return {
		close: () => {
			stdoutSocket.close()
			stderrSocket.close()
			statsSocket.close()
			stdinSocket.close()
		},
		stdout: stdoutSocket,
		stderr: stderrSocket,
		stats: statsSocket,
		stdin: stdinSocket
	}
}

export function connectToPlayersSocket(
	serverId: string,
	handlers: {
		onPlayerJoin: (playerName: string) => void
		onPlayerLeave: (playerName: string) => void
	}
): {
	playerEvents: WebSocket
	close: () => void
} {
	const playerEventsSocket = new WebSocket(toWebSocketUrl(baseUrl, 'server:players:' + serverId))

	playerEventsSocket.onopen = () =>
		debug('Connected to player events channel for server ' + serverId)

	playerEventsSocket.onmessage = (event) => {
		const type = event.data.split(':')[0]
		const playerName = event.data.split(':')[1]
		if (type === 'join') {
			handlers.onPlayerJoin(playerName)
		} else if (type === 'leave') {
			handlers.onPlayerLeave(playerName)
		} else {
			debug('Unknown player event type: ' + type)
		}
	}

	playerEventsSocket.onclose = () =>
		debug('Disconnected from player events channel for server ' + serverId)

	return {
		close: () => playerEventsSocket.close(),
		playerEvents: playerEventsSocket
	}
}
