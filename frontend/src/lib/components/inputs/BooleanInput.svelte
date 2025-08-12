<script lang="ts">
	import { fly } from 'svelte/transition';
	import { tick } from 'svelte';

	export let id: string;
	export let value: string = 'false';
	export let error: string = '';
	export let label: string = '';
	export let onBlur: () => void = () => {};

	let open = false;
	let query = '';
	let options = [
		{ label: 'True', value: 'true' },
		{ label: 'False', value: 'false' }
	];
	let filtered = options;
	let focusedIndex = -1;

	function toggle() {
		open = !open;
		if (open) {
			filtered = options;
			focusedIndex = options.findIndex(o => o.value === value);
		}
	}

	function close() {
		open = false;
		onBlur();
	}

	function select(option: { label: string; value: string }) {
		value = option.value;
		open = false;
	}

	function onKeyDown(e: KeyboardEvent) {
		if (!open) {
			if (e.key === 'ArrowDown' || e.key === 'Enter' || e.key === ' ') {
				e.preventDefault();
				toggle();
			}
			return;
		}

		switch (e.key) {
			case 'ArrowDown':
				e.preventDefault();
				focusedIndex = (focusedIndex + 1) % filtered.length;
				break;
			case 'ArrowUp':
				e.preventDefault();
				focusedIndex = (focusedIndex - 1 + filtered.length) % filtered.length;
				break;
			case 'Enter':
			case ' ':
				e.preventDefault();
				if (focusedIndex >= 0 && focusedIndex < filtered.length) {
					select(filtered[focusedIndex]);
				}
				break;
			case 'Escape':
				e.preventDefault();
				close();
				break;
		}
	}

	$: filtered = options.filter(o =>
		o.label.toLowerCase().includes(query.toLowerCase())
	);

	async function onInput(e: Event) {
		const target = e.target as HTMLInputElement;
		query = target.value;
		await tick();
		focusedIndex = 0;
	}

	function onFocus() {
		open = true;
		filtered = options;
		focusedIndex = options.findIndex(o => o.value === value);
	}

	function onBlurHandler() {
		setTimeout(() => {
			if (!document.activeElement?.closest(`#${id}-dropdown`)) {
				close();
			}
		}, 150);
	}
</script>

<div class="flex flex-col w-full" id={id}>
	<label
		for={`${id}-input`}
		class="mb-1 text-xs font-medium text-gray-500 select-none"
	>{label}</label
	>

	<div
		tabindex="0"
		id={`${id}-input`}
		class="relative flex w-full cursor-pointer items-center rounded border border-gray-200 bg-white px-3 py-2 text-base text-gray-900 transition focus-within:ring-2 focus-within:ring-blue-200 focus-within:outline-none {error ? 'border-red-400' : ''}"
		on:keydown={onKeyDown}
		on:focus={onFocus}
		on:blur={onBlurHandler}
		role="combobox"
		aria-expanded={open}
		aria-haspopup="listbox"
		aria-controls={`${id}-dropdown`}
	>
		<input
			class="w-full border-none p-0 m-0 text-base text-gray-900 bg-white focus:outline-none placeholder-gray-400 rounded"
			autocomplete="off"
			placeholder="Select..."
			value={query || options.find(o => o.value === value)?.label}
			on:input={onInput}
			readonly={false}
		/>
		<span class="material-symbols-outlined pointer-events-none select-none text-gray-300 ml-2 text-lg transition-transform duration-200 {open ? 'rotate-180' : ''}">
			arrow_drop_down
		</span>

		{#if open}
			<ul
				id={`${id}-dropdown`}
				class="absolute top-full left-0 z-20 mt-2 max-h-40 w-full overflow-auto rounded border border-gray-200 bg-white animate-fade-in"
				transition:fly={{ y: 2, duration: 80 }}
			>
				{#if filtered.length === 0}
					<li class="px-4 py-2 text-gray-400 text-sm">No options found</li>
				{:else}
					{#each filtered as option, i}
						<li
							class="cursor-pointer px-4 py-2 text-sm transition-colors duration-75 hover:bg-blue-50 {i === focusedIndex ? 'bg-blue-100 font-medium' : ''} rounded flex items-center justify-between"
						>
							<button class="w-full text-left bg-transparent border-none p-0 m-0 flex items-center justify-between" on:mousedown|preventDefault={() => select(option)}>
								{option.label}
								{#if value === option.value}
									<span class="material-symbols-outlined text-blue-400 text-base ml-2">check</span>
								{/if}
							</button>
						</li>
					{/each}
				{/if}
			</ul>
		{/if}
	</div>
	{#if error}
		<p class="mt-1 text-sm text-red-500">{error}</p>
	{/if}
</div>

<style>
    input:focus {
        outline: none;
    }
</style>
