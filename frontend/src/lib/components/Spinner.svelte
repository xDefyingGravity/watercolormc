<script lang="ts">
	export let size: number = 16
	export let color: string = 'currentColor'
</script>

<span
	role="status"
	aria-label="Loading"
	class="inline-block animate-spin rounded-full border-2 border-t-transparent"
	style="
		width: {size}px;
		height: {size}px;
		border-color: {color};
		border-top-color: transparent;
	"
>
	<!-- No Content -->
</span>
