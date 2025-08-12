<script lang="ts">
	import ServerSidebar from '$lib/components/ServerSidebar.svelte'
	import Spinner from '$lib/components/Spinner.svelte'
	import { getServerProperties, updateServerProperties } from '$lib/properties'
	import { validateServerProperties } from '$lib/properties/validate'
	import { onMount } from 'svelte'
	import { fly, fade } from 'svelte/transition'
	import { error } from '$lib/logging'
	import Input from '$lib/components/inputs/Input.svelte'
	import BooleanInput from '$lib/components/inputs/BooleanInput.svelte'
	import DifficultyInput from '$lib/components/inputs/DifficultyInput.svelte'
	import GamemodeInput from '$lib/components/inputs/GamemodeInput.svelte'

	let properties: Record<string, string> = {}
	let originalProperties: Record<string, string> = {}
	let validationErrors: Record<string, string> = {}
	let loading = true
	let saving = false
	let saveSuccess = false
	let saveError = ''
	const serverId = localStorage.getItem('selectedServer')

	onMount(async () => {
		if (!serverId) return
		try {
			properties = await getServerProperties(serverId)
			originalProperties = { ...properties }
		} catch (e) {
			error('Failed to load server properties', e)
		} finally {
			loading = false
		}

		window.addEventListener('keydown', async (e) => {
			if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
				e.preventDefault();
				if (!loading && !saving && JSON.stringify(properties) !== JSON.stringify(originalProperties)) {
					await save();
				}
			}
		});

	})

	function validateAll() {
		const errors = validateServerProperties(properties)
		validationErrors = {}
		errors.forEach((e) => {
			validationErrors[e.key] = e.message
		})
	}

	async function save() {
		if (saving || JSON.stringify(properties) === JSON.stringify(originalProperties)) return

		saveError = ''
		saving = true
		validateAll()
		if (Object.keys(validationErrors).length > 0) {
			saving = false
			saveError = 'Please fix validation errors.'
			return
		}
		try {
			if (!serverId) {
				saveError = 'No server selected.'
				return
			}

			await updateServerProperties(serverId, properties)
			originalProperties = { ...properties }
			saveSuccess = true
			setTimeout(() => (saveSuccess = false), 1500)
		} catch (e) {
			error('Failed to save server properties', e)
			saveError = 'Failed to save changes.'
		} finally {
			saving = false
		}
	}

	function isBoolean(value: string): boolean {
		return value === 'true' || value === 'false';
	}
</script>

<ServerSidebar />
<div
	class="relative flex min-h-screen min-w-full flex-col items-center justify-center bg-white px-2 font-sans text-gray-900"
>
	{#if loading}
		<div class="flex flex-col items-center justify-center gap-4">
			<Spinner size={32} color="#2563eb" />
			<div class="text-lg text-gray-600">Loading properties...</div>
		</div>
	{:else}
		<form
			class="flex w-full max-w-[420px] min-w-[320px] flex-col items-center justify-center gap-4 rounded-xl border border-gray-100 bg-white px-6 py-8 text-center shadow"
			on:submit|preventDefault={save}
			in:fly={{ y: 16, duration: 300 }}
			out:fade
		>
			<div class="mb-4 flex items-center gap-2 text-xl font-semibold tracking-tight text-blue-800">
				<span class="material-symbols-outlined text-2xl text-blue-700">settings</span>
				Server Properties
			</div>
			{#each Object.entries(properties) as [key, value]}
				<div class="flex w-full flex-col items-start gap-1">
					<label class="ml-1 text-xs font-medium text-gray-500" for={key}>
						{key.replace(/_/g, ' ')}
					</label>

					{#if isBoolean(value)}
						<BooleanInput
							bind:value={properties[key]}
							error={validationErrors[key]}
							onBlur={validateAll}
							id={key}
						/>
					{:else if key === 'difficulty'}
						<DifficultyInput
							bind:value={properties[key]}
							error={validationErrors[key]}
							onBlur={validateAll}
							id={key}
						/>
					{:else if key === 'gamemode'}
						<GamemodeInput
							bind:value={properties[key]}
							error={validationErrors[key]}
							onBlur={validateAll}
							id={key}
						/>
					{:else}
						<Input
							bind:value={properties[key]}
							error={validationErrors[key]}
							onBlur={validateAll}
							placeholder=""
							id={key}
						/>
					{/if}
				</div>
			{/each}
			<button
				type="submit"
				class="mt-4 flex items-center justify-center gap-2 rounded bg-blue-600 px-6 py-2 font-semibold text-white transition hover:bg-blue-700 disabled:opacity-60"
				disabled={saving}
			>
				{#if saving}
					<Spinner size={18} color="#fff" />
				{/if}
				Save Changes
			</button>
			{#if saveError}
				<div class="mt-2 text-xs text-red-500">{saveError}</div>
			{/if}
		</form>
	{/if}

	{#if saveSuccess}
		<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/10">
			<div
				class="animate-fade-in flex flex-col items-center gap-2 rounded-xl bg-white px-8 py-6 shadow-lg"
			>
				<span class="material-symbols-outlined text-4xl text-green-600">check_circle</span>
				<div class="text-lg font-semibold text-green-700">Saved!</div>
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
