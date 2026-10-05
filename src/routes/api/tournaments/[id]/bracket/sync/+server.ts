/* eslint-disable @typescript-eslint/no-explicit-any */
import { db } from '#lib/server/db/client';
import {
	activeBattles,
	brackets,
	entrants,
	nodes,
	participants,
	type Node
} from '#lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { and, eq, inArray } from 'drizzle-orm';
import { getJson } from '#lib/utils.js';
import type { GroupedMatches } from '../../../../../../types/tournament.js';
import { getChallongeMatches } from '#lib/api/challonge.js';
import { CHALLONGE_API_KEY } from '$app/env/private';

export async function GET({ params, fetch }) {
	const tournamentId = params.id;

	if (!tournamentId) {
		return json({ message: 'Tournament ID is required.' }, { status: 400 });
	}

	try {
		const bracket = db
			.select({
				id: brackets.id,
				challongeTournamentId: brackets.url
			})
			.from(brackets)
			.where(eq(brackets.tournamentId, tournamentId))
			.limit(1)
			.all()[0];

		if (!bracket) {
			return json({ message: 'Bracket not found.' }, { status: 404 });
		}

		const matches = await getChallongeMatches(bracket.challongeTournamentId, CHALLONGE_API_KEY);
		const challongeParticipantIds = [
			...new Set(
				matches.flatMap((match) =>
					match.attributes.points_by_participant
						.map((entry) => String(entry.participant_id))
						.filter(Boolean)
				)
			)
		] as string[];

		const localParticipants = db
			.select({
				id: participants.id,
				challongeParticipantId: participants.challongeParticipantId
			})
			.from(participants)
			.where(
				and(
					eq(participants.tournamentId, tournamentId),
					inArray(participants.challongeParticipantId, challongeParticipantIds)
				)
			)
			.all();

		const participantMap = new Map(
			localParticipants.map((participant) => [participant.challongeParticipantId, participant.id])
		);

		const matchIds = matches.map((match: any) => match.id);

		const existingNodes = db
			.select()
			.from(nodes)
			.where(and(eq(nodes.bracketId, bracket.id), inArray(nodes.matchId, matchIds)))
			.all();

		db.transaction((tx) => {
			for (const match of matches) {
				const { state, round, points_by_participant, winner_id } = match.attributes;

				// Get the existing two nodes for this match.
				const matchNodes = existingNodes.filter((node) => node.matchId === match.id);

				for (const [index, entry] of points_by_participant.entries()) {
					const participantId = entry.participant_id
						? participantMap.get(String(entry.participant_id))
						: null;

					if (entry.participant_id && !participantId) {
						throw new Error(
							`Local participant not found for Challonge participant ${entry.participant_id}`
						);
					}

					const existingNode = matchNodes[index];

					const values = {
						bracketId: bracket.id,
						participantId,
						matchId: match.id,
						status: state,
						round,
						winnerId: winner_id != null ? String(winner_id) : null
					};

					if (existingNode) {
						tx.update(nodes)
							.set({
								participantId: values.participantId,
								status: values.status as Node['status'],
								round: values.round,
								winnerId: values.winnerId
							})
							.where(eq(nodes.id, existingNode.id))
							.run();
					} else {
						tx.insert(nodes).values(values).run();
					}
				}
			}
		});
		const { matches: groupedMatches } = await getJson<{ matches: GroupedMatches }>(
			fetch,
			`/api/tournaments/${tournamentId}/bracket/matches`
		);

		db.transaction((tx) => {
			for (const round of groupedMatches) {
				for (const match of round) {
					if (match.length !== 2) {
						console.log('match length is not 2');
						continue;
					}

					const p1 = match[0];
					const p2 = match[1];

					if (!p1.participantId || !p2.participantId) {
						continue;
					}

					if (!p1.challongeParticipantId || !p2.challongeParticipantId) {
						continue;
					}

					if (p1.status === 'complete' || p2.status === 'complete') {
						continue;
					}

					const p1Entrant = tx
						.select({
							score: entrants.score
						})
						.from(entrants)
						.where(
							and(
								eq(entrants.tournamentId, tournamentId),
								eq(entrants.participantId, p1.participantId)
							)
						)
						.limit(1)
						.all()[0];

					const p2Entrant = tx
						.select({
							score: entrants.score
						})
						.from(entrants)
						.where(
							and(
								eq(entrants.tournamentId, tournamentId),
								eq(entrants.participantId, p2.participantId)
							)
						)
						.limit(1)
						.all()[0];

					if (!p1Entrant || !p2Entrant) {
						continue;
					}

					const existingBattles = tx
						.select({
							id: activeBattles.id
						})
						.from(activeBattles)
						.where(eq(activeBattles.matchId, p1.matchId))
						.all();

					if (existingBattles.length >= 2) {
						// console.log(`active battle already exists for match ${p1.matchId}`);
						continue;
					}

					let p1Leads: boolean;
					let p2Leads: boolean;

					if (p1Entrant.score === p2Entrant.score) {
						p1Leads = Math.random() < 0.5;
						p2Leads = !p1Leads;
					} else {
						p1Leads = p1Entrant.score > p2Entrant.score;
						p2Leads = !p1Leads;
					}

					tx.insert(activeBattles)
						.values([
							{
								matchId: p1.matchId,
								participantId: p1.participantId,
								leadStatus: p1Leads ? 1 : 0
							},
							{
								matchId: p2.matchId,
								participantId: p2.participantId,
								leadStatus: p2Leads ? 1 : 0
							}
						])
						.run();
				}
			}
		});

		const battles = db
			.select({
				id: activeBattles.id,
				matchId: activeBattles.matchId,
				participantId: activeBattles.participantId,
				leadStatus: activeBattles.leadStatus
			})
			.from(activeBattles)
			.innerJoin(nodes, eq(nodes.matchId, activeBattles.matchId))
			.innerJoin(brackets, eq(brackets.id, nodes.bracketId))
			.where(eq(brackets.tournamentId, tournamentId))
			.all();

		return json({
			matches: groupedMatches,
			activeBattles: battles
		});
	} catch (error) {
		console.error('Failed to sync bracket:', error);

		return json({ message: 'Failed to sync bracket.' }, { status: 500 });
	}
}
