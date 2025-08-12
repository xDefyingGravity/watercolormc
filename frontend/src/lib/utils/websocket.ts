export function toWebSocketUrl(baseUrl: string, channel: string) {
	if (!baseUrl) {
		const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
		return `${protocol}://${window.location.host}/channels/${channel}`
	}

	const wsUrl = baseUrl.replace(/^http/, 'ws')
	return `${wsUrl}/channels/${channel}`
}
