import { db } from '#lib/server/db/client';
import { entrants, tournaments, participants, brackets } from '#lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { CHALLONGE_API_KEY } from '$app/env/private';
import { and, eq, inArray } from 'drizzle-orm';

export async function GET({ params }) {
	const tournamentId = params.id;

	if (!tournamentId) {
		return json({ message: 'Tournament ID is required.' }, { status: 400 });
	}

	try {
		const tournament = db
			.select({
				size: tournaments.size,
				challongeTournamentId: brackets.url
			})
			.from(tournaments)
			.leftJoin(brackets, eq(tournaments.id, brackets.tournamentId))
			.where(eq(tournaments.id, tournamentId))
			.limit(1)
			.all()[0];

		if (!tournament) {
			return json({ message: 'Tournament not found.' }, { status: 404 });
		}

		const tournamentEntrants = db
			.select({
				id: entrants.id,
				tournamentId: entrants.tournamentId,
				participantId: entrants.participantId,
				status: entrants.status,
				score: entrants.score,
				name: participants.name,
				challongeParticipantId: participants.challongeParticipantId
			})
			.from(entrants)
			.leftJoin(participants, eq(entrants.participantId, participants.id))
			.where(eq(entrants.tournamentId, tournamentId))
			.all();

		const sorted = tournamentEntrants.sort((a, b) => b.score - a.score);
		const qualifiedEntrants = sorted.slice(0, tournament.size);
		const qualifiedIds = qualifiedEntrants.map((entrant) => entrant.id);

		const challongeResponse = await fetch(
			`https://api.challonge.com/v2.1/tournaments/${tournament.challongeTournamentId}/participants/bulk_add.json`,
			{
				method: 'POST',
				headers: {
					Accept: 'application/json',
					'Content-Type': 'application/vnd.api+json',
					'Authorization-Type': 'v1',
					Authorization: CHALLONGE_API_KEY
				},
				body: JSON.stringify({
					data: {
						type: 'Participants',
						attributes: {
							participants: qualifiedEntrants.map((entrant) => ({
								name: entrant.name
							}))
						}
					}
				})
			}
		);

		const challongeData = await challongeResponse.json();

		if (!challongeResponse.ok) {
			console.error(
				'Failed to add participants to Challonge:',
				JSON.stringify(challongeData, null, 2)
			);

			return json(
				{
					message: 'Failed to create Challonge participants.',
					error: challongeData
				},
				{ status: challongeResponse.status }
			);
		}

		const { data } = challongeData;

		const challongeIdsByName: Record<string, string> = {};
		for (const participant of data) {
			challongeIdsByName[participant.attributes.name] = participant.id;
		}

		let updatedTournament;

		db.transaction((tx) => {
			tx.update(entrants)
				.set({
					status: 'not_qualified',
					updatedAt: new Date()
				})
				.where(eq(entrants.tournamentId, tournamentId))
				.run();

			if (qualifiedIds.length > 0) {
				tx.update(entrants)
					.set({
						status: 'qualified',
						updatedAt: new Date()
					})
					.where(and(eq(entrants.tournamentId, tournamentId), inArray(entrants.id, qualifiedIds)))
					.run();
			}

			for (const entrant of qualifiedEntrants) {
				const challongeParticipantId = challongeIdsByName[entrant.name!];

				if (!challongeParticipantId) {
					throw new Error(`Challonge participant not found: ${entrant.name}`);
				}

				tx.update(participants)
					.set({
						challongeParticipantId,
						updatedAt: new Date()
					})
					.where(
						and(
							eq(participants.id, entrant.participantId),
							eq(participants.tournamentId, tournamentId)
						)
					)
					.run();
			}

			[updatedTournament] = tx
				.update(tournaments)
				.set({
					status: 'elimination',
					updatedAt: new Date()
				})
				.where(eq(tournaments.id, tournamentId))
				.returning()
				.all();
		});

		return json({
			message: 'Qualifying completed.',
			tournament: updatedTournament
		});
	} catch (error) {
		console.error('Failed to complete qualifying:', error);

		return json({ message: 'Failed to complete qualifying.' }, { status: 500 });
	}
}
