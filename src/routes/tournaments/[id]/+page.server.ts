import { getEntrants, syncTournament } from '$lib/api/tournaments.js';
import { getJson } from '$lib/utils';
import type { CriteriaCatrgory } from '../../../types/index.js';
import type { GroupedMatches, Tournament } from '../../../types/tournament';

type CriteriaConfig = {
	name: string;
	categories: Record<string, CriteriaCatrgory>;
};

export async function load({ fetch, params }) {
	const id = params.id;

	const [tournament, entrants, criteriaConfig] = await Promise.all([
		getJson<Tournament>(fetch, `/api/tournaments/${id}`),
		getEntrants(fetch, id),
		getJson<CriteriaConfig>(fetch, '/conf/default.json')
	]);

	let matches: GroupedMatches = [];
	if (tournament.status === 'elimination') {
		matches = await syncTournament(fetch, id);
	}

	return {
		tournament,
		entrants,
		matches,
		criteriaConfig
	};
}
