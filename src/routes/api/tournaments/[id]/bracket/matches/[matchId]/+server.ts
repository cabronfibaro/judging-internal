import { db } from '#lib/server/db/client';
import { brackets, nodes } from '#lib/server/db/schema';
import { json } from '@sveltejs/kit';
import { and, eq } from 'drizzle-orm';

export async function PATCH({ params, request }) {
	const tournamentId = params.id;
	const matchId = params.matchId;

	if (!tournamentId || !matchId) {
		return json({ message: 'Tournament ID and match ID are required.' }, { status: 400 });
	}

	const body = await request.json();

	const { participantId, score, winnerId, status, appendScore = false } = body;

	if (!participantId) {
		return json({ message: 'Participant ID is required.' }, { status: 400 });
	}

	if (score !== undefined && typeof score !== 'number') {
		return json({ message: 'Score must be a number.' }, { status: 400 });
	}

	const bracket = db
		.select({ id: brackets.id })
		.from(brackets)
		.where(eq(brackets.tournamentId, tournamentId))
		.limit(1)
		.all()[0];

	if (!bracket) {
		return json({ message: 'Bracket not found.' }, { status: 404 });
	}

	const node = db
		.select({
			id: nodes.id,
			score: nodes.score
		})
		.from(nodes)
		.where(
			and(
				eq(nodes.bracketId, bracket.id),
				eq(nodes.matchId, matchId),
				eq(nodes.participantId, participantId)
			)
		)
		.limit(1)
		.all()[0];

	if (!node) {
		return json({ message: 'Participant node not found.' }, { status: 404 });
	}

	const updatedScore =
		score !== undefined ? (appendScore ? (node.score ?? 0) + score : score) : undefined;

	const updated = db
		.update(nodes)
		.set({
			...(updatedScore !== undefined && { score: updatedScore }),
			...(winnerId !== undefined && { winnerId }),
			...(status !== undefined && { status })
		})
		.where(eq(nodes.id, node.id))
		.returning()
		.all()[0];

	return json(updated);
}
