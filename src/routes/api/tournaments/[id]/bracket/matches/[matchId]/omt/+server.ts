import { db } from '$lib/server/db/client';
import { activeBattles, brackets, entrants, nodes } from '$lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { and, eq } from 'drizzle-orm';

export async function GET({ params }) {
	const tournamentId = params.id;
	const matchId = params.matchId;

	if (!tournamentId || !matchId) {
		return json({ message: 'Tournament ID and match ID are required.' }, { status: 400 });
	}

	const result = db.transaction((tx) => {
		const bracket = tx
			.select({ id: brackets.id })
			.from(brackets)
			.where(eq(brackets.tournamentId, tournamentId))
			.limit(1)
			.all()[0];

		if (!bracket) {
			return {
				error: json({ message: 'Bracket not found.' }, { status: 404 })
			};
		}

		const battles = tx
			.select({
				id: activeBattles.id,
				participantId: activeBattles.participantId,
				rerun: activeBattles.rerun
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
			return {
				error: json(
					{
						message: 'Active battle must contain exactly two participants.'
					},
					{ status: 400 }
				)
			};
		}

		const [p1, p2] = battles;

		const p1Entrant = tx
			.select({
				score: entrants.score
			})
			.from(entrants)
			.where(
				and(eq(entrants.tournamentId, tournamentId), eq(entrants.participantId, p1.participantId))
			)
			.limit(1)
			.all()[0];

		const p2Entrant = tx
			.select({
				score: entrants.score
			})
			.from(entrants)
			.where(
				and(eq(entrants.tournamentId, tournamentId), eq(entrants.participantId, p2.participantId))
			)
			.limit(1)
			.all()[0];

		if (!p1Entrant || !p2Entrant) {
			return {
				error: json({ message: 'Entrant data not found.' }, { status: 404 })
			};
		}

		const rerun = (p1.rerun ?? 0) + 1;

		const p1Score = p1Entrant.score ?? 0;
		const p2Score = p2Entrant.score ?? 0;

		const higherScoreLeads = rerun % 2 === 0;

		let p1Leads: boolean;

		if (p1Score === p2Score) {
			p1Leads = rerun % 2 === 0;
		} else if (higherScoreLeads) {
			p1Leads = p1Score > p2Score;
		} else {
			p1Leads = p1Score < p2Score;
		}

		tx.update(nodes)
			.set({ score: 0 })
			.where(
				and(
					eq(nodes.bracketId, bracket.id),
					eq(nodes.matchId, matchId),
					eq(nodes.participantId, p1.participantId)
				)
			)
			.run();

		tx.update(nodes)
			.set({ score: 0 })
			.where(
				and(
					eq(nodes.bracketId, bracket.id),
					eq(nodes.matchId, matchId),
					eq(nodes.participantId, p2.participantId)
				)
			)
			.run();

		tx.update(activeBattles)
			.set({
				rerun,
				leadStatus: p1Leads ? 1 : 0
			})
			.where(eq(activeBattles.id, p1.id))
			.run();

		tx.update(activeBattles)
			.set({
				rerun,
				leadStatus: p1Leads ? 0 : 1
			})
			.where(eq(activeBattles.id, p2.id))
			.run();

		return {
			rerun,
			p1: {
				participantId: p1.participantId,
				score: p1Score,
				leadStatus: p1Leads ? 1 : 0
			},
			p2: {
				participantId: p2.participantId,
				score: p2Score,
				leadStatus: p1Leads ? 0 : 1
			}
		};
	});

	if ('error' in result) {
		return json({ error: result.error });
	}

	return json(result);
}
