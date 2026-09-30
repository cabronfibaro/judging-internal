import { CHALLONGE_API_KEY } from '$env/static/private';
import { db } from '$lib/server/db/client.js';
import { brackets } from '$lib/server/db/schema.js';
import { json } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';

const states = [
	'process_checkin',
	'abort_checkin',
	'start_group_stage',
	'finalize_group_stage',
	'reset_group_stage',
	'start',
	'finalize',
	'reset',
	'open_predictions'
];

export async function POST({ params, request }) {
	const tournamentId = params.id;
	const body = await request.json();

	const { state } = body;

	if (!tournamentId) {
		return json({ message: 'Tournament ID is required.' }, { status: 400 });
	}

	if (!states.includes(state)) {
		return json({ message: `State '${state}' is invalid` });
	}

	const bracket = db
		.select()
		.from(brackets)
		.where(eq(brackets.tournamentId, tournamentId))
		.limit(1)
		.all()[0];

	if (!bracket) {
		return json({ message: 'Bracket not found.' }, { status: 404 });
	}

	try {
		const challongeResponse = await fetch(
			`https://api.challonge.com/v2.1/tournaments/${bracket.challongeId}/change_state.json`,
			{
				method: 'PUT',
				headers: {
					Accept: 'application/json',
					'Content-Type': 'application/vnd.api+json',
					'Authorization-Type': 'v1',
					Authorization: CHALLONGE_API_KEY
				},
				body: JSON.stringify({
					data: {
						type: 'TournamentState',
						attributes: {
							state
						}
					}
				})
			}
		);

		const challongeData = await challongeResponse.json();

		if (!challongeResponse.ok) {
			console.error(
				'Failed to change the tournament state',
				JSON.stringify(challongeData, null, 2)
			);

			return json(
				{
					message: 'Failed to change the tournament state',
					error: challongeData
				},
				{ status: challongeResponse.status }
			);
		}

		// update bracket status to start as well
		if (state === 'start') {
			const [updated] = db
				.update(brackets)
				.set({ status: 'started' })
				.where(eq(brackets.tournamentId, tournamentId))
				.returning({ status: brackets.status })
				.all();

			if (!updated) {
				return json({ message: 'Cannot update bracket status to started.' }, { status: 404 });
			}
		}

		return json({ message: 'OK' });
	} catch (error) {
		console.error('Failed to start the tournament:', error);
		return json({ message: 'Failed to start the tournament.' }, { status: 500 });
	}
}
