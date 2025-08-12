<script lang="ts">
	import { funFacts } from '$lib/constants/funFacts'
	import { onMount, onDestroy } from 'svelte'
	import { writable } from 'svelte/store'
	import { fade } from 'svelte/transition'

	export let message: string = 'Loading...'
	export let showFunFact: boolean = true

	const cycleInterval = 5000
	const dotsInterval = 500
	const spinnerSpeed = writable(2.5)

	let dotsAmount = 0
	let maxDots = 3

	let currentFactIndex = funFacts.length > 0 ? Math.floor(Math.random() * funFacts.length) : 0
	let cycleIntervalId: ReturnType<typeof setInterval>
	let dotsIntervalId: ReturnType<typeof setInterval>

	function cycleFact() {
		currentFactIndex = (currentFactIndex + 1) % funFacts.length
	}

	function updateDots() {
		dotsAmount = (dotsAmount + 1) % (maxDots + 1)
	}

	onMount(() => {
		cycleIntervalId = setInterval(cycleFact, cycleInterval)
		dotsIntervalId = setInterval(updateDots, dotsInterval)
	})

	onDestroy(() => {
		clearInterval(cycleIntervalId)
		clearInterval(dotsIntervalId)
	})
</script>

<div
	class="flex min-h-screen flex-col items-center justify-center overflow-hidden bg-gray-100 p-6 select-none sm:p-10"
>
	<div
		class="mb-6 h-14 w-14 cursor-pointer rounded-full border-4 border-blue-600 border-t-transparent transition-transform duration-300 hover:scale-125 hover:border-purple-600"
		aria-label="Loading spinner"
		style="animation: spin linear infinite; animation-duration: {(1 / $spinnerSpeed) * 2}s"
	></div>

	<div
		class="mb-3 cursor-default text-xl font-semibold text-gray-800 transition-colors duration-300 select-text hover:text-blue-700"
		aria-live="polite"
	>
		{#if message.endsWith('...')}
			<span>
				{message.slice(0, -3).trimEnd()}{'.'.repeat(dotsAmount)}
			</span>
		{:else}
			{message}
		{/if}
	</div>

	{#if showFunFact}
		<div
			class="max-w-lg px-4 text-center text-base text-gray-600 select-text sm:px-8"
			aria-live="polite"
			aria-atomic="true"
			transition:fade={{ duration: 500 }}
		>
			{funFacts[currentFactIndex]}
		</div>
	{/if}
</div>

<style>
	@keyframes spin {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}
</style>
