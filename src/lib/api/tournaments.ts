import type { Entrant, Tournament } from '../../types/tournament';

export async function getEntrants(tournamentId: string): Promise<Entrant[]> {
	const response = await fetch(`/api/tournaments/${tournamentId}/entrants`);

	if (!response.ok) {
		throw new Error('Failed to load entrants.');
	}

	const data = await response.json();

	return data.entrants;
}

export async function getTournament(tournamentId: string) {
	const response = await fetch(`/api/tournaments/${tournamentId}`);

	const data = await response.json();

	if (!response.ok) {
		throw new Error(data.message ?? 'Failed to get tournament.');
	}

	return data.tournament;
}

export async function updateTournament(tournamentId: string, tournament: Partial<Tournament>) {
	const response = await fetch(`/api/tournaments/${tournamentId}`, {
		method: 'PATCH',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			qualifyingOrder: tournament.qualifyingOrder
		})
	});

	const data = await response.json();

	if (!response.ok) {
		throw new Error(data.message ?? 'Failed to update tournament.');
	}

	return data;
}

export async function deleteTournament(tournamentId: string) {
	const stateResponse = await fetch(`/api/tournaments/${tournamentId}/state`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			state: 'reset'
		})
	});

	const stateData = await stateResponse.json();
	if (!stateResponse.ok) {
		throw new Error(stateData.message ?? 'Failed to change tournament state.');
	}

	const response = await fetch(`/api/tournaments/${tournamentId}`, {
		method: 'DELETE'
	});

	const data = await response.json();
	if (!response.ok) {
		throw new Error(data.message ?? 'Failed to delete tournament.');
	}

	return data;
}

export async function startTournament(tournamentId: string) {
	const bracketResponse = await fetch(`/api/tournaments/${tournamentId}/bracket`);
	const bracketData = await bracketResponse.json();

	if (!bracketResponse.ok) {
		throw new Error(bracketData.message ?? 'Failed to get tournament bracket.');
	}

	if (bracketData.bracket.status === 'started') {
		console.info('tournament has already been started');
		return;
	}

	const response = await fetch(`/api/tournaments/${tournamentId}/state`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			state: 'start'
		})
	});

	const data = await response.json();
	if (!response.ok) {
		throw new Error(data.message ?? 'Failed to start tournament.');
	}

	return data;
}
