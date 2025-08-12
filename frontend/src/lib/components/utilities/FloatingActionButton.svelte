<script lang="ts">
	export let onClick: () => void = () => {}
	export let description: string = 'Add new server'
	export let id: string = ""
</script>

<div class="fixed right-6 bottom-6">
	<button
		class="flex h-14 w-14 items-center justify-center rounded-full bg-blue-600 text-white shadow-lg transition-all duration-200 hover:scale-104 hover:bg-blue-700 focus:ring-2 focus:ring-blue-400 focus:outline-none active:bg-blue-800"
		aria-label={description}
		on:click={onClick}
		id={id}
	>
		<svg
			class="h-6 w-6"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			viewBox="0 0 24 24"
			xmlns="http://www.w3.org/2000/svg"
		>
			<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
		</svg>
	</button>
</div>
