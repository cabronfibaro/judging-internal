import { json } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';
import { participants, qualifiers, tournaments, tracks } from '$lib/server/db/schema';
import { db } from '$lib/server/db/client';

export async function GET() {
	try {
		const result = await db
			.select({
				id: tournaments.id,
				name: tournaments.name,
				startingDate: tournaments.startingDate,
				track: tracks.name
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

	const { name, startingDate, track, qualifiers: qualifierNames } = body;

	if (
		typeof name !== 'string' ||
		typeof startingDate !== 'string' ||
		typeof track !== 'string' ||
		!Array.isArray(qualifierNames)
	) {
		return json({ message: 'Invalid request body.' }, { status: 400 });
	}

	const names = qualifierNames
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
					trackId: trackRow.id
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

			tx.insert(qualifiers).values(qualifierRows).run();

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
