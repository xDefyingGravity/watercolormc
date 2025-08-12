<script lang="ts">
	import Spinner from './Spinner.svelte'
	export let stats: any
	export let status: string
</script>

<section class="mb-8" transition:fade={{ duration: 300, delay: 180 }}>
	<h2 class="mb-4 text-xl font-semibold text-gray-800">Server Stats</h2>
	{#if stats}
		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
			{#each [{ label: 'CPU Usage', value: `${(stats.stats?.CPUPercent ?? 0).toFixed(1)}%` }, { label: 'Memory Usage', value: `${(stats.stats?.MemoryMB ?? 0).toFixed(1)} MB` }, { label: 'Thread Count', value: stats.stats?.ThreadCount ?? 0 }, { label: 'Uptime', value: `${Math.floor((stats.stats?.Uptime ?? 0) / 1_000_000_000)} sec` }] as stat}
				<div class="rounded-lg bg-gray-100 p-4 shadow-sm">
					<h3 class="text-base font-medium text-gray-700">{stat.label}</h3>
					<p class="selectable mt-1 text-2xl font-bold text-gray-900">{stat.value}</p>
				</div>
			{/each}
		</div>
	{:else if status === 'online'}
		<div class="flex items-center justify-center gap-2 p-8">
			<Spinner />
			<div class="text-gray-500">Loading stats...</div>
		</div>
	{:else}
		<div class="text-gray-500">Server is offline, no stats available.</div>
	{/if}
</section>
