import { db } from '$lib/server/db/client';
import { entrants } from '$lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { and, eq } from 'drizzle-orm';

export async function PATCH({ params, request }) {
	const tournamentId = params.id;
	const entrantId = params.entrantId;

	if (!tournamentId || !entrantId) {
		return json({ message: 'Tournament ID and entrant ID are required.' }, { status: 400 });
	}

	const body = await request.json();

	const { status, score } = body;

	if (
		status !== undefined &&
		status !== 'qualified' &&
		status !== 'not_qualified' &&
		status !== 'pending' &&
		status !== 'scored'
	) {
		return json({ message: 'Invalid entrant status.' }, { status: 400 });
	}

	if (score !== undefined && typeof score !== 'number') {
		return json({ message: 'Score must be a number.' }, { status: 400 });
	}

	if (status === undefined && score === undefined) {
		return json({ message: 'Nothing to update.' }, { status: 400 });
	}

	try {
		const [updatedEntrant] = await db
			.update(entrants)
			.set({
				...(status !== undefined && { status }),
				...(score !== undefined && { score })
			})
			.where(and(eq(entrants.id, entrantId), eq(entrants.tournamentId, tournamentId)))
			.returning();

		if (!updatedEntrant) {
			return json({ message: 'Entrant not found.' }, { status: 404 });
		}

		return json({
			entrant: updatedEntrant
		});
	} catch (error) {
		console.error('Failed to update qualifier:', error);

		return json({ message: 'Failed to update qualifier.' }, { status: 500 });
	}
}
