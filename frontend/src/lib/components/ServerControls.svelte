<script lang="ts">
	import Spinner from './Spinner.svelte'
	export let status: string
	export let isStarting: boolean
	export let isStopping: boolean
	export let onStart: () => void
	export let onStop: () => void
</script>

<section class="mb-8 flex gap-4" transition:fade={{ duration: 300, delay: 120 }}>
	<button
		class="flex items-center gap-2 rounded bg-green-600 px-5 py-2 text-white shadow-sm transition-colors duration-200 hover:bg-green-700 disabled:cursor-not-allowed disabled:opacity-50"
		disabled={status === 'online' || isStarting}
		on:click={onStart}
	>
		{#if isStarting}
			<Spinner />
		{/if}
		Start
	</button>

	<button
		class="flex items-center gap-2 rounded bg-red-600 px-5 py-2 text-white shadow-sm transition-colors duration-200 hover:bg-red-700 disabled:cursor-not-allowed disabled:opacity-50"
		disabled={status !== 'online' || isStopping}
		on:click={onStop}
	>
		{#if isStopping}
			<Spinner />
		{/if}
		Stop
	</button>
</section>
