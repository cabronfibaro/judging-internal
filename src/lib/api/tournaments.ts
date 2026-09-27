export async function updateTournament(tournamentId: string, qualifyingOrder: number[]) {
	const response = await fetch(`/api/tournaments/${tournamentId}`, {
		method: 'PATCH',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			qualifyingOrder: qualifyingOrder.join(', ')
		})
	});

	const data = await response.json();

	if (!response.ok) {
		throw new Error(data.message ?? 'Failed to update tournament.');
	}

	return data;
}
