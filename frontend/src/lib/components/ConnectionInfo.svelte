<script lang="ts">
	export let host: string
	export let port: number
	export let version: string
	export let privateIp: string

	function copyAddress(text: string) {
		navigator.clipboard.writeText(text)
	}
</script>

<section class="mb-8" transition:slide={{ duration: 350, delay: 80 }}>
	<h2 class="mb-2 text-xl">Connection Info</h2>
	<div class="flex items-center gap-2">
		<code class="rounded bg-gray-100 px-3 py-1 font-mono text-base">
			{host === '0.0.0.0' ? privateIp : host}:{port}
		</code>
		<button
			class="ml-1 rounded p-1 text-gray-400 transition hover:bg-blue-50 hover:text-blue-600"
			on:click={() => copyAddress(`${host === '0.0.0.0' ? privateIp : host}:${port}`)}
			title="Copy address"
			type="button"
			aria-label="Copy server address"
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
	<div class="mt-1 text-sm text-gray-500">Version: {version}</div>
</section>
