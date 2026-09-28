import { db } from '$lib/server/db/client';
import { entrants, tournaments } from '$lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { and, eq } from 'drizzle-orm';

export async function GET({ params }) {
	const tournamentId = params.id;

	if (!tournamentId) {
		return json({ message: 'Tournament ID is required.' }, { status: 400 });
	}

	try {
		let updatedTournament;

		db.transaction((tx) => {
			const [tournament] = tx
				.select({
					size: tournaments.size
				})
				.from(tournaments)
				.where(eq(tournaments.id, tournamentId))
				.all();

			if (!tournament) {
				throw new Error('Tournament not found.');
			}

			const tournamentEntrants = tx
				.select()
				.from(entrants)
				.where(eq(entrants.tournamentId, tournamentId))
				.all();

			const sorted = tournamentEntrants.sort((a, b) => b.score - a.score);

			for (const [index, entrant] of sorted.entries()) {
				tx.update(entrants)
					.set({
						status: index < tournament.size ? 'qualified' : 'not_qualified',
						updatedAt: new Date()
					})
					.where(and(eq(entrants.id, entrant.id), eq(entrants.tournamentId, tournamentId)))
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
