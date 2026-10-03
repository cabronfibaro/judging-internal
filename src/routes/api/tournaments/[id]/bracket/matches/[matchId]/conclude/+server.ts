import { db } from '$lib/server/db/client';
import { activeBattles, brackets, nodes, participants } from '$lib/server/db/schema';
import { CHALLONGE_API_KEY } from '$env/static/private';
import { json } from '@sveltejs/kit';
import { and, eq } from 'drizzle-orm';

export async function GET({ params }) {
	const tournamentId = params.id;
	const matchId = params.matchId;

	if (!tournamentId || !matchId) {
		return json({ message: 'Tournament ID and match ID are required.' }, { status: 400 });
	}

	const matchNodes = db
		.select({
			nodeId: nodes.id,
			participantId: nodes.participantId,
			score: nodes.score,
			challongeParticipantId: participants.challongeParticipantId
		})
		.from(nodes)
		.innerJoin(participants, eq(nodes.participantId, participants.id))
		.where(and(eq(nodes.matchId, matchId), eq(participants.tournamentId, tournamentId)))
		.all();

	if (matchNodes.length !== 2) {
		return json(
			{
				message: 'Match must contain exactly two participants.',
				nodes: matchNodes
			},
			{ status: 400 }
		);
	}

	const [p1, p2] = matchNodes;

	if (!p1.challongeParticipantId || !p2.challongeParticipantId) {
		return json(
			{
				message: 'One or more participants do not have a Challonge participant ID.'
			},
			{ status: 400 }
		);
	}

	const p1Score = p1.score ?? 0;
	const p2Score = p2.score ?? 0;

	// TODO: make conclude force later

	if (p1Score === p2Score) {
		return json(
			{
				message: 'Match cannot be concluded because the scores are tied.',
				p1Score,
				p2Score
			},
			{ status: 400 }
		);
	}

	const winner = p1Score > p2Score ? p1 : p2;
	const loser = p1Score > p2Score ? p2 : p1;

	const bracket = db
		.select({
			url: brackets.url,
			challongeId: brackets.challongeId
		})
		.from(brackets)
		.where(eq(brackets.tournamentId, tournamentId))
		.limit(1)
		.all()[0];

	if (!bracket) {
		return json({ message: 'Bracket not found.' }, { status: 404 });
	}

	// update local winner.
	db.update(nodes)
		.set({
			winnerId: winner.challongeParticipantId
		})
		.where(eq(nodes.matchId, matchId))
		.run();

	// update Challonge.
	const response = await fetch(
		`https://api.challonge.com/v2.1/tournaments/${bracket.url}/matches/${matchId}.json`,
		{
			method: 'PUT',
			headers: {
				'Content-Type': 'application/vnd.api+json',
				Accept: 'application/json',
				'Authorization-Type': 'v1',
				Authorization: CHALLONGE_API_KEY
			},
			body: JSON.stringify({
				data: {
					type: 'match',
					attributes: {
						match: [
							{
								participant_id: winner.challongeParticipantId,
								score_set: String(winner.score ?? 0),
								rank: 1,
								advancing: true
							},
							{
								participant_id: loser.challongeParticipantId,
								score_set: String(loser.score ?? 0),
								rank: 2
							}
						],
						tie: false
					}
				}
			})
		}
	);

	const data = await response.json().catch(() => null);

	if (!response.ok) {
		return json(
			{
				message: 'Failed to update Challonge match.',
				challonge: data
			},
			{ status: response.status }
		);
	}

	// match is concluded, so remove its active battle state
	db.delete(activeBattles).where(eq(activeBattles.matchId, matchId)).run();

	return json({
		winner: {
			participantId: winner.participantId,
			challongeParticipantId: winner.challongeParticipantId,
			score: winner.score
		},
		loser: {
			participantId: loser.participantId,
			challongeParticipantId: loser.challongeParticipantId,
			score: loser.score
		},
		challonge: data
	});
}
