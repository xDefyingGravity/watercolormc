<script lang="ts">
	export let name: string
	export let status: string
	export let description: string
	export let statusColors: Record<string, string>
	export let statusLabels: Record<string, string>
</script>

<header class="mb-8" transition:slide={{ duration: 350 }}>
	<h1 class="flex items-center gap-3 text-3xl">
		<span class="relative flex h-3 w-3 items-center" title={statusLabels[status]}>
			<span class={`h-3 w-3 rounded-full ${statusColors[status]}`}></span>
		</span>
		<span class="truncate">{name}</span>
	</h1>
	{#if description}
		<p class="mt-2 text-gray-700">{description}</p>
	{/if}
</header>
