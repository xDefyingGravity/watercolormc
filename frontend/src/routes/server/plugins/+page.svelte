<script>
	import ServerSidebar from '$lib/components/ServerSidebar.svelte'
	import { goto } from '$app/navigation'
	import { servers } from '$lib/stores'
	import Plugins from '$lib/components/Plugins.svelte'

	const serverId = localStorage.getItem('selectedServer')
	if (!serverId) {
		goto('/servers')
	}

	const server = $servers.find((s) => s.id === serverId)
	if (!server) {
		goto('/servers')
	}
</script>

<ServerSidebar />

{#if server?.type === 'papermc'}
	<Plugins serverId={serverId || ''} />
{:else}
	<div class="flex min-h-screen w-full flex-col items-center justify-center bg-gray-100 px-4">
		<div class="w-full max-w-3xl rounded-lg bg-white p-6 shadow-md">
			<h2 class="mb-4 text-xl font-semibold text-gray-700">Plugins</h2>
			<p class="mb-4 text-gray-700">
				Vanilla servers do not support plugins. If you want to use plugins, consider creating a new
				server with
				<a
					href="https://papermc.io/"
					class="text-blue-400 underline"
					style="color: #60a5fa !important;"
					target="_blank"
					rel="noopener noreferrer"
				>
					Paper
				</a>.
			</p>
		</div>
	</div>
{/if}
