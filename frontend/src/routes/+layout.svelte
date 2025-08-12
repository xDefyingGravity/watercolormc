<script lang="ts">
	import '../app.css'
	import Header from '$lib/components/Header.svelte'
	import { init } from '$lib/init'
	import LoadingScreen from '$lib/components/LoadingScreen.svelte'
	import { afterNavigate } from '$app/navigation'
	import { backendStatus, currentPage, pages } from '$lib/stores'

	let { children } = $props()

	let loaded = $state(false)

	const start = performance.now()

	init()
		.then(() => {
			const elapsed = performance.now() - start
			const remaining = 300 - elapsed

			setTimeout(
				() => {
					loaded = true
				},
				Math.max(0, remaining)
			)
		})
		.catch((error) => {
			console.error('Initialization failed:', error)
			loaded = true
		})

	afterNavigate(({ to }) => {
		const page = pages.find((x) => x.route === to?.route.id) || {
			route: to?.route.id || '',
			name: ''
		}
		currentPage.set(page as never)
	})
</script>

{#if $backendStatus === 'ok'}
	{#if loaded}
		<div class="bg-gradient-to-br">
			<Header />
			{@render children()}
		</div>
	{:else}
		<LoadingScreen />
	{/if}
{:else}
	<div class="flex min-h-screen items-center justify-center bg-gray-100 px-4">
		<div class="max-w-md text-center">
			<h1 class="mb-4 text-6xl font-bold text-red-500">503</h1>
			<h2 class="mb-2 text-2xl font-semibold text-gray-800">Server Unavailable</h2>
			<p class="mb-6 text-gray-600">
				The local server isn't responding. Make sure the program is running and try refreshing.
			</p>
		</div>
	</div>
{/if}
