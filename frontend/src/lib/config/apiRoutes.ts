import { baseUrl } from '$lib/config/index'

export const apiRoutes = {
	upStatusCheck: {
		path: baseUrl + '/api/upstatus',
		method: 'GET',
		description: 'Check if the backend API is up and running'
	},
	getAllServers: {
		path: baseUrl + '/api/servers',
		method: 'GET',
		description: 'Get a list of all servers'
	},
	createServer: {
		path: baseUrl + '/api/servers',
		method: 'POST',
		description: 'Create a new server'
	},
	getIpInfo: {
		path: baseUrl + '/api/ipinfo',
		method: 'GET',
		description: 'Get IP information of the client'
	},
	channel: {
		path: baseUrl + '/channels/',
		method: 'WS' /* WebSocket, not HTTP */,
		description: 'WebSocket connection for real-time updates on server logs'
	}
}
