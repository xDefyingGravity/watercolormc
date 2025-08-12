/**
 * version.ts
 *
 * Description:
 *   frontend — version.ts module.
 *
 * Created: 7/11/25
 * Author: Will Ballantine
 *
 * @packageDocumentation
 * @copyright 2025-present Will Ballantine
 */

export function extractPaperVersion(url: string): string {
	const match = url.match(/paper-(\d+\.\d+\.\d+(?:-\d+)?)/)
	if (!match) throw new Error('invalid paper jar url')
	return `paper-${match[1]}`
}

export function getBaseVersion(version: string): string {
	return version.split('-')[1]
}
