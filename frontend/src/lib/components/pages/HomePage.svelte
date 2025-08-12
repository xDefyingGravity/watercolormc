<script>
	import FloatingActionButton from '$lib/components/utilities/FloatingActionButton.svelte'
	import { goto } from '$app/navigation'
	import ServerList from '$lib/components/ServerList.svelte'

	import { startTour } from '$lib/tutorial'
	import { onMount } from 'svelte'

	onMount(() => {
		if (localStorage.getItem('tutorial') === null) {
			startTour();
		}
	});
</script>

<ServerList />

<FloatingActionButton onClick={() => goto('/create-server')} description="Add new server" id="add-button"/>
