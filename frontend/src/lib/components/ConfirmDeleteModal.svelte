<script lang="ts">
	export let serverName: string
	export let onConfirm: () => Promise<void>
	export let onCancel: () => void

	let input = ''
	let error = ''
	let loading = false

	async function handleConfirm() {
		if (input === serverName) {
			loading = true
			error = ''
			try {
				await onConfirm()
			} catch (e) {
				error = 'Failed to delete. Please try again.'
			}
			loading = false
		} else {
			error = 'Server name does not match.'
		}
	}
</script>

<div class="animate-fadeIn fixed inset-0 z-50 flex items-center justify-center bg-black/40">
	<div class="animate-scaleIn relative w-full max-w-md rounded-xl bg-white p-8 shadow-2xl">
		<button
			class="absolute top-4 right-4 text-gray-400 transition hover:text-gray-600"
			aria-label="Close"
			on:click={onCancel}
			disabled={loading}
			type="button"
		>
			<svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
			</svg>
		</button>
		<div class="flex flex-col items-center gap-2">
			<div class="mb-2 flex h-12 w-12 items-center justify-center rounded-full bg-red-50">
				<span class="material-symbols-outlined text-3xl leading-none text-red-500"> delete </span>
			</div>
			<h2 class="mb-1 text-center text-xl font-bold text-red-600">Delete Server</h2>
			<p class="mb-2 text-center text-gray-700">
				This will <span class="font-bold text-red-600">immediately erase</span> the entire server
				<span class="comfortaa-nunito">{serverName}</span> and all its data.<br />
				<span class="text-gray-500">This action cannot be undone.</span>
			</p>
			<label class="mb-1 w-full text-sm font-medium text-gray-600" for="confirm-input">
				Type the server name to confirm:
			</label>
			<input
				id="confirm-input"
				class="mb-1 w-full rounded border border-gray-300 px-3 py-2 transition focus:ring-2 focus:ring-red-200 focus:outline-none"
				bind:value={input}
				placeholder="Enter server name"
				disabled={loading}
				autocomplete="off"
				autofocus
			/>
			{#if error}
				<div class="mb-1 w-full text-left text-sm text-red-500">{error}</div>
			{/if}
			<div class="mt-2 flex w-full justify-end gap-2">
				<button
					class="rounded bg-gray-100 px-4 py-2 text-gray-700 transition hover:bg-gray-200"
					on:click={onCancel}
					disabled={loading}
					type="button"
				>
					Cancel
				</button>
				<button
					class="flex items-center gap-2 rounded bg-red-600 px-4 py-2 text-white transition hover:bg-red-700 disabled:opacity-60"
					on:click={handleConfirm}
					disabled={loading || input !== serverName}
					type="button"
				>
					{#if loading}
						<svg class="h-4 w-4 animate-spin" viewBox="0 0 24 24">
							<circle
								class="opacity-25"
								cx="12"
								cy="12"
								r="10"
								stroke="white"
								stroke-width="4"
								fill="none"
							/>
							<path class="opacity-75" fill="white" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
						</svg>
					{/if}
					Delete
				</button>
			</div>
		</div>
	</div>
</div>

<style>
	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}
	@keyframes scaleIn {
		from {
			transform: scale(0.95);
			opacity: 0;
		}
		to {
			transform: scale(1);
			opacity: 1;
		}
	}
	.animate-fadeIn {
		animation: fadeIn 0.2s;
	}
	.animate-scaleIn {
		animation: scaleIn 0.2s;
	}
	.animate-spin {
		animation: spin 0.7s linear infinite;
	}
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
