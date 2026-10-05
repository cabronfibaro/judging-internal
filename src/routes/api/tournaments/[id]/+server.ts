import { CHALLONGE_API_KEY } from '$app/env/private';
import { db } from '#lib/server/db/client';
import { brackets, tournaments } from '#lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';

export async function GET({ params }) {
	const id = params.id;
	if (!id) {
		return json({ message: 'Tournament ID is required.' }, { status: 400 });
	}

	try {
		const tournament = db
			.select()
			.from(tournaments)
			.where(eq(tournaments.id, id))
			.limit(1)
			.all()[0];

		if (!tournament) {
			return json({ message: 'Tournament not found.' }, { status: 404 });
		}

		return json(tournament);
	} catch (error) {
		console.error('Failed to get tournament:', error);
		return json({ message: 'Failed to get tournament.' }, { status: 500 });
	}
}

export async function DELETE({ params }) {
	const id = params.id;

	if (!id) {
		return json({ message: 'Tournament ID is required.' }, { status: 400 });
	}

	const bracket = db.select().from(brackets).where(eq(brackets.tournamentId, id)).limit(1).all()[0];
	if (!bracket) {
		return json({ message: 'Bracket not found fhor this tournament.' }, { status: 404 });
	}

	try {
		const deleted = await db.delete(tournaments).where(eq(tournaments.id, id)).returning({
			id: tournaments.id
		});

		if (deleted.length === 0) {
			return json({ message: 'Tournament not found.' }, { status: 404 });
		}

		await fetch(`https://api.challonge.com/v2.1/tournaments/${bracket.challongeId}.json`, {
			method: 'DELETE',
			headers: {
				'Authorization-Type': 'v1',
				Accept: 'application/json',
				'Content-Type': 'application/vnd.api+json',
				Authorization: CHALLONGE_API_KEY
			}
		});

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
	const { qualifyingOrder, status } = body;

	const updates = {
		...(qualifyingOrder !== undefined && { qualifyingOrder }),
		...(status !== undefined && { status })
	};

	try {
		const [updated] = db
			.update(tournaments)
			.set(updates)
			.where(eq(tournaments.id, id))
			.returning({
				id: tournaments.id,
				qualifyingOrder: tournaments.qualifyingOrder,
				status: tournaments.status
			})
			.all();

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
