import { debug, error } from '$lib/logging'

export async function safeFetch<T = unknown>(
	url: string,
	options?: RequestInit
): Promise<T | undefined> {
	try {
		const response = await fetch(url, options)

		if (!response.ok) {
			error(`fetch failed: ${response.status} ${response.statusText}`)
			return undefined
		}

		debug(
			`fetch successful: ${response.status} ${response.statusText}\n` +
				`response url: ${response.url}\n` +
				'response headers:\n' +
				JSON.stringify(Object.fromEntries(response.headers.entries()), null, 2)
		)
		const contentType = response.headers.get('content-type') || ''

		if (contentType.includes('application/json')) {
			return (await response.json()) as T
		}

		return (await response.text()) as T
	} catch (e: unknown) {
		error('fetch failed:', e instanceof Error ? e.message : String(e))
		return undefined
	}
}
