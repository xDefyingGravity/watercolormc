<script lang="ts">
	import ServerSidebar from '$lib/components/ServerSidebar.svelte'
	import { onMount } from 'svelte'
	import { getPlugin, type Project } from '$lib/modrinth/plugins'
	import Markdown from '$lib/components/Markdown.svelte'
	import Spinner from '$lib/components/Spinner.svelte'
	import { goto } from '$app/navigation'
	import { addPluginGraph, removePlugin } from '$lib/paper/plugins'
	import { servers } from '$lib/stores'
	import { getBaseVersion } from '$lib/paper/version'
	import { getManifest, type Manifest } from '$lib/paper/manifest'

	const searchParams = new URLSearchParams(window.location.search)
	const pluginId = searchParams.get('i') || ''
	let plugin: Project | null = null
	let loading = true
	let addLoading = false;
	let addError = '';
	let addSuccess = false;

	let removeLoading = false;
	let removeError = '';
	let removeSuccess = false;

	let isInstalled = false;

	const serverId = localStorage.getItem('selectedServer')
	if (!serverId) {
		goto('/servers');
	}

	const server = $servers.find(s => s.id === serverId);
	if (!server) {
		goto('/servers');
	}

	let manifest: Manifest | null = null;

	onMount(async () => {
		if (pluginId) {
			loading = true;
			try {
				plugin = await getPlugin(pluginId);
				if (!plugin) {
					goto('/server/plugins');
					return;
				}

				manifest = await getManifest(serverId!);
				isInstalled = manifest.plugins.some(p => p.id === plugin!.id);
			} catch (error) {
				console.error('Failed to fetch plugin:', error);
				addError = 'Failed to fetch plugin.';
			} finally {
				loading = false;
			}
		}
	});

	async function addToServer() {
		if (!plugin) return;
		addLoading = true;
		addError = '';
		addSuccess = false;
		try {
			await addPluginGraph(serverId!, plugin.id, getBaseVersion(server?.version || ""));
			addSuccess = true;
			setTimeout(() => addSuccess = false, 2000);
		} catch (e) {
			addError = 'Failed to add plugin.';
			console.error(e);
		} finally {
			addLoading = false;
		}
	}

	async function removeFromServer() {
		removeLoading = true;
		removeError = '';
		removeSuccess = false;
		try {
			const jarName = manifest?.plugins.find(p => p.id === plugin?.id)?.jar_name || '';
			await removePlugin(serverId!, jarName);
			removeSuccess = true;
			setTimeout(() => removeSuccess = false, 2000);
		} catch (e) {
			removeError = 'Failed to remove plugin.';
			console.error(e);
		} finally {
			removeLoading = false;
		}
	}
</script>

<ServerSidebar />

<div class="plugin-page">
	{#if loading}
		<div class="flex h-screen items-center justify-center">
			<Spinner size={64} color="blue" />
		</div>
	{:else if plugin}
		<div class="flex min-h-screen w-full flex-col items-center justify-center bg-gray-50">
			<div class="plugin-card w-full max-w-4xl rounded-3xl bg-white p-8 shadow-xl">
				<div class="mb-6 flex items-center gap-6">
					<img
						src={plugin.icon_url}
						alt={plugin.title}
						class="h-20 w-20 rounded-xl border bg-gray-100 object-contain"
					/>
					<div class="flex flex-col gap-2">
						<h1 class="text-2xl font-bold text-gray-900">{plugin.title}</h1>
						<p class="text-base text-gray-600">{plugin.description}</p>
						<div class="flex gap-4 text-sm text-gray-500">
							<span
								>Downloads: <span class="font-semibold text-gray-700"
									>{plugin.downloads.toLocaleString()}</span
								></span
							>
							<span
								>Followers: <span class="font-semibold text-gray-700"
									>{plugin.followers.toLocaleString()}</span
								></span
							>
							<span
								>Type: <span class="font-semibold text-gray-700">{plugin.project_type}</span></span
							>
						</div>
					</div>
				</div>
				<div class="mb-6">
					<h2 class="mb-2 text-lg font-semibold text-gray-800">About</h2>
					<p class="whitespace-pre-line text-gray-700">
						<Markdown markdown={plugin.body} />
					</p>
				</div>
				<div class="mb-6">
					<h2 class="mb-2 text-lg font-semibold text-gray-800">Game Versions</h2>
					<div class="flex flex-wrap gap-2">
						{#each plugin.game_versions as version}
							<span class="rounded bg-blue-50 px-2 py-1 text-xs font-medium text-blue-700"
								>{version}</span
							>
						{/each}
					</div>
				</div>
				{#if plugin.gallery && plugin.gallery.length}
					<div class="mb-6">
						<h2 class="mb-2 text-lg font-semibold text-gray-800">Gallery</h2>
						<div class="flex gap-4 overflow-x-auto pb-2">
							{#each plugin.gallery as image}
								<img src={image.url} alt={image.title} class="h-32 rounded-lg border shadow" />
							{/each}
						</div>
					</div>
				{/if}
				<div class="mt-8 flex items-center justify-between">
					<div class="flex gap-4">
						{#if plugin.source_url}
							<a href={plugin.source_url} target="_blank" rel="noopener noreferrer" class="btn-link"
								>Source</a
							>
						{/if}
						{#if plugin.wiki_url}
							<a href={plugin.wiki_url} target="_blank" rel="noopener noreferrer" class="btn-link"
								>Wiki</a
							>
						{/if}
						{#if plugin.issues_url}
							<a href={plugin.issues_url} target="_blank" rel="noopener noreferrer" class="btn-link"
								>Issues</a
							>
						{/if}
						{#if plugin.discord_url}
							<a
								href={plugin.discord_url}
								target="_blank"
								rel="noopener noreferrer"
								class="btn-link">Discord</a
							>
						{/if}
					</div>
					{#if !isInstalled}
						<button
							class="rounded bg-blue-500 px-4 py-2 text-white transition-colors duration-300 hover:bg-blue-600 flex items-center gap-2 relative"
							onclick={addToServer}
							disabled={addLoading}
						>
							{#if addLoading}
								<Spinner size={16} color="white" />
							{/if}
							<span>{addLoading ? 'Adding...' : 'Add To Server'}</span>
						</button>
					{:else}
						<button
							class="rounded bg-red-500 px-4 py-2 text-white transition-colors duration-300 hover:bg-red-600 flex items-center gap-2 relative"
							onclick={removeFromServer}
							disabled={removeLoading}
						>
							{#if removeLoading}
								<Spinner size={16} color="white" />
							{/if}
							<span>{removeLoading ? 'Removing...' : 'Remove From Server'}</span>
						</button>
					{/if}
				</div>
				{#if addError}
					<div class="error-message bg-red-100 text-red-700 rounded px-4 py-2 mt-4 shadow">{addError}</div>
				{/if}
				{#if addSuccess}
					<div class="success-overlay">
						<div class="success-popup">
							<span class="material-symbols-outlined" style="font-size:48px;color:#22c55e;">check_circle</span>
							<p class="text-lg font-semibold text-green-700">Plugin added successfully!</p>
						</div>
					</div>
				{/if}
				{#if removeError}
					<div class="error-message bg-red-100 text-red-700 rounded px-4 py-2 mt-4 shadow">{removeError}</div>
				{/if}
				{#if removeSuccess}
					<div class="success-overlay">
						<div class="success-popup">
							<span class="material-symbols-outlined" style="font-size:48px;color:#22c55e;">check_circle</span>
							<p class="text-lg font-semibold text-green-700">Plugin removed successfully!</p>
						</div>
					</div>
				{/if}
			</div>
		</div>
	{:else}
		<div class="flex h-screen items-center justify-center">
			<p class="text-lg text-gray-500">Plugin not found.</p>
		</div>
	{/if}
</div>

<style>
	.plugin-page {
		min-height: 100vh;
		width: 100vw;
		background: #f7fafc;
		position: relative;
	}
	.plugin-card {
		animation: fadeIn 0.3s ease;
	}
	.btn-link {
		background: #f1f5f9;
		color: #2563eb;
		padding: 0.5rem 1.25rem;
		border-radius: 0.5rem;
		font-weight: 500;
		transition: background 0.15s;
		text-decoration: none;
		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
	}
	.btn-link:hover {
		background: #e0e7ef;
	}
	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(16px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
	.success-overlay {
		position: fixed;
		inset: 0;
		background: rgba(255,255,255,0.7);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 50;
	}
	.success-popup {
		background: #fff;
		border-radius: 1rem;
		box-shadow: 0 4px 24px rgba(34,197,94,0.12);
		padding: 2rem 3rem;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
	}
	.error-message {
		animation: fadeIn 0.2s;
	}
</style>
