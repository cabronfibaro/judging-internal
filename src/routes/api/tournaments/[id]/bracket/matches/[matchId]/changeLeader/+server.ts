import { db } from '#lib/server/db/client';
import { activeBattles, brackets, nodes } from '#lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { and, eq } from 'drizzle-orm';

export async function GET({ params }) {
	const tournamentId = params.id;
	const matchId = params.matchId;

	if (!tournamentId || !matchId) {
		return json({ message: 'Tournament ID and match ID are required.' }, { status: 400 });
	}

	const bracket = db
		.select({
			id: brackets.id
		})
		.from(brackets)
		.where(eq(brackets.tournamentId, tournamentId))
		.limit(1)
		.all()[0];

	if (!bracket) {
		return json({ message: 'Bracket not found.' }, { status: 404 });
	}

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
				eq(nodes.bracketId, bracket.id),
				eq(nodes.participantId, activeBattles.participantId)
			)
		)
		.where(eq(activeBattles.matchId, matchId))
		.all();

	if (battles.length !== 2) {
		return json(
			{
				message: 'Active battle must contain exactly two participants.',
				battles
			},
			{ status: 400 }
		);
	}

	const p1 = battles[0];
	const p2 = battles[1];

	db.transaction((tx) => {
		// 1, 0 -> -1, 1
		if (p1.leadStatus === 1 && p2.leadStatus === 0) {
			tx.update(activeBattles).set({ leadStatus: -1 }).where(eq(activeBattles.id, p1.id)).run();

			tx.update(activeBattles).set({ leadStatus: 1 }).where(eq(activeBattles.id, p2.id)).run();
		}

		// 0, 1 -> 1, -1
		else if (p1.leadStatus === 0 && p2.leadStatus === 1) {
			tx.update(activeBattles).set({ leadStatus: 1 }).where(eq(activeBattles.id, p1.id)).run();

			tx.update(activeBattles).set({ leadStatus: -1 }).where(eq(activeBattles.id, p2.id)).run();
		}

		// -1, 1 -> -1, -1
		else if (p1.leadStatus === -1 && p2.leadStatus === 1) {
			tx.update(activeBattles).set({ leadStatus: -1 }).where(eq(activeBattles.id, p2.id)).run();
		}

		// 1, -1 -> -1, -1
		else if (p1.leadStatus === 1 && p2.leadStatus === -1) {
			tx.update(activeBattles).set({ leadStatus: -1 }).where(eq(activeBattles.id, p1.id)).run();
		}

		// -1, -1 -> do nothing
	});

	const updatedBattles = db
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
				eq(nodes.bracketId, bracket.id),
				eq(nodes.participantId, activeBattles.participantId)
			)
		)
		.where(eq(activeBattles.matchId, matchId))
		.all();

	return json(updatedBattles);
}
