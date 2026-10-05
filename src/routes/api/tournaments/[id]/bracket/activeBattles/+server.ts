import { db } from '#lib/server/db/client';
import { activeBattles, brackets, nodes } from '#lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { and } from 'drizzle-orm';
import { eq } from 'drizzle-orm';

export async function GET({ params }) {
	const tournamentId = params.id;
	const battles = db
		.select({
			id: activeBattles.id,
			matchId: activeBattles.matchId,
			participantId: activeBattles.participantId,
			leadStatus: activeBattles.leadStatus
		})
		.from(activeBattles)
		.innerJoin(
			nodes,
			and(
				eq(nodes.matchId, activeBattles.matchId),
				eq(nodes.participantId, activeBattles.participantId)
			)
		)
		.innerJoin(brackets, eq(brackets.id, nodes.bracketId))
		.where(eq(brackets.tournamentId, tournamentId))
		.all();

	return json(battles);
}
