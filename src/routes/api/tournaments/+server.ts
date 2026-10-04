import { json } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';
import { CHALLONGE_API_KEY } from '$env/static/private';
import { participants, entrants, tournaments, tracks, brackets } from '$lib/server/db/schema';
import { db } from '$lib/server/db/client';
import { snakeCase } from 'es-toolkit';
import { getJson } from '$lib/utils.js';
import type { GroupedMatches } from '../../../types/tournament.js';

export async function GET({ fetch }) {
	try {
		const result = await db
			.select({
				id: tournaments.id,
				name: tournaments.name,
				startingDate: tournaments.startingDate,
				track: tracks.name,
				qualifyingOrder: tournaments.qualifyingOrder,
				size: tournaments.size,
				status: tournaments.status,
				url: brackets.url
			})
			.from(tournaments)
			.leftJoin(tracks, eq(tournaments.trackId, tracks.id))
			.leftJoin(brackets, eq(tournaments.id, brackets.tournamentId));

		const tournamentsWithProgress = await Promise.all(
			result.map(async (tournament) => {
				const { matches } = await getJson<{ matches: GroupedMatches }>(
					fetch,
					`/api/tournaments/${tournament.id}/bracket/matches`
				);

				const allMatches = matches.flat();

				const completedMatches = allMatches.filter((match) =>
					match.every((node) => node.status === 'complete')
				).length;

				const progress = allMatches.length === 0 ? 0 : completedMatches / allMatches.length;

				return {
					...tournament,
					progress
				};
			})
		);

		return json({
			tournaments: tournamentsWithProgress
		});
	} catch (error) {
		console.error('Failed to get tournaments:', error);
		return json({ message: 'Failed to get tournaments.' }, { status: 500 });
	}
}

export async function POST({ request }) {
	const body = await request.json();

	const {
		name,
		size,
		startingDate,
		track,
		entrants: entrantsNames,
		randomizeQualifyingOrder
	} = body;

	if (
		typeof name !== 'string' ||
		typeof startingDate !== 'string' ||
		typeof track !== 'string' ||
		!Array.isArray(entrantsNames)
	) {
		return json({ message: 'Invalid request body.' }, { status: 400 });
	}

	const names = entrantsNames
		.filter((name): name is string => typeof name === 'string')
		.map((name) => name.trim())
		.filter(Boolean);

	const tournamentName = name.trim();
	const trackName = track.trim();

	if (!tournamentName || !startingDate || !trackName || names.length === 0) {
		return json({ message: 'Missing required fields.' }, { status: 400 });
	}

	const date = new Date(startingDate);

	if (Number.isNaN(date.getTime())) {
		return json({ message: 'Invalid starting date.' }, { status: 400 });
	}

	try {
		const challongeResponse = await fetch('https://api.challonge.com/v2.1/tournaments.json', {
			method: 'POST',
			headers: {
				Accept: 'application/json',
				'Content-Type': 'application/vnd.api+json',
				'Authorization-Type': 'v1',
				Authorization: CHALLONGE_API_KEY
			},
			body: JSON.stringify({
				data: {
					type: 'tournaments',
					attributes: {
						name: tournamentName,
						url: `${snakeCase(name)}_${Date.now()}`,
						tournament_type: 'single elimination',
						game_name: 'CarX Drift Racing', // TODO: make this configurable
						private: false,
						starts_at: date.toISOString(),
						description: `Top ${size} head-to-head tandem battle knockout bracket.`,
						match_options: {
							consolation_matches_target_rank: 3,
							accept_attachments: false
						},
						registration_options: {
							open_signup: false
						},
						seeding_options: {
							hide_seeds: false,
							sequential_pairings: false // TODO: make this configurable also
						},
						station_options: {
							auto_assign: false,
							only_start_matches_with_assigned_stations: false
						},
						group_stage_enabled: false
					}
				}
			})
		});

		const challongeData = await challongeResponse.json();

		if (!challongeResponse.ok) {
			console.error(
				'Failed to create Challonge tournament:',
				JSON.stringify(challongeData, null, 2)
			);

			return json(
				{
					message: 'Failed to create Challonge tournament.',
					error: challongeData
				},
				{ status: challongeResponse.status }
			);
		}

		const challongeUrl = challongeData.data.attributes.url;
		const challongeId = challongeData.data.id;

		const tournament = db.transaction((tx) => {
			const existingTrack = tx
				.select()
				.from(tracks)
				.where(eq(tracks.name, trackName))
				.limit(1)
				.all();

			let trackRow = existingTrack[0];

			if (!trackRow) {
				[trackRow] = tx
					.insert(tracks)
					.values({
						name: trackName
					})
					.returning()
					.all();
			}

			const qualifierRows = [];
			const participantIds = [];

			const [createdTournament] = tx
				.insert(tournaments)
				.values({
					name: tournamentName,
					startingDate: date,
					trackId: trackRow.id,
					qualifyingOrder: '',
					size,
					randomizeQualifyingOrder
				})
				.returning()
				.all();

			for (const participantName of names) {
				const existingParticipant = tx
					.select()
					.from(participants)
					.where(eq(participants.name, participantName))
					.limit(1)
					.all();

				let participant = existingParticipant[0];

				if (!participant) {
					[participant] = tx
						.insert(participants)
						.values({
							name: participantName,
							tournamentId: createdTournament.id
						})
						.returning()
						.all();
				}

				participantIds.push(participant.id);

				qualifierRows.push({
					tournamentId: createdTournament.id,
					participantId: participant.id
				});
			}

			if (randomizeQualifyingOrder) {
				for (let i = participantIds.length - 1; i > 0; i--) {
					const j = Math.floor(Math.random() * (i + 1));
					[participantIds[i], participantIds[j]] = [participantIds[j], participantIds[i]];
				}
			}

			const qualifyingOrder = participantIds.join(', ');

			tx.update(tournaments)
				.set({ qualifyingOrder })
				.where(eq(tournaments.id, createdTournament.id))
				.run();

			tx.insert(entrants).values(qualifierRows).run();
			tx.insert(brackets)
				.values({
					tournamentId: createdTournament.id,
					url: challongeUrl,
					challongeId: challongeId
				})
				.run();

			return createdTournament;
		});

		return json({ tournament }, { status: 201 });
	} catch (error) {
		console.error('Failed to create tournament:', error);

		return json({ message: 'Failed to create tournament.' }, { status: 500 });
	}
}
