import { getJson } from '#lib/utils';
import type { Entrant, GroupedMatches, Tournament } from '../../types/tournament';

export async function getEntrants(
	fetch: typeof globalThis.fetch,
	tournamentId: string
): Promise<Entrant[]> {
	const { entrants } = await getJson<{ entrants: Entrant[] }>(
		fetch,
		`/api/tournaments/${tournamentId}/entrants`
	);

	return entrants;
}

export async function getTournaments() {
	const { tournaments } = await getJson<{ tournaments: Tournament[] }>(fetch, `/api/tournaments`);
	return tournaments;
}

export async function getTournament(tournamentId: string) {
	const { tournament } = await getJson<{ tournament: Tournament }>(
		fetch,
		`/api/tournaments/${tournamentId}`
	);

	return tournament;
}

export async function updateTournament(tournamentId: string, tournament: Partial<Tournament>) {
	return getJson(fetch, `/api/tournaments/${tournamentId}`, {
		method: 'PATCH',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			qualifyingOrder: tournament.qualifyingOrder
		})
	});
}

export async function deleteTournament(tournamentId: string, status: Tournament['status']) {
	if (status !== 'qualifying') {
		await getJson(fetch, `/api/tournaments/${tournamentId}/state`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				state: 'reset'
			})
		});
	}

	return getJson(fetch, `/api/tournaments/${tournamentId}`, {
		method: 'DELETE'
	});
}

export async function getMatches(tournamentId: string) {
	const { matches } = await getJson<{ matches: GroupedMatches }>(
		fetch,
		`/api/tournaments/${tournamentId}/bracket/matches`
	);

	return matches;
}

export async function startTournament(tournamentId: string) {
	await getJson(fetch, `/api/tournaments/${tournamentId}/state`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			state: 'start'
		})
	});
}

export type Battles = {
	id: string;
	matchId: string;
	participantId: string;
	leadStatus: number;
}[];

export async function syncTournament(fetch: typeof globalThis.fetch, tournamentId: string) {
	const { matches, battles } = await getJson<{ matches: GroupedMatches; battles: Battles }>(
		fetch,
		`/api/tournaments/${tournamentId}/bracket/sync`
	);

	return { matches, battles };
}

export async function getActiveBattles(tournamentId: string) {
	const battles = await getJson<Battles>(
		fetch,
		`/api/tournaments/${tournamentId}/bracket/activeBattles`
	);

	return battles;
}
