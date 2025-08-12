/**
 * index.ts
 *
 * Description:
 *   frontend — index.ts module.
 *
 * Created: 7/10/25
 * Author: Will Ballantine
 *
 * @packageDocumentation
 * @copyright 2025-present Will Ballantine
 */
import { safeFetch } from '$lib/utils/fetch'
import { baseUrl } from '$lib/config'
import { error } from '$lib/logging'

export async function getUUIDFromPlayerName(playerName: string): Promise<string> {
	const response = await safeFetch<string>(baseUrl + `/api/minecraft/uuid/${playerName}`)

	if (response === undefined) {
		throw new Error(`Error fetching UUID for player: ${playerName}`)
	}

	try {
		const data = JSON.parse(response as string)
		return data.id
	} catch (e) {
		error(`Error parsing UUID response: ${e}`)
		return response
	}
}
