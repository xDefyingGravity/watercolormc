/**
 * plugins.ts
 *
 * Description:
 *   frontend — plugins.ts module.
 *
 * Created: 7/11/25
 * Author: Will Ballantine
 *
 * @packageDocumentation
 * @copyright 2025-present Will Ballantine
 */
import { safeFetch } from '$lib/utils/fetch'
import { baseUrl } from '$lib/config'
import { getProjectDependencyGraph, type PluginDownloadInfo } from '$lib/modrinth/plugins'
import { addPluginToManifest, type ManifestPlugin, removePluginFromManifest } from '$lib/paper/manifest'
import { debug } from '$lib/logging'

export async function addPlugins(serverId: string, plugins: PluginDownloadInfo[]): Promise<void> {
	const response = await safeFetch<string>(`${baseUrl}/api/servers/${serverId}/plugins`, {
		method: 'POST',
		body: JSON.stringify({
			plugins: plugins.map(plugin => plugin.downloadUrl),
		}),
		headers: {
			'Content-Type': 'application/json'
		}
	})

	if (response !== 'ok') {
		throw new Error(`Failed to add plugins: ${response}`)
	}

	for (const plugin of plugins) {
		await addPluginToManifest(serverId, {
			id: plugin.id,
			jar_name: plugin.downloadUrl.split('/').pop() ?? ''
		})
	}
}

export async function removePlugin(serverId: string, plugin: string): Promise<void> {
	debug(`Removing plugin ${plugin} from server ${serverId}`)

	const encodedPlugin = encodeURIComponent(plugin)
	const response = await safeFetch<string>(
		`${baseUrl}/api/servers/${serverId}/plugins/${encodedPlugin}`,
		{
			method: 'DELETE'
		}
	)

	if (response !== 'ok') {
		throw new Error(`Failed to remove plugins: ${response}`)
	}

	await removePluginFromManifest(serverId, plugin)
}

export async function getPlugins(serverId: string): Promise<ManifestPlugin[]> {
	const response = await safeFetch<ManifestPlugin[]>(`${baseUrl}/api/servers/${serverId}/plugins/manifest`, {
		method: 'GET'
	})

	if (!Array.isArray(response)) {
		throw new Error(`Failed to fetch plugins: ${response}`)
	}

	return response
}

export async function addPluginGraph(serverId: string, plugin: string, mcVersion: string): Promise<void> {
	const dependencyGraph = await getProjectDependencyGraph(plugin, mcVersion)
	if (!dependencyGraph) {
		throw new Error(`Failed to fetch dependency graph for plugin: ${plugin}`)
	}

	await addPlugins(serverId, dependencyGraph)
}