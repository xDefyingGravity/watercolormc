<script lang="ts">
	import { getPlugin, getTopPlugins, type PluginInfo } from '$lib/modrinth/plugins'
	import { onMount } from 'svelte'
	import { getManifest } from '$lib/paper/manifest'
	import type { Project } from '@xmcl/modrinth'

	let { serverId }: { serverId: string } = $props()

	let plugins: PluginInfo[] = $state<PluginInfo[]>([])
	let myPlugins: Project[] = $state<Project[]>([])

	let searchQuery = $state<string>('')

	let page: number = $state<number>(1)
	let pageSize = 10

	async function loadPlugins() {
		try {
			plugins = await getTopPlugins({
				limit: pageSize,
				offset: (page - 1) * pageSize,
				query: searchQuery
			})

			myPlugins = await Promise.all(
				(await getManifest(serverId)).plugins.map(plugin => getPlugin(plugin.id))
			)
		} catch (error) {
			console.error('Failed to fetch plugins:', error)
		}
	}

	onMount(loadPlugins)

	async function onSearch() {
		page = 1
		await loadPlugins()
	}

	function nextPage() {
		page++
		loadPlugins()
	}

	function prevPage() {
		if (page > 1) {
			page--
			loadPlugins()
		}
	}
</script>

<div class="ml-16 flex min-h-screen flex-col items-center justify-start bg-gray-100 px-6 py-12">
	<div class="mb-12 border-b-2 border-gray-200 pb-4">
		<h1 class="text-2xl font-bold text-gray-900 text-center mb-8">Installed Plugins</h1>
		{#if myPlugins.length > 0}
			{#each myPlugins as plugin}
				<a href={`/server/plugins/view?i=${plugin.id}`} class="no-underline">
					<div
						class="flex gap-4 rounded-2xl bg-white p-5 shadow-sm transition-shadow hover:shadow-md"
					>
						<img
							src={plugin.icon_url}
							alt={plugin.title}
							class="h-16 w-16 rounded-xl bg-gray-50 object-contain"
						/>
						<div>
							<h2 class="text-lg font-semibold text-gray-900">{plugin.title}</h2>
							<p class="line-clamp-2 text-sm text-gray-600">{plugin.description}</p>
							<span class="text-xs text-gray-500">{plugin.downloads.toLocaleString()} downloads</span>
						</div>
					</div>
				</a>
			{/each}
		{:else}
			<p class="text-gray-600">No plugins installed on this server. You should install some!</p>
		{/if}
	</div>

	<h1 class="text-2xl font-bold text-gray-900 mb-8">Browse Plugins</h1>

	<!-- Search Bar -->
	<div class="mb-8 flex w-full max-w-2xl">
		<input
			type="text"
			placeholder="Search plugins..."
			bind:value={searchQuery}
			onkeydown={(event) => {
				if (event.key === 'Enter') onSearch()
			}}
			class="flex-1 rounded-l-lg border border-gray-300 p-3 focus:ring-2 focus:ring-blue-500 focus:outline-none"
		/>
		<button
			class="rounded-r-lg bg-blue-600 px-6 text-white transition-colors hover:bg-blue-700"
			onclick={onSearch}
		>
			Search
		</button>
	</div>

	<div class="grid w-full max-w-6xl grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
		{#each plugins as plugin}
			<a href={`/server/plugins/view?i=${plugin.id}`} class="no-underline">
				<div
					class="flex gap-4 rounded-2xl bg-white p-5 shadow-sm transition-shadow hover:shadow-md"
				>
					<img
						src={plugin.icon}
						alt={plugin.name}
						class="h-16 w-16 rounded-xl bg-gray-50 object-contain"
					/>
					<div>
						<h2 class="text-lg font-semibold text-gray-900">{plugin.name}</h2>
						<p class="line-clamp-2 text-sm text-gray-600">{plugin.description}</p>
						<span class="text-xs text-gray-500">{plugin.downloads.toLocaleString()} downloads</span>
					</div>
				</div>
			</a>
		{/each}
	</div>

	<!-- Pagination -->
	<div class="mt-10 flex items-center gap-4">
		<button
			class="rounded bg-gray-200 px-4 py-2 text-gray-700 hover:bg-gray-300 disabled:opacity-50"
			onclick={prevPage}
			disabled={page === 1}
		>
			Previous
		</button>
		<span class="font-medium text-gray-700">Page {page}</span>
		<button
			class="rounded bg-gray-200 px-4 py-2 text-gray-700 hover:bg-gray-300"
			onclick={nextPage}
		>
			Next
		</button>
	</div>
</div>
