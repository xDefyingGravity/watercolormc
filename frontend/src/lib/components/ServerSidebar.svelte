<script lang="ts">
	import { goto } from '$app/navigation'
	import { servers } from '$lib/stores'

	const buttons = [
		{ aria: 'Home', icon: 'home', path: '/server' },
		{ aria: 'World', icon: 'globe', path: '/server/world' },
		{ aria: 'Properties', icon: 'tune', path: '/server/properties' },
		{ aria: 'Players', icon: 'group', path: '/server/players' },
		{ aria: 'Backups', icon: 'save_alt', path: '/server/backups' },
		{ aria: 'Settings', icon: 'settings', path: '/server/settings' }
	]

	const serverId = localStorage.getItem('selectedServer')
	if (!serverId) {
		goto('/servers')
	}

	const server = $servers.find((s) => s.id === serverId)
	if (server?.type === 'papermc') {
		buttons.splice(4, 0, { aria: 'Plugins', icon: 'extension', path: '/server/plugins' })
	}

	async function navigate(path: string) {
		await goto(path)
	}
</script>

<div
	class="absolute top-12 left-0 z-10 flex h-full w-16 flex-col items-center bg-white px-4 shadow-md"
>
	{#each buttons as { aria, icon, path }}
		<button
			class="mt-3 flex h-10 w-10 items-center justify-center rounded-full p-2
             text-gray-500/75 transition-colors
             duration-150 hover:bg-blue-100 hover:text-blue-600
             focus:ring-0 focus:outline-none active:ring-0 active:outline-none"
			type="button"
			aria-label={aria}
			onclick={() => navigate(path)}
		>
			<span class="material-symbols-outlined text-[18px] leading-none select-none">{icon}</span>
		</button>
	{/each}
</div>
