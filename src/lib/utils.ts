export { cn } from 'cn';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, 'child'> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any } ? Omit<T, 'children'> : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & { ref?: U | null };

export async function getJson<T>(
	fetch: typeof globalThis.fetch,
	url: string,
	options?: RequestInit
): Promise<T> {
	const response = await fetch(url, options);

	let data: T | { message?: string } | null = null;

	try {
		data = await response.json();
	} catch {
		// Non-JSON response
	}

	if (!response.ok) {
		const message =
			data && typeof data === 'object' && 'message' in data ? data.message : undefined;

		throw new Error(message ?? `Request failed: ${response.status}`);
	}

	return data as T;
}
