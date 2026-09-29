import { db } from '$lib/server/db/client';
import { brackets } from '$lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';

export async function GET({ params }) {
	const tournamentId = params.id;

	if (!tournamentId) {
		return json({ message: 'Tournament ID is required.' }, { status: 400 });
	}

	try {
		const bracket = db
			.select()
			.from(brackets)
			.where(eq(brackets.tournamentId, tournamentId))
			.limit(1)
			.all()[0];

		if (!bracket) {
			return json({ message: 'Bracket not found.' }, { status: 404 });
		}

		return json({
			bracket
		});
	} catch (error) {
		console.error('Failed to get bracket:', error);

		return json({ message: 'Failed to get bracket.' }, { status: 500 });
	}
}
