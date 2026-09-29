/* eslint-disable @typescript-eslint/no-explicit-any */
import { db } from '$lib/server/db/client';
import { brackets, nodes, participants } from '$lib/server/db/schema';
import { CHALLONGE_API_KEY } from '$env/static/private';
import { json } from '@sveltejs/kit';
import { and, eq, inArray } from 'drizzle-orm';

export async function GET({ params }) {
	const tournamentId = params.id;

	if (!tournamentId) {
		return json({ message: 'Tournament ID is required.' }, { status: 400 });
	}

	try {
		const bracket = db
			.select({
				id: brackets.id,
				challongeTournamentId: brackets.url
			})
			.from(brackets)
			.where(eq(brackets.tournamentId, tournamentId))
			.limit(1)
			.all()[0];

		if (!bracket) {
			return json({ message: 'Bracket not found.' }, { status: 404 });
		}

		const response = await fetch(
			`https://api.challonge.com/v2.1/tournaments/${bracket.challongeTournamentId}/matches.json`,
			{
				method: 'GET',
				headers: {
					'Content-Type': 'application/vnd.api+json',
					Accept: 'application/json',
					'Authorization-Type': 'v1',
					Authorization: CHALLONGE_API_KEY
				}
			}
		);

		const data = await response.json();

		if (!response.ok) {
			console.error('Failed to fetch Challonge matches:', JSON.stringify(data, null, 2));

			return json(
				{
					message: 'Failed to fetch Challonge matches.',
					error: data
				},
				{ status: response.status }
			);
		}

		const matches = data.data;

		const challongeParticipantIds = [
			...new Set(
				matches.flatMap((match: any) =>
					match.attributes.points_by_participant
						.map((entry: any) => String(entry.participant_id))
						.filter(Boolean)
				)
			)
		] as string[];

		const localParticipants = db
			.select({
				id: participants.id,
				challongeParticipantId: participants.challongeParticipantId
			})
			.from(participants)
			.where(
				and(
					eq(participants.tournamentId, tournamentId),
					inArray(participants.challongeParticipantId, challongeParticipantIds)
				)
			)
			.all();

		const participantMap = new Map(
			localParticipants.map((participant) => [participant.challongeParticipantId, participant.id])
		);

		db.transaction((tx) => {
			for (const match of matches) {
				const { state, round, points_by_participant } = match.attributes;

				for (const entry of points_by_participant) {
					const participantId = entry.participant_id
						? participantMap.get(String(entry.participant_id))
						: null;

					if (entry.participant_id && !participantId) {
						throw new Error(
							`Local participant not found for Challonge participant ${entry.participant_id}`
						);
					}

					tx.insert(nodes)
						.values({
							bracketId: bracket.id,
							participantId,
							matchId: match.id,
							status: state === 'open' ? 'open' : 'pending',
							round
						})
						.run();
				}
			}
		});

		return json({
			message: 'Bracket synchronized.'
		});
	} catch (error) {
		console.error('Failed to sync bracket:', error);

		return json({ message: 'Failed to sync bracket.' }, { status: 500 });
	}
}
