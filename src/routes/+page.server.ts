import { getJson } from '../lib/utils.js';
import type { Tournament } from '../types/tournament.js';

export async function load({ fetch }) {
	const [{ tournaments }] = await Promise.all([
		getJson<{ tournaments: Tournament[] }>(fetch, '/api/tournaments')
	]);

	return { tournaments };
}
