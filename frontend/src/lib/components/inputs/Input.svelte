<script lang="ts">
	export let id: string;
	export let value: string = '';
	export let placeholder: string = '';
	export let error: string = '';
	export let onInput: (e: Event) => void = () => {};
	export let onBlur: () => void = () => {};
</script>

<div class="relative w-full">
	<input
		id={id}
		type="text"
		bind:value={value}
		placeholder={placeholder}
		on:input={onInput}
		on:blur={onBlur}
		class="w-full rounded border border-gray-200 px-3 py-2 text-base transition focus:ring-2 focus:ring-blue-200 focus:outline-none {error ? 'border-red-400' : ''}"
	/>
	{#if error}
		<p class="mt-1 text-sm text-red-500">{error}</p>
	{/if}
</div>
