<script lang="ts">
	import { servers } from '$lib/stores'
	import Server from '$lib/components/parts/Server.svelte'
	import { error, info } from '$lib/logging'
	import { baseUrl } from '$lib/config'
	import { safeFetch } from '$lib/utils/fetch'
	import ConfirmDeleteModal from '$lib/components/ConfirmDeleteModal.svelte'

	let showModal = false
	let serverToDelete: { id: string; name: string } | null = null

	function handleDeleteRequest(server: { id: string; name: string }) {
		serverToDelete = { id: server.id, name: server.name }
		showModal = true
	}

	async function confirmDelete() {
		if (serverToDelete) {
			await deleteServer(serverToDelete.id)
			showModal = false
			serverToDelete = null
		}
	}

	function cancelDelete() {
		showModal = false
		serverToDelete = null
	}

	async function deleteServer(serverId: string) {
		info('Deleting server with ID: ' + serverId)

		const response = await safeFetch<string>(baseUrl + `/api/servers/${serverId}`, {
			method: 'DELETE'
		})

		if (response != 'ok') {
			error('Failed to delete server: ' + response)
		} else {
			info('Server deleted successfully')
			servers.update((current) => current.filter((s) => s.id !== serverId))
		}
	}
</script>

<div class="p-6" id="server-list">
	{#if $servers.length === 0}
		<div class="py-12 text-center text-gray-500">
			<p class="text-lg font-medium">No servers found.</p>
			<p class="text-sm">You can add a new server using the + button.</p>
		</div>
	{:else}
		<div class="grid gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			{#each $servers as server (server.id)}
				<Server {server} ondelete={() => handleDeleteRequest(server)} />
			{/each}
		</div>
	{/if}
</div>

{#if showModal && serverToDelete}
	<ConfirmDeleteModal
		serverName={serverToDelete.name}
		onConfirm={confirmDelete}
		onCancel={cancelDelete}
	/>
{/if}
