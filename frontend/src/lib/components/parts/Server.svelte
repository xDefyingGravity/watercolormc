<script lang="ts">
	import type { Server } from '$lib/types/server'
	import ConnectPopup from '$lib/components/parts/ConnectPopup.svelte'
	import { ipinfo } from '$lib/stores'
	import { goto } from '$app/navigation'
	export let server: Server
	export let ondelete: (id: string) => void

	let showConnect = false

	function copyAddress() {
		navigator.clipboard.writeText(`${$ipinfo.privateIp}:${server.port}`)
	}
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

	function gotoServer() {
		localStorage.setItem('selectedServer', server.id)
		goto('/server')
	}
</script>

<div
	tabindex="0"
	role="button"
	class="group relative flex w-full cursor-pointer flex-col gap-3 rounded-lg border border-gray-200 bg-white p-6 shadow-sm
    transition focus-within:ring-2
    focus-within:ring-blue-200 hover:shadow-md"
	on:keydown={(e) => {
		if (e.key === 'Enter' || e.key === ' ') {
			e.preventDefault()
			gotoServer()
		}
	}}
	on:click={gotoServer}
>
	<button
		class="absolute top-3 right-3 z-10 h-8 w-8 rounded-full p-2 text-gray-400 transition hover:bg-red-50 hover:text-red-500"
		type="button"
		aria-label="Delete server"
		on:click|stopPropagation={() => ondelete(server.id)}
	>
		<span
			class="material-symbols-outlined"
			style="font-size: 16px!important; line-height: 1; display: inline-block;"
		>
			delete
		</span>
	</button>

	<div class="flex items-center gap-2">
		<span class="relative flex h-2.5 w-2.5 items-center" title={statusLabels[server.status]}>
			<span class={`h-2.5 w-2.5 rounded-full ${statusColors[server.status]}`}></span>
		</span>
		<h2 class="flex items-center gap-2 truncate text-lg text-gray-900">
			{server.name}
			<button
				class="flex items-center justify-center rounded p-1 text-gray-400 transition hover:bg-blue-50 hover:text-blue-600"
				type="button"
				aria-label="Connect"
				on:click|stopPropagation={() => (showConnect = true)}
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-4 w-4"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="2"
				>
					<path stroke-linecap="round" stroke-linejoin="round" d="M8 5v14l11-7z" />
				</svg>
			</button>
		</h2>
	</div>

	<div class="flex flex-col gap-1">
		<div class="flex items-center gap-2 text-sm text-gray-700">
			<span class="font-medium text-gray-500">Address:</span>
			<span class="truncate"
				>{server.host === '0.0.0.0' ? $ipinfo.privateIp : server.host}:{server.port}</span
			>
			<button
				class="ml-1 rounded p-1 text-gray-400 transition hover:bg-blue-50 hover:text-blue-600"
				on:click|stopPropagation={copyAddress}
				title="Copy address"
				type="button"
				aria-label="Copy server address"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-4 w-4"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="2"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2M16 8h2a2 2 0 012 2v8a2 2 0 01-2 2h-8a2 2 0 01-2-2v-2"
					/>
				</svg>
			</button>
		</div>
		<div class="text-sm text-gray-700">
			<span class="font-medium text-gray-500">Version:</span>
			{server.version}
		</div>
		{#if server.description}
			<div class="truncate text-sm text-gray-400">{server.description}</div>
		{/if}
	</div>
</div>

{#if showConnect}
	<ConnectPopup
		serverName={server.name}
		port={server.port}
		host={server.host}
		onClose={() => (showConnect = false)}
	/>
{/if}

<style>
	div[role='button']:hover,
	div[role='button']:focus-within {
		box-shadow: 0 4px 24px 0 rgba(80, 80, 120, 0.08);
		border-color: #60a5fa;
	}
</style>
