<script lang="ts">
	import ServerSidebar from '$lib/components/ServerSidebar.svelte'
	import { onMount } from 'svelte'
	import { getBackups, backupServer, deleteServerBackup, restoreServerBackup } from '$lib/backups'
	import Spinner from '$lib/components/Spinner.svelte'
	import { error as error_ } from '$lib/logging'

	let serverId: string = localStorage.getItem('selectedServer') ?? ''
	let backups: string[] = []
	let loading = false
	let error: string | null = null
	let success: string | null = null

	async function loadBackups() {
		loading = true
		error = null
		try {
			backups = await getBackups(serverId)
		} catch (e) {
			error_('Failed to load backups:', e)
			error = 'Failed to load backups.'
		} finally {
			loading = false
		}
	}

	async function handleBackup() {
		loading = true
		error = null
		try {
			await backupServer(serverId)
			success = 'Backup created successfully!'
			await loadBackups()
			setTimeout(() => (success = null), 2000)
		} catch (e) {
			error_('Failed to create backup:', e)
			error = 'Failed to create backup.'
		} finally {
			loading = false
		}
	}

	async function handleRestore(backup: string) {
		loading = true
		error = null
		try {
			await restoreServerBackup(serverId, backup)
			success = 'Backup restored successfully!'
			setTimeout(() => (success = null), 2000)
		} catch (e) {
			error_('Failed to restore backup:', e)
			error = 'Failed to restore backup.'
		} finally {
			loading = false
		}
	}

	async function handleDelete(backup: string) {
		loading = true
		error = null
		try {
			await deleteServerBackup(serverId, backup)
			success = 'Backup deleted.'
			await loadBackups()
			setTimeout(() => (success = null), 2000)
		} catch (e) {
			error_('Failed to delete backup:', e)
			error = 'Failed to delete backup.'
		} finally {
			loading = false
		}
	}

	onMount(loadBackups)
</script>

<ServerSidebar />

<div class="flex min-h-screen w-full flex-col items-center justify-center bg-white px-0">
	<div
		class="animate-fade-in w-full max-w-3xl rounded-xl border border-gray-200 bg-white p-10 shadow-lg"
	>
		<h2
			class="mb-8 flex items-center justify-center gap-2 text-center text-2xl font-semibold text-gray-900"
		>
			<span class="material-symbols-outlined text-3xl text-gray-700">cloud_sync</span>
			Server Backups
		</h2>
		<p class="mb-8 text-center text-base text-gray-500">
			Manage your server backups below. Create, restore, or delete backups as needed.
		</p>

		{#if loading}
			<div class="flex justify-center py-8"><Spinner size={28} /></div>
		{/if}

		{#if error}
			<div class="animate-fade-in-slow mb-4 text-center text-base font-medium text-red-500">
				{error}
			</div>
		{/if}

		{#if success}
			<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/10">
				<div
					class="animate-fade-in flex flex-col items-center gap-2 rounded-xl bg-white px-8 py-6 shadow-lg"
				>
					<span class="material-symbols-outlined text-4xl text-green-600">check_circle</span>
					<div class="text-lg font-semibold text-green-700">{success}</div>
				</div>
			</div>
		{/if}

		<div class="mb-8 flex justify-end">
			<button
				class="flex items-center gap-2 rounded-xl bg-blue-600 px-6 py-2 font-semibold text-white shadow transition-all duration-300 hover:bg-blue-700 focus:ring-1 focus:ring-blue-400 focus:outline-none active:scale-98"
				on:click={handleBackup}
				disabled={loading}
			>
				<span class="material-symbols-outlined text-xl">add_circle</span>
				Create Backup
			</button>
		</div>

		{#if backups.length > 0}
			<ul class="grid gap-4">
				{#each backups as backup}
					<li
						class="animate-fade-in-slow flex items-center justify-between rounded-lg border border-gray-200 bg-gray-50 px-6 py-4 shadow-sm transition-all duration-300 hover:shadow-md"
					>
						<div class="flex items-center gap-3">
							<span class="material-symbols-outlined text-2xl text-gray-700">folder</span>
							<span class="text-base font-medium text-gray-900">{backup}</span>
						</div>
						<div class="flex gap-2">
							<button
								class="flex items-center gap-1 rounded border border-green-200 bg-green-50 px-3 py-1 text-green-700 shadow-sm transition-all duration-300 hover:bg-green-100"
								title="Restore"
								on:click={() => handleRestore(backup)}
								disabled={loading}
							>
								<span class="material-symbols-outlined text-base">restore</span>
								Restore
							</button>
							<button
								class="flex items-center gap-1 rounded border border-red-200 bg-red-50 px-3 py-1 text-red-700 shadow-sm transition-all duration-300 hover:bg-red-100"
								title="Delete"
								on:click={() => handleDelete(backup)}
								disabled={loading}
							>
								<span class="material-symbols-outlined text-base">delete</span>
								Delete
							</button>
						</div>
					</li>
				{/each}
			</ul>
		{:else}
			<div class="animate-fade-in-slow flex flex-col items-center justify-center py-12">
				<span class="material-symbols-outlined mb-4 text-4xl text-gray-300">cloud_off</span>
				<p class="text-base text-gray-500">No backups found for this server.</p>
			</div>
		{/if}
	</div>
</div>

<style>
	@keyframes fade-in-slow {
		from {
			opacity: 0;
			transform: translateY(10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
	.animate-fade-in-slow {
		animation: fade-in-slow 1s cubic-bezier(0.4, 0, 0.2, 1);
	}

	@keyframes fade-in {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	.animate-fade-in {
		animation: fade-in 0.2s ease-out forwards;
	}
</style>
