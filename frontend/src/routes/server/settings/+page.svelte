<script lang="ts">
	import ServerSidebar from '$lib/components/ServerSidebar.svelte'
	import Spinner from '$lib/components/Spinner.svelte'
	import { getServerSettings, updateServerSettings, type ServerSettings } from '$lib/settings/'
	import { onMount } from 'svelte'
	import { fly } from 'svelte/transition'

	let loading = true
	let saving = false
	let saveSuccess = false
	let saveError = ''
	const serverId = localStorage.getItem('selectedServer')

	if (!serverId) {
		window.location.href = '/servers'
	}

	let settings: ServerSettings = {
		Versions: {
			WatercolorVersion: '0.1.0',
			MinecraftVersion: '1.21.7'
		},
		JavaSettings: {
			Memory: {
				Min: 2048,
				Max: 4096
			},
			JavaPath: '',
			JvmArgs: []
		}
	}

	onMount(async () => {
		if (!serverId) return
		try {
			settings = await getServerSettings(serverId)
			console.log(JSON.stringify(settings, null, 2))
		} catch (e) {
			console.error('Failed to load server settings:', e)
			loading = false
			saveError = 'Failed to load server settings.'
		} finally {
			loading = false
		}
	})

	function validate(): boolean {
		const min = Number(settings.JavaSettings.Memory.Min)
		const max = Number(settings.JavaSettings.Memory.Max)

		if (!Number.isInteger(min) || min < 512 || min % 512 !== 0) {
			saveError = 'Min Memory must be a multiple of 512 and at least 512 MB.'
			return false
		}

		if (!Number.isInteger(max) || max < 512 || max % 512 !== 0) {
			saveError = 'Max Memory must be a multiple of 512 and at least 512 MB.'
			return false
		}

		if (max < min) {
			saveError = 'Max Memory must be greater than or equal to Min Memory.'
			return false
		}

		return true
	}

	async function save() {
		saveError = ''
		if (!validate()) return
		saving = true
		try {
			settings.JavaSettings.Memory.Min = Number(settings.JavaSettings.Memory.Min)
			settings.JavaSettings.Memory.Max = Number(settings.JavaSettings.Memory.Max)
			await updateServerSettings(serverId as string, settings)
			saveSuccess = true
			setTimeout(() => (saveSuccess = false), 1500)
		} catch (e) {
			saveSuccess = false
			saveError =
				e instanceof Error
					? e.message[0].toLocaleUpperCase() + e.message.slice(1)
					: 'Failed to save server settings.'
		} finally {
			saving = false
		}
	}
</script>

<ServerSidebar />
<div
	class="relative flex min-h-screen min-w-full flex-col items-center justify-center bg-white px-2 font-sans text-gray-900"
>
	{#if loading}
		<div class="flex flex-col items-center justify-center gap-4">
			<Spinner size={32} color="#2563eb" />
			<div class="text-lg text-gray-600">Loading settings...</div>
		</div>
	{:else}
		<form
			class="flex w-full max-w-[420px] min-w-[320px] flex-col items-center justify-center gap-4 rounded-xl border border-gray-100 bg-white px-6 py-8 text-center shadow"
			on:submit|preventDefault={save}
			in:fly={{ y: 16, duration: 300 }}
		>
			<div class="mb-4 flex items-center gap-2 text-xl font-semibold tracking-tight text-blue-800">
				<span class="material-symbols-outlined text-2xl text-blue-700">settings</span>
				Server Settings
			</div>
			<div class="flex w-full flex-col items-start gap-1">
				<label class="ml-1 text-xs font-medium text-gray-500" for="">JVM Options</label>
				<input
					type="text"
					class="w-full rounded border border-gray-200 px-3 py-2 text-base transition focus:ring-2 focus:ring-blue-200 focus:outline-none"
					bind:value={settings.JavaSettings.JvmArgs}
				/>
			</div>
			<div class="flex w-full flex-col items-start gap-1">
				<label class="ml-1 text-xs font-medium text-gray-500" for="">Min Memory (MB)</label>
				<input
					type="text"
					class="w-full rounded border border-gray-200 px-3 py-2 text-base transition focus:ring-2 focus:ring-blue-200 focus:outline-none"
					bind:value={settings.JavaSettings.Memory.Min}
				/>
			</div>
			<div class="flex w-full flex-col items-start gap-1">
				<label class="ml-1 text-xs font-medium text-gray-500" for="">Max Memory (MB)</label>
				<input
					type="text"
					class="w-full rounded border border-gray-200 px-3 py-2 text-base transition focus:ring-2 focus:ring-blue-200 focus:outline-none"
					bind:value={settings.JavaSettings.Memory.Max}
				/>
			</div>
			<div class="flex w-full flex-col items-start gap-1">
				<label class="ml-1 text-xs font-medium text-gray-500" for="">Java Path</label>
				<input
					type="text"
					class="w-full rounded border border-gray-200 px-3 py-2 text-base transition focus:ring-2 focus:ring-blue-200 focus:outline-none"
					bind:value={settings.JavaSettings.JavaPath}
				/>
			</div>
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
