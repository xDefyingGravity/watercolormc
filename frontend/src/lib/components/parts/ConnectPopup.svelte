<script lang="ts">
	import { fade, scale } from 'svelte/transition'
	import { ipinfo } from '$lib/stores'
	export let port: number
	export let serverName: string
	export let host: string
	export let onClose: () => void

	function copy(text: string) {
		navigator.clipboard.writeText(text)
	}
</script>

<div
	class="fixed inset-0 z-50 flex items-center justify-center bg-black/20 p-6"
	transition:fade={{ duration: 200 }}
>
	<div
		class="w-full max-w-lg rounded-lg bg-white p-6 shadow-lg"
		transition:scale={{ duration: 250, start: 0.96 }}
	>
		<h2 class="mb-2 text-xl font-semibold text-gray-900">
			Connect to <span class="text-blue-600">{serverName}</span>
		</h2>
		<p class="mb-5 text-sm text-gray-600">Choose the option that fits your situation:</p>

		<ol class="space-y-4">
			<li class="group flex items-start gap-3">
				<div
					class="mt-1 flex h-5 w-5 flex-shrink-0 items-center justify-center rounded-full bg-blue-50 text-xs font-bold text-blue-500"
				>
					1
				</div>
				<div class="flex-1">
					<div class="font-medium text-gray-800">For me</div>
					<div class="mt-1 flex items-center gap-1">
						<code class="rounded bg-gray-100 px-2 py-0.5 font-mono text-sm"
							>{`localhost:${port}`}</code
						>
						<button
							class="ml-1 rounded p-1 text-gray-400 opacity-0 transition group-hover:opacity-100 hover:bg-blue-50 hover:text-blue-600"
							on:click={() => copy(`localhost:${port}`)}
							title="Copy"
							type="button"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-4 w-4"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								stroke-width="2"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2M16 8h2a2 2 0 012 2v8a2 2 0 01-2 2h-8a2 2 0 01-2-2v-2"
								/>
							</svg>
						</button>
					</div>
				</div>
			</li>
			<li class="group flex items-start gap-3 border-t border-gray-100 pt-4">
				<div
					class="mt-1 flex h-5 w-5 flex-shrink-0 items-center justify-center rounded-full bg-blue-50 text-xs font-bold text-blue-500"
				>
					2
				</div>
				<div class="flex-1">
					<div class="font-medium text-gray-800">For people in my house</div>
					<div class="mt-1 flex items-center gap-1">
						<code class="rounded bg-gray-100 px-2 py-0.5 font-mono text-sm"
							>{`${host === '0.0.0.0' ? $ipinfo.privateIp : host}:${port}`}</code
						>
						<button
							class="ml-1 rounded p-1 text-gray-400 opacity-0 transition group-hover:opacity-100 hover:bg-blue-50 hover:text-blue-600"
							on:click={() => copy(`${host === '0.0.0.0' ? $ipinfo.privateIp : host}:${port}`)}
							title="Copy"
							type="button"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-4 w-4"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								stroke-width="2"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2M16 8h2a2 2 0 012 2v8a2 2 0 01-2 2h-8a2 2 0 01-2-2v-2"
								/>
							</svg>
						</button>
					</div>
				</div>
			</li>
			<li class="group flex items-start gap-3 border-t border-gray-100 pt-4">
				<div
					class="mt-1 flex h-5 w-5 flex-shrink-0 items-center justify-center rounded-full bg-blue-50 text-xs font-bold text-blue-500"
				>
					3
				</div>
				<div class="flex-1">
					<div class="font-medium text-gray-800">For friends not in my house</div>
					<div class="mb-1 text-xs text-red-400">
						Requires port forwarding or VPN setup on your router.
					</div>
					<div class="flex items-center gap-1">
						<code class="rounded bg-gray-100 px-2 py-0.5 font-mono text-sm"
							>{`${$ipinfo.publicIp}:${port}`}</code
						>
						<button
							class="ml-1 rounded p-1 text-gray-400 opacity-0 transition group-hover:opacity-100 hover:bg-blue-50 hover:text-blue-600"
							on:click={() => copy(`${$ipinfo.publicIp}:${port}`)}
							title="Copy"
							type="button"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-4 w-4"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								stroke-width="2"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2M16 8h2a2 2 0 012 2v8a2 2 0 01-2 2h-8a2 2 0 01-2-2v-2"
								/>
							</svg>
						</button>
					</div>
				</div>
			</li>
		</ol>

		<button
			class="mt-6 w-full rounded bg-blue-600 py-2 font-semibold text-white hover:bg-blue-700 focus:ring-2 focus:ring-blue-400 focus:outline-none"
			on:click={onClose}
		>
			Got it
		</button>
	</div>
</div>
