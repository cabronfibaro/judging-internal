import { db } from '$lib/server/db/client';
import { brackets, nodes, participants } from '$lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';

export async function GET({ params }) {
	const tournamentId = params.id;

	if (!tournamentId) {
		return json({ message: 'Tournament ID is required.' }, { status: 400 });
	}

	try {
		const rows = db
			.select({
				matchId: nodes.matchId,
				participant_name: participants.name,
				status: nodes.status,
				round: nodes.round,
				winnerId: nodes.winnerId,
				participantId: participants.challongeParticipantId
			})
			.from(nodes)
			.innerJoin(brackets, eq(nodes.bracketId, brackets.id))
			.leftJoin(participants, eq(participants.id, nodes.participantId))
			.where(eq(brackets.tournamentId, tournamentId))
			.all();

		const rounds = Object.groupBy(rows, ({ round }) => round);
		const groupedMatches = Object.values(rounds).map((round) => {
			if (!round) return [];
			return Object.values(Object.groupBy(round, ({ matchId }) => matchId));
		});

		return json({
			matches: groupedMatches
		});
	} catch (error) {
		console.error('Failed to get bracket matches:', error);
		return json({ message: 'Failed to get bracket matches.' }, { status: 500 });
	}
}
