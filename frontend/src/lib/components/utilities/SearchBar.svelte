<script lang="ts">
	export let onSelect: (value: string) => void = () => {}
	export let options: string[] = []
	export let placeholder: string = 'Search...'

	let query = ''
	let open = false
	let filterMode: 'both' | 'vanilla' | 'paper' = 'both'

	$: filtered = options
		.filter((opt) => {
			if (filterMode === 'vanilla') return !opt.startsWith('paper-')
			if (filterMode === 'paper') return opt.startsWith('paper-')
			return true // both
		})
		.filter((opt) => (query ? opt.toLowerCase().includes(query.toLowerCase()) : true))
</script>

<div class="relative w-full max-w-sm">
	<!-- Toggle buttons -->
	<div class="mb-2 flex gap-2 [&>button]:shadow-sm">
		<button
			class="rounded px-3 py-1 transition-colors duration-200"
			class:selected={filterMode === 'both'}
			on:click={() => (filterMode = 'both')}
			type="button"
		>
			Both
		</button>
		<button
			class="rounded px-3 py-1 transition-colors duration-200"
			class:selected={filterMode === 'vanilla'}
			on:click={() => (filterMode = 'vanilla')}
			type="button"
		>
			Vanilla
		</button>
		<button
			class="rounded px-3 py-1 transition-colors duration-200"
			class:selected={filterMode === 'paper'}
			on:click={() => (filterMode = 'paper')}
			type="button"
		>
			Paper
		</button>
	</div>

	<input
		type="text"
		bind:value={query}
		{placeholder}
		on:focus={() => (open = true)}
		on:blur={() => setTimeout(() => (open = false), 100)}
		class="w-full rounded-md border border-gray-300 px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none"
	/>

	{#if open && filtered.length > 0}
		<ul
			class="absolute z-10 mt-1 max-h-60 w-full overflow-y-auto rounded-md border border-gray-300 bg-white shadow-lg"
		>
			{#each filtered as option}
				<li class="w-full">
					<button
						class="w-full cursor-pointer px-4 py-2 text-start hover:bg-blue-100"
						on:mousedown={() => {
							onSelect(option)
							query = option
							open = false
						}}
					>
						{option}
					</button>
				</li>
			{/each}
		</ul>
	{:else if open && query.trim().length > 0}
		<div
			class="absolute z-10 mt-1 w-full rounded-md border border-gray-300 bg-white px-4 py-2 text-gray-500 shadow-lg"
		>
			No matches found
		</div>
	{/if}
</div>

<style lang="postcss">
    button.selected {
        background-color: #3b82f6;
        color: white;
    }

		button {
				border: none;
		}
</style>

