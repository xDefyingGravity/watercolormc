import { writable } from 'svelte/store'
import type { Server } from '$lib/types/server'

export const servers = writable<Server[]>([])

export const pages = [
	{ name: 'servers', route: '/' },
	{ name: 'create', route: '/create-server' },
	{ name: 'settings', route: '/settings' }
] as const

export type page = (typeof pages)[number]

export const currentPage = writable<page>(pages[0])

export const ipinfo = writable<{
	privateIp: string
	publicIp: string
}>({
	privateIp: '',
	publicIp: ''
})

export const backendStatus = writable<'ok' | 'error'>('ok')
