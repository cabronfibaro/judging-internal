import { db } from '$lib/server/db/client';
import { entrants } from '$lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { and, eq } from 'drizzle-orm';

export async function PATCH({ params, request }) {
	const tournamentId = params.id;
	const entrantId = params.entrantId;

	if (!tournamentId || !entrantId) {
		return json({ message: 'Tournament ID and qualifier ID are required.' }, { status: 400 });
	}

	const body = await request.json();

	const { status, score } = body;

	if (
		status !== undefined &&
		status !== 'qualified' &&
		status !== 'not_qualified' &&
		status !== 'pending'
	) {
		return json({ message: 'Invalid qualifier status.' }, { status: 400 });
	}

	if (score !== undefined && typeof score !== 'number') {
		return json({ message: 'Score must be a number.' }, { status: 400 });
	}

	if (status === undefined && score === undefined) {
		return json({ message: 'Nothing to update.' }, { status: 400 });
	}

	try {
		const [updatedQualifier] = await db
			.update(entrants)
			.set({
				...(status !== undefined && { status }),
				...(score !== undefined && { score })
			})
			.where(and(eq(entrants.id, entrantId), eq(entrants.tournamentId, tournamentId)))
			.returning();

		if (!updatedQualifier) {
			return json({ message: 'Qualifier not found.' }, { status: 404 });
		}

		return json({
			qualifier: updatedQualifier
		});
	} catch (error) {
		console.error('Failed to update qualifier:', error);

		return json({ message: 'Failed to update qualifier.' }, { status: 500 });
	}
}
