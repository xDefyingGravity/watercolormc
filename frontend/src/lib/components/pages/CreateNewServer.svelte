<script lang="ts">
	import { fade, fly } from 'svelte/transition'
	import SearchBar from '$lib/components/utilities/SearchBar.svelte'
	import { versions } from '$lib/versions'
	import { createServer } from '$lib/tasks/create'
	import { goto } from '$app/navigation'

	let name = ''
	let version = ''
	let description = ''
	let loading = false
	let error = ''

	function validate() {
		if (!name.trim() || !version.trim() || !description.trim()) {
			error = 'Please fill all fields.'
			return false
		}
		if (!versions.includes(version)) {
			error = 'Invalid version.'
			return false
		}
		return true
	}

	async function submit() {
		error = ''
		if (!validate()) return
		loading = true
		try {
			const result = await createServer(name, description, version)
			if (result.failed) {
				error = result.message || 'Failed to create server. Please try again.'
				return
			}
			name = ''
			version = ''
			description = ''
			await goto('/')
		} catch (e) {
			console.error('Error creating server:', e)
			error = 'Failed to create server. Please try again.'
		} finally {
			loading = false
		}
	}
</script>

<div class="animate-fadein flex min-h-screen items-center justify-center p-6">
	<form
		on:submit|preventDefault={submit}
		class="animate-slidein w-full max-w-md space-y-6 rounded-lg bg-white p-8 shadow-md"
		transition:fly={{ y: 40, duration: 500 }}
		autocomplete="off"
	>
		<h1 class="mb-6 text-center text-2xl font-bold text-gray-900" transition:fade>
			Create New Server
		</h1>

		<label class="block" transition:fade>
			<span class="mb-1 block font-medium text-gray-700">Name</span>
			<input
				type="text"
				bind:value={name}
				placeholder="Server Name"
				class="w-full rounded-md border border-gray-300 px-3 py-2 transition-all duration-300 focus:ring-2 focus:ring-blue-500 focus:outline-none"
				required
			/>
		</label>

		<label class="block" transition:fade>
			<span class="mb-1 block font-medium text-gray-700">Version</span>
			<SearchBar
				options={versions}
				placeholder="Select version..."
				onSelect={(v) => (version = v)}
			/>
		</label>

		<label class="block" transition:fade>
			<span class="mb-1 block font-medium text-gray-700">Description</span>
			<textarea
				bind:value={description}
				placeholder="Short description"
				class="w-full rounded-md border border-gray-300 px-3 py-2 transition-all duration-300 focus:ring-2 focus:ring-blue-500 focus:outline-none"
				required
				rows="3"
			></textarea>
		</label>

		{#if error}
			<div
				class="animate-shake rounded border border-red-300 bg-red-100 px-3 py-2 text-sm text-red-600"
				transition:fade
			>
				{error}
			</div>
		{/if}

		<button
			type="submit"
			class="mt-4 flex w-full items-center justify-center gap-2 rounded-md bg-blue-600 py-3 font-semibold text-white transition-all duration-300 hover:bg-blue-700 active:bg-blue-800 disabled:opacity-60"
			disabled={loading}
			transition:fade
		>
			{#if loading}
				<svg class="h-5 w-5 animate-spin text-white" viewBox="0 0 24 24" fill="none">
					<circle
						class="opacity-25"
						cx="12"
						cy="12"
						r="10"
						stroke="currentColor"
						stroke-width="4"
					/>
					<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
				</svg>
				Creating...
			{:else}
				Create Server
			{/if}
		</button>
	</form>
</div>
