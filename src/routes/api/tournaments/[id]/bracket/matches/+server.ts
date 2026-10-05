import { db } from '#lib/server/db/client';
import { activeBattles, brackets, nodes, participants } from '#lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { eq, sql } from 'drizzle-orm';

export async function GET({ params }) {
	const tournamentId = params.id;

	if (!tournamentId) {
		return json({ message: 'Tournament ID is required.' }, { status: 400 });
	}

	try {
		const battleStatus = db
			.select({
				matchId: activeBattles.matchId,
				complete: sql<boolean>`
					COUNT(*) = 2
					AND SUM(
						CASE
							WHEN ${activeBattles.leadStatus} = -1 THEN 1
							ELSE 0
						END
					) = 2`.as('complete')
			})
			.from(activeBattles)
			.groupBy(activeBattles.matchId)
			.as('battle_status');

		const rows = db
			.select({
				matchId: nodes.matchId,
				participant_name: participants.name,
				status: nodes.status,
				round: nodes.round,
				winnerId: nodes.winnerId,
				challongeParticipantId: participants.challongeParticipantId,
				participantId: participants.id,
				score: nodes.score,
				bothRunFinished: sql<boolean>`COALESCE(${battleStatus.complete}, 1)`
			})
			.from(nodes)
			.innerJoin(brackets, eq(nodes.bracketId, brackets.id))
			.leftJoin(battleStatus, eq(battleStatus.matchId, nodes.matchId))
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
