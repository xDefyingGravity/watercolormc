<script lang="ts">
	import { fly } from 'svelte/transition'

	export let id: string
	export let value: string = 'normal'
	export let error: string = ''
	export let label: string = ''
	export let onBlur: () => void = () => {
	}

	const difficulties = [
		{ label: 'Peaceful', value: 'peaceful' },
		{ label: 'Easy', value: 'easy' },
		{ label: 'Normal', value: 'normal' },
		{ label: 'Hard', value: 'hard' }
	]
	let open = false
	let focusedIndex = -1

	function toggle() {
		open = !open
		if (open) focusedIndex = difficulties.findIndex(d => d.value === value)
	}

	function close() {
		open = false
		onBlur()
	}

	function select(option: { label: string; value: string }) {
		value = option.value
		open = false
	}

	function onKeyDown(e: KeyboardEvent) {
		if (!open) {
			if (['ArrowDown', 'Enter', ' '].includes(e.key)) {
				e.preventDefault()
				toggle()
			}
			return
		}
		switch (e.key) {
			case 'ArrowDown':
				focusedIndex = (focusedIndex + 1) % difficulties.length
				e.preventDefault()
				break
			case 'ArrowUp':
				focusedIndex = (focusedIndex - 1 + difficulties.length) % difficulties.length
				e.preventDefault()
				break
			case 'Enter':
			case ' ':
				if (focusedIndex >= 0) select(difficulties[focusedIndex])
				e.preventDefault()
				break
			case 'Escape':
				close()
				e.preventDefault()
				break
		}
	}

	function onFocus() {
		open = true
		focusedIndex = difficulties.findIndex(d => d.value === value)
	}

	function onBlurHandler() {
		setTimeout(() => {
			close()
		}, 150)
	}
</script>
<div class="flex flex-col w-full">
	<label for={id} class="mb-1 text-xs font-medium text-gray-500 select-none">{label}</label>
	<div
		tabindex="0"
		id={id}
		class="relative flex w-full cursor-pointer items-center rounded border border-gray-200 bg-white px-3 py-2 text-base text-gray-900 transition focus-within:ring-2 focus-within:ring-blue-200 focus-within:outline-none {error ? 'border-red-400' : ''}"
		on:keydown={onKeyDown}
		on:focus={onFocus}
		on:blur={onBlurHandler}
		role="combobox"
		aria-expanded={open}
		aria-haspopup="listbox"
		aria-controls={`${id}-dropdown`}
	>
		<span
			class="w-full text-left text-base text-gray-900 select-none">{difficulties.find(d => d.value === value)?.label}</span>
		<span
			class="material-symbols-outlined pointer-events-none select-none text-gray-300 ml-2 text-lg transition-transform duration-200 {open ? 'rotate-180' : ''}">arrow_drop_down</span>
		{#if open}
			<ul
				id={`${id}-dropdown`}
				class="absolute top-full left-0 z-20 mt-2 max-h-40 w-full overflow-auto rounded border border-gray-200 bg-white animate-fade-in"
				transition:fly={{ y: 2, duration: 80 }}
			>
				{#each difficulties as option, i}
					<li
						class="cursor-pointer px-4 py-2 text-sm transition-colors duration-75 hover:bg-blue-50 {i === focusedIndex ? 'bg-blue-100 font-medium' : ''} rounded flex items-center justify-between">
						<button class="w-full text-left bg-transparent border-none p-0 m-0 flex items-center justify-between"
										on:mousedown|preventDefault={() => select(option)}>
							{option.label}
							{#if value === option.value}
								<span class="material-symbols-outlined text-blue-400 text-base ml-2">check</span>
							{/if}
						</button>
					</li>
				{/each}
			</ul>
		{/if}
	</div>
	{#if error}
		<p class="mt-1 text-sm text-red-500">{error}</p>
	{/if}
</div>
