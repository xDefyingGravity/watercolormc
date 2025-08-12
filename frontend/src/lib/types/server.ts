export interface Server {
	type: 'vanilla' | 'papermc'
	id: string
	name: string
	port: number
	host: string
	status: 'online' | 'offline' | 'unknown'
	description: string
	version: string
	createdAt: string
}
