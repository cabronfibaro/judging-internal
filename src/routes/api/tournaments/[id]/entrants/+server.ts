import { db } from '#lib/server/db/client';
import { participants, entrants } from '#lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';

export async function GET({ params }) {
	const tournamentId = params.id;

	if (!tournamentId) {
		return json({ message: 'Tournament ID is required.' }, { status: 400 });
	}

	try {
		const result = await db
			.select({
				id: entrants.id,
				participantId: entrants.participantId,
				name: participants.name,
				status: entrants.status,
				score: entrants.score
			})
			.from(entrants)
			.innerJoin(participants, eq(entrants.participantId, participants.id))
			.where(eq(entrants.tournamentId, tournamentId));

		return json({
			entrants: result
		});
	} catch (error) {
		console.error('Failed to get tournament entrants:', error);

		return json({ message: 'Failed to get tournament entrants.' }, { status: 500 });
	}
}
