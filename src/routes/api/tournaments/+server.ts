import { json } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';
import { participants, entrants, tournaments, tracks } from '$lib/server/db/schema';
import { db } from '$lib/server/db/client';

export async function GET() {
	try {
		const result = await db
			.select({
				id: tournaments.id,
				name: tournaments.name,
				startingDate: tournaments.startingDate,
				track: tracks.name,
				qualifyingOrder: tournaments.qualifyingOrder
			})
			.from(tournaments)
			.leftJoin(tracks, eq(tournaments.trackId, tracks.id));

		return json({
			tournaments: result
		});
	} catch (error) {
		console.error('Failed to get tournaments:', error);

		return json({ message: 'Failed to get tournaments.' }, { status: 500 });
	}
}

export async function POST({ request }) {
	const body = await request.json();

	const { name, startingDate, track, entrants: entrantsNames } = body;

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

	const qualifyingOrder = Array.from({ length: names.length }, (_, i) => i);
	for (let i = qualifyingOrder.length - 1; i > 0; i--) {
		const j = Math.floor(Math.random() * (i + 1));
		[qualifyingOrder[i], qualifyingOrder[j]] = [qualifyingOrder[j], qualifyingOrder[i]];
	}

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

			const [createdTournament] = tx
				.insert(tournaments)
				.values({
					name: tournamentName,
					startingDate: date,
					trackId: trackRow.id,
					qualifyingOrder: qualifyingOrder.join(', ')
				})
				.returning()
				.all();

			const qualifierRows = [];

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

				qualifierRows.push({
					tournamentId: createdTournament.id,
					participantId: participant.id
				});
			}

			tx.insert(entrants).values(qualifierRows).run();

			return createdTournament;
		});

		return json(
			{
				tournament
			},
			{ status: 201 }
		);
	} catch (error) {
		console.error('Failed to create tournament:', error);

		return json({ message: 'Failed to create tournament.' }, { status: 500 });
	}
}
