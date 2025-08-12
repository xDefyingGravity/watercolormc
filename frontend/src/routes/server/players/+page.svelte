<script lang="ts">
	import ServerSidebar from '$lib/components/ServerSidebar.svelte'
	import { onDestroy, onMount } from 'svelte'
	import { connectToPlayersSocket } from '$lib/utils/serverSockets.js'
	import { getUUIDFromPlayerName } from '$lib/uuid'
	import { baseUrl } from '$lib/config'
	import { safeFetch } from '$lib/utils/fetch'

	const serverId = localStorage.getItem('selectedServer')

	interface Player {
		name: string
		uuid: string
	}

	let players: Player[] = []

	async function onPlayerJoin(player: string) {
		players = [
			...players,
			{
				name: player,
				uuid: await getUUIDFromPlayerName(player)
			}
		]
	}
	function onPlayerLeave(player: string) {
		players = players.filter((p) => p.name !== player)
	}

	let socket: ReturnType<typeof connectToPlayersSocket> | null = null

	onMount(async () => {
		if (!serverId) {
			window.location.href = '/servers'
		}

		socket = connectToPlayersSocket(serverId!, {
			onPlayerJoin,
			onPlayerLeave
		})

		const response = await safeFetch<string[] | string>(
			`${baseUrl}/api/servers/${serverId}/players`
		)

		if (Array.isArray(response)) {
			players = await Promise.all(
				response.map(async (player) => ({
					name: player,
					uuid: await getUUIDFromPlayerName(player)
				}))
			)
		} else {
			console.error('Failed to fetch players:', response)
		}
	})

	onDestroy(() => {
		if (socket) {
			socket.close()
		}
	})
</script>

<ServerSidebar />

<div class="flex min-h-screen w-full flex-col items-center justify-center bg-white px-0">
	<div
		class="animate-fade-in w-full max-w-4xl rounded-xl border border-gray-200 bg-white p-10 shadow-lg"
	>
		<h2
			class="mb-8 flex items-center justify-center gap-3 text-center text-3xl font-semibold text-gray-900"
		>
			<span class="material-symbols-outlined text-3xl text-blue-600">group</span>
			Players Online
		</h2>
		<p class="mb-8 text-center text-lg text-gray-700">
			Current players connected to the server <span class="font-semibold text-blue-600"
				>{serverId}</span
			>.
		</p>

		{#if players.length > 0}
			<ul class="grid gap-6 sm:grid-cols-2">
				{#each players as player}
					<li
						class="animate-fade-in flex items-center gap-4 rounded-lg border border-blue-100 bg-blue-50 px-6 py-4 shadow-sm transition hover:shadow-md"
					>
						<img
							class="h-12 w-12 rounded-full bg-blue-200"
							src={`https://crafatar.com/avatars/${player.uuid}?size=64&default=MHF_Steve`}
							alt={player.name}
						/>
						<span class="truncate text-lg font-medium text-gray-900">{player.name}</span>
					</li>
				{/each}
			</ul>
		{:else}
			<div class="animate-fade-in flex flex-col items-center justify-center py-12">
				<span class="material-symbols-outlined mb-4 text-5xl text-gray-300">sentiment_neutral</span>
				<p class="text-lg text-gray-500">No players are currently online.</p>
			</div>
		{/if}
	</div>
</div>

<style>
	@keyframes fade-in {
		from {
			opacity: 0;
			transform: translateY(10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
	.animate-fade-in {
		animation: fade-in 0.5s ease;
	}
</style>
