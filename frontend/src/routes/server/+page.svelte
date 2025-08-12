<script lang="ts">
	import { servers, ipinfo } from '$lib/stores'
	import { get } from 'svelte/store'
	import { fade } from 'svelte/transition'
	import type { Server } from '$lib/types/server'
	import { baseUrl } from '$lib/config'
	import { error } from '$lib/logging'
	import { safeFetch } from '$lib/utils/fetch'
	import { onDestroy, onMount, tick } from 'svelte'
	import ConnectionInfo from '$lib/components/ConnectionInfo.svelte'
	import ServerHeader from '$lib/components/ServerHeader.svelte'
	import ServerControls from '$lib/components/ServerControls.svelte'
	import ServerStats from '$lib/components/ServerStats.svelte'

	import {
		connectToServerSockets,
		type ServerStatsMessage,
		type ServerLog
	} from '$lib/utils/serverSockets'
	import { highlight } from '$lib/highlight'
	import ServerSidebar from '$lib/components/ServerSidebar.svelte'

	let server = $state<Server | null>(null)
	let stats = $state<ServerStatsMessage | null>(null)

	let serverLogs = $state<ServerLog[]>([])

	let logsElement: HTMLDivElement | null = $state(null)

	const storedId = localStorage.getItem('selectedServer')
	const allServers = get(servers)

	server = storedId ? allServers.find((x) => x.id === storedId) || allServers[0] : allServers[0]

	const statusColors = {
		online: 'bg-green-400',
		offline: 'bg-red-300',
		unknown: 'bg-gray-300'
	}
	const statusLabels = {
		online: 'Online',
		offline: 'Offline',
		unknown: 'Unknown'
	}

	let socketsConnection: {
		close: () => void
		stdout: WebSocket
		stderr: WebSocket
		stats: WebSocket
		stdin: WebSocket
	} | null = null

	let errorStartingServerModal = $state(false)
	let startingServerError = $state('')

	function setupSockets() {
		if (socketsConnection) socketsConnection.close()
		socketsConnection = connectToServerSockets(server!.id, {
			onStdout: (message: string) => {
				serverLogs.push({ message, type: 'stdout' })
				tick().then(() => {
					logsElement?.scrollTo({ top: logsElement.scrollHeight + 400, behavior: 'smooth' })
				})
			},
			onStderr: (message: string) => {
				serverLogs.push({ message, type: 'stderr' })
				tick().then(() => {
					logsElement?.scrollTo({ top: logsElement.scrollHeight + 400, behavior: 'smooth' })
				})
			},
			onStats: (data: ServerStatsMessage) => {
				if (data.online) {
					stats = data
				} else {
					setServerStatus(server!.id, 'offline')
					disconnectFromSockets()
					stats = null
				}
			}
		})
	}

	let commandInput = $state('')
	function sendCommandToServer(command: string) {
		if (socketsConnection && socketsConnection.stdin) {
			socketsConnection.stdin.send(command)
		} else {
			error('Server is not connected')
		}
	}

	function disconnectFromSockets() {
		if (socketsConnection) {
			socketsConnection.close()
			socketsConnection = null
		}
	}

	// Loading states
	let isStarting = $state(false)
	let isStopping = $state(false)

	let showEulaModal = $state(false)
	let eulaAccepted = false

	function setServerStatus(id: string, status: 'online' | 'offline' | 'unknown') {
		servers.update((currentServers) =>
			currentServers.map((s) => (s.id === id ? { ...s, status } : s))
		)
		server!.status = status
	}

	async function startServer() {
		if (!server) {
			error('No server selected')
			return
		}

		if (server.status === 'online') {
			error('Server is already online')
			return
		}

		serverLogs = []
		if (!eulaAccepted && localStorage.getItem('server:eulaAccepted:' + server!.id) !== 'true') {
			showEulaModal = true
			return
		}

		isStarting = true
		try {
			const response = await safeFetch<string>(`${baseUrl}/api/servers/start/${server!.id}`, {
				method: 'POST'
			})
			if (response !== 'ok') {
				startingServerError = response || 'Failed to start server'
				errorStartingServerModal = true
				error(response || 'Failed to start server')
			}
			setupSockets()
			setServerStatus(server!.id, 'online')
		} finally {
			isStarting = false
		}
	}

	function acceptEulaAndStart() {
		eulaAccepted = true
		showEulaModal = false
		localStorage.setItem('server:eulaAccepted:' + server!.id, 'true')
		startServer()
	}

	function closeEulaModal() {
		showEulaModal = false
	}

	async function stopServer() {
		isStopping = true
		try {
			const timeStart = performance.now()
			const response = await safeFetch<string>(`${baseUrl}/api/servers/stop/${server!.id}`, {
				method: 'POST'
			})
			if (response !== 'ok') {
				error(response || 'Failed to stop server')
			}
			await new Promise((res) => setTimeout(res, 1000 - (performance.now() - timeStart)))
		} catch (err) {
			error('Error stopping server:', err)
		} finally {
			isStopping = false
		}
	}

	async function getPreviousServerLogs(id: string): Promise<ServerLog[]> {
		const logs = await safeFetch<string>(baseUrl + `/api/servers/logs/${id}`, {
			method: 'GET'
		})
		if (!logs || !Array.isArray(logs)) {
			error('Failed to fetch previous server logs')
			return []
		}
		return logs.map((log: string) => ({
			message: log,
			type: 'stdout'
		}))
	}

	onMount(async () => {
		if (server?.status === 'online') {
			serverLogs = [...serverLogs, ...((await getPreviousServerLogs(server.id)) || [])]
			setupSockets()
			await tick()
			logsElement?.scrollTo({ top: logsElement.scrollHeight + 400, behavior: 'smooth' })
		}
	})

	onDestroy(() => {
		disconnectFromSockets()
	})
</script>

{#if server}
	<ServerSidebar />
	<div class="mx-auto flex min-h-screen max-w-3xl flex-col p-6" transition:fade={{ duration: 250 }}>
		<ServerHeader
			name={server.name}
			status={server.status}
			description={server.description}
			{statusColors}
			{statusLabels}
		/>

		<ConnectionInfo
			host={server.host}
			port={server.port}
			version={server.version}
			privateIp={$ipinfo.privateIp}
		/>

		<ServerControls
			status={server.status}
			{isStarting}
			{isStopping}
			onStart={startServer}
			onStop={stopServer}
		/>

		<ServerStats {stats} status={server.status} />

		<section transition:fade={{ duration: 300, delay: 220 }}>
			<h2 class="mb-4 text-xl font-semibold">Server Console</h2>
			<div
				class="selectable max-h-64 space-y-2 overflow-y-auto rounded-lg border border-gray-200 bg-gray-50 px-4 py-3 font-mono text-sm"
				bind:this={logsElement}
				style="min-height: 160px;"
			>
				{#if serverLogs.length === 0}
					{#if server.status === 'online'}
						<div class="text-gray-500">No logs yet.</div>
					{:else}
						<div class="text-gray-500">Server is offline. Start the server to see logs.</div>
					{/if}
				{:else}
					{#each serverLogs as log}
						<div
							class={`selectable whitespace-pre-wrap ${log.type === 'stderr' ? 'text-red-600' : 'text-gray-800'}`}
						>
							{@html highlight(log.message)}
						</div>
					{/each}
				{/if}
			</div>
			<form
				class="mt-3 flex gap-2"
				onsubmit={(e) => {
					e.preventDefault()
					if (commandInput.trim() && server?.status === 'online') {
						sendCommandToServer(commandInput.trim())
						commandInput = ''
					}
				}}
				autocomplete="off"
			>
				<input
					class="flex-1 rounded border border-gray-300 px-3 py-2 font-mono text-sm transition focus:ring-2 focus:ring-blue-200 focus:outline-none"
					type="text"
					placeholder="Type a command and press Enter"
					bind:value={commandInput}
					autocomplete="off"
					spellcheck="false"
					aria-label="Send command to server"
					disabled={server?.status !== 'online'}
				/>
				<button
					class="rounded bg-blue-600 px-4 py-2 text-white transition hover:bg-blue-700 disabled:opacity-60"
					type="submit"
					disabled={server?.status !== 'online' || !commandInput.trim()}
				>
					{server?.status === 'online' ? 'Send' : 'Server Offline'}
				</button>
			</form>
		</section>
	</div>
{:else}
	<div class="p-8 text-center text-gray-500">
		{@html $servers.length === 0
			? 'No servers available. Please <a href="/create-server" class="text-blue-600 hover:underline">add a server</a> first.'
			: 'Server not found. Please select a valid server.'}
	</div>
{/if}

{#if showEulaModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40" transition:fade>
		<div class="animate-pop w-full max-w-md rounded-lg bg-white p-8 shadow-lg">
			<h2 class="mb-4 text-xl font-bold">End User License Agreement</h2>
			<div class="mb-4 max-h-48 overflow-y-auto text-gray-700">
				Please read and accept the <a
					href="https://www.minecraft.net/eula"
					class="text-blue-500 underline">Minecraft EULA</a
				> before starting the server.
			</div>
			<div class="flex justify-end gap-2">
				<button class="rounded bg-gray-200 px-4 py-2 hover:bg-gray-300" onclick={closeEulaModal}
					>Cancel</button
				>
				<button
					class="rounded bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
					onclick={acceptEulaAndStart}>Accept</button
				>
			</div>
		</div>
	</div>
{/if}

{#if errorStartingServerModal}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 transition-opacity duration-200"
	>
		<div class="animate-pop relative w-full max-w-md rounded-xl bg-white p-8 shadow-2xl">
			<button
				class="absolute top-4 right-4 text-gray-400 transition hover:text-gray-600"
				aria-label="Close"
				onclick={() => {
					startingServerError = ''
					errorStartingServerModal = false
				}}
				type="button"
			>
				<svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
				</svg>
			</button>
			<div class="flex flex-col items-center gap-2">
				<div class="mb-2 flex h-12 w-12 items-center justify-center rounded-full bg-red-50">
					<svg
						class="h-7 w-7 text-red-500"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						viewBox="0 0 24 24"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</div>
				<h2 class="mb-1 text-center text-xl font-bold text-red-600">Server Start Failed</h2>
				<p class="mb-2 text-center text-base whitespace-pre-wrap text-red-600">
					{startingServerError}
				</p>
				<button
					class="mt-4 w-full rounded bg-gray-100 px-5 py-2.5 font-medium text-gray-700 transition hover:bg-gray-200"
					onclick={() => {
						startingServerError = ''
						errorStartingServerModal = false
					}}
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	button:focus {
		outline: none;
		box-shadow: 0 0 0 2px #60a5fa;
	}
	section:not(:last-child) {
		border-bottom: 1px solid #f3f4f6;
		padding-bottom: 2rem;
	}

	@keyframes pop {
		0% {
			transform: scale(0.8);
			opacity: 0;
		}
		100% {
			transform: scale(1);
			opacity: 1;
		}
	}
	.animate-pop {
		animation: pop 0.2s cubic-bezier(0.4, 0, 0.2, 1);
	}
</style>
