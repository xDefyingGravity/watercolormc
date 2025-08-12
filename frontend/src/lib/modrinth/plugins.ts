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

import { client } from '$lib/modrinth/client'
import type { SearchProjectOptions, SearchResultHit } from '@xmcl/modrinth'

export interface ProjectDependency {
	version_id: string | null
	project_id: string
	dependency_type: 'required' | 'optional' | 'incompatible' | 'embedded'
}

export interface PluginInfo {
	id: string
	name: string
	description: string
	icon: string
	downloads: number
}

export interface TopPluginOptions {
	limit?: number
	offset?: number
	query?: string
}

export async function getTopPlugins(opts: TopPluginOptions): Promise<PluginInfo[]> {
	const options: SearchProjectOptions = {
		query: opts.query ?? '',
		limit: opts.limit ?? 10,
		offset: opts.offset ?? 0,
		index: 'downloads',
		facets: JSON.stringify([['project_type:plugin']])
	}

	const result = await client.searchProjects(options)

	return result.hits.map(
		(project: SearchResultHit): PluginInfo => ({
			id: project.project_id,
			name: project.title,
			description: project.description,
			icon: project.icon_url,
			downloads: project.downloads
		})
	)
}

export async function getPlugin(id: string): Promise<ReturnType<typeof client.getProject>> {
	return client.getProject(id)
}

export type PluginDownloadInfo = {
	id: string;
	downloadUrl: string;
};

export async function getProjectDependencyGraph(
	root: string,
	mcVersion: string
): Promise<PluginDownloadInfo[]> {
	const visited = new Map<string, PluginDownloadInfo>();

	async function crawl(slug: string) {
		if (visited.has(slug)) return;

		const versions = await client.getProjectVersions(slug, {
			gameVersions: [mcVersion],
		});
		if (!versions.length) {
			throw new Error(`No compatible versions found for ${slug} on MC ${mcVersion}`);
		}

		const latest = versions[0];
		const primaryFile = latest.files.find(f => f.primary);
		if (!primaryFile) {
			throw new Error(`No primary file found for ${slug} version ${latest.version_number}`);
		}

		visited.set(slug, {
			id: slug,
			downloadUrl: primaryFile.url,
		});

		const deps = latest.dependencies.filter(dep => dep.dependency_type === 'required');
		for (const dep of deps) {
			await crawl(dep.project_id);
		}
	}

	await crawl(root);

	return Array.from(visited.values());
}

export type { Project } from '@xmcl/modrinth'
