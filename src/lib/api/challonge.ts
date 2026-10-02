import { getJson } from '$lib/utils';

export interface ParticipantReference {
	data: {
		id: string;
		type: string;
	};
}

export interface ParticipantScore {
	participant_id: number;
	scores: number[];
}

export interface MatchTimestamps {
	created_at: string;
	updated_at: string;
}

export interface MatchRelationships {
	player1: ParticipantReference;
	player2: ParticipantReference;
}

export interface MatchAttributes {
	state: string; // or 'complete' | 'pending' if strictly typed
	round: number;
	identifier: string;
	suggested_play_order: number;
	scores: string;
	score_in_sets: [number, number][];
	points_by_participant: ParticipantScore[];
	timestamps: MatchTimestamps;
	winner_id: number;
	relationships: MatchRelationships;
}

export interface MatchData {
	id: string;
	type: string;
	attributes: MatchAttributes;
}

export interface MatchResponse {
	data: MatchData[];
}

export async function getChallongeMatches(challongeTournamentId: string, apiKey: string) {
	const { data } = await getJson<MatchResponse>(
		fetch,
		`https://api.challonge.com/v2.1/tournaments/${challongeTournamentId}/matches.json`,
		{
			method: 'GET',
			headers: {
				'Content-Type': 'application/vnd.api+json',
				Accept: 'application/json',
				'Authorization-Type': 'v1',
				Authorization: apiKey
			}
		}
	);
	return data;
}
