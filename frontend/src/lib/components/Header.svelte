<script lang="ts">
	import { pages, currentPage } from '$lib/stores'
	import { goto } from '$app/navigation'

	function setPage(pageName: string) {
		const page = pages.find((p) => p.name === pageName)
		if (page) {
			currentPage.set(page)
			goto(page.route)
		}
	}
</script>

<nav class="mx-auto flex h-12 items-center space-x-6 bg-white p-4 shadow">
	{#each pages as page}
		<button
			class="rounded px-3 py-1 text-sm font-medium
        {$currentPage.name === page.name
				? 'bg-blue-600 text-white'
				: 'text-gray-700 hover:bg-gray-100'}"
			id="{page.name}-button"
			on:click={() => setPage(page.name)}
			aria-current={$currentPage.name === page.name ? 'page' : undefined}
		>
			{page.name.charAt(0).toUpperCase() + page.name.slice(1)}
		</button>
	{/each}
</nav>
