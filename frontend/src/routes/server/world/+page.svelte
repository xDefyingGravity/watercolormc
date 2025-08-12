<script lang="ts">
	import { onMount } from 'svelte'
	import { getWorld, type World, uploadWorld } from '$lib/world'
	import { fly, fade } from 'svelte/transition'
	import ServerSidebar from '$lib/components/ServerSidebar.svelte'
	import Spinner from '$lib/components/Spinner.svelte'
	import { servers } from '$lib/stores'

	const serverId = localStorage.getItem('selectedServer')
	if (!serverId) {
		window.location.href = '/servers'
	}

	const server = $servers.find((s) => s.id === serverId)

	let world: World | null = null
	let loading = true
	let uploading = false
	let uploadError = ''
	let showSuccess = false

	onMount(async () => {
		world = await getWorld(serverId || '')
		loading = false
		if (!world) {
			window.location.href = '/servers'
		}
	})

	async function handleUpload(event: Event) {
		uploadError = ''
		const input = event.target as HTMLInputElement
		if (!input.files || input.files.length === 0) return
		const file = input.files[0]
		try {
			uploading = true
			await uploadWorld(serverId || '', file)
			world = await getWorld(serverId || '')
			showSuccess = true
			setTimeout(() => (showSuccess = false), 1800)
		} catch (e) {
			uploadError = e instanceof Error ? e.message : 'Failed to upload world.'
		} finally {
			uploading = false
		}
	}
</script>

<ServerSidebar />
{#if server?.status === 'online'}
	<div class="flex min-h-screen w-full flex-col items-center justify-center bg-gray-100 px-4">
		<div class="w-full max-w-3xl rounded-lg bg-white p-6 shadow-md">
			<h2 class="mb-4 text-xl font-semibold text-gray-900">World Management</h2>
			<p class="mb-4 text-gray-700">
				Manage your Minecraft world for <strong>{server.name}</strong>. You can view the current
				world details or upload a new one.
			</p>

			{#if server.status === 'online'}
				<div
					class="mb-6 flex items-center rounded-md border border-yellow-300 bg-yellow-100 p-3 text-yellow-800"
					role="alert"
				>
					<span class="material-symbols-outlined font-extralight text-yellow-700">dangerous</span>
					<p class="text-sm font-medium">
						Please turn off your server before accessing or uploading worlds.
					</p>
				</div>
			{/if}

			<slot />
		</div>
	</div>
{:else}
	<div
		class="relative flex min-h-screen min-w-full flex-col items-center justify-center bg-white px-2 font-sans text-gray-900"
	>
		{#if loading || uploading}
			<div class="flex flex-col items-center justify-center gap-4">
				<Spinner size={36} color="#2563eb" />
				<div class="text-lg text-gray-600">
					{uploading ? 'Uploading world...' : 'Loading world...'}
				</div>
			</div>
		{:else if world}
			<div
				class="flex w-full max-w-[400px] min-w-[300px] flex-col items-center justify-center rounded-xl border border-gray-100 bg-white px-6 py-8 text-center shadow transition-all"
				in:fly={{ y: 16, duration: 300 }}
				out:fade
			>
				<div class="mb-4 flex items-center justify-center gap-2">
					<span class="material-symbols-outlined text-2xl text-blue-700">globe</span>
					<div class="text-xl font-semibold tracking-tight text-blue-800">{world.name}</div>
				</div>
				<div class="mb-4 flex w-full flex-col items-center justify-center gap-4">
					<div class="flex flex-col items-center">
						<div class="mb-1 text-xs font-semibold tracking-widest text-gray-400 uppercase">
							Seed
						</div>
						<div class="text-base font-medium text-blue-900">{world.seed}</div>
					</div>
					<div class="flex flex-col items-center">
						<div class="mb-1 text-xs font-semibold tracking-widest text-gray-400 uppercase">
							World Type
						</div>
						<div class="text-base font-medium text-blue-900">{world.type}</div>
					</div>
				</div>
				<label class="mt-2 flex cursor-pointer flex-col items-center gap-1">
					<span class="material-symbols-outlined text-xl text-blue-600">upload</span>
					<span class="text-xs font-medium text-blue-700">Upload New World</span>
					<input
						type="file"
						accept=".zip,.tar,.tar.gz,.mcworld"
						class="hidden"
						on:change={handleUpload}
					/>
				</label>
				{#if uploadError}
					<div class="mt-2 text-xs text-red-500">{uploadError}</div>
				{/if}
			</div>
		{:else}
			<div in:fade class="text-lg text-gray-600">World not found.</div>
		{/if}

		{#if showSuccess}
			<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/10">
				<div
					class="animate-fade-in flex flex-col items-center gap-2 rounded-xl bg-white px-8 py-6 shadow-lg"
				>
					<span class="material-symbols-outlined text-4xl text-green-600">check_circle</span>
					<div class="text-lg font-semibold text-green-700">World uploaded successfully!</div>
				</div>
			</div>
		{/if}
	</div>
	<style>
		@keyframes fade-in {
			from {
				opacity: 0;
				transform: scale(0.98);
			}
			to {
				opacity: 1;
				transform: scale(1);
			}
		}
		.animate-fade-in {
			animation: fade-in 0.3s;
		}
	</style>
{/if}
