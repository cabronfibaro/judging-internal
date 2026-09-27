import { db } from '$lib/server/db/client';
import { tournaments } from '$lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';

export async function DELETE({ params }) {
	const id = params.id;

	if (!id) {
		return json({ message: 'Tournament ID is required.' }, { status: 400 });
	}

	try {
		const deleted = await db.delete(tournaments).where(eq(tournaments.id, id)).returning({
			id: tournaments.id
		});

		if (deleted.length === 0) {
			return json({ message: 'Tournament not found.' }, { status: 404 });
		}

		return json({
			id: deleted[0].id
		});
	} catch (error) {
		console.error('Failed to delete tournament:', error);

		return json({ message: 'Failed to delete tournament.' }, { status: 500 });
	}
}

export async function PATCH({ params, request }) {
	const id = params.id;

	if (!id) {
		return json({ message: 'Tournament ID is required.' }, { status: 400 });
	}

	const body = await request.json();
	const { qualifyingOrder } = body;

	try {
		const [updated] = await db
			.update(tournaments)
			.set({
				qualifyingOrder: qualifyingOrder
			})
			.where(eq(tournaments.id, id))
			.returning({
				id: tournaments.id,
				qualifyingOrder: tournaments.qualifyingOrder
			});

		if (!updated) {
			return json({ message: 'Tournament not found.' }, { status: 404 });
		}

		return json({
			tournament: updated
		});
	} catch (error) {
		console.error('Failed to update tournament:', error);

		return json({ message: 'Failed to update tournament.' }, { status: 500 });
	}
}
