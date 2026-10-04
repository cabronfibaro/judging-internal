export type Tournament = {
	id: string;
	name: string;
	startingDate: Date | string;
	track: string | null;
	qualifyingOrder: string;
	status: 'qualifying' | 'elimination';
	url: string;
	progress: number;
};

export type Entrant = {
	id: string;
	challongeParticipantId: string;
	participantId: string;
	name: string;
	status: 'qualified' | 'not_qualified' | 'pending' | 'scored';
	score: number;
};

export type Match = {
	matchId: string;
	participant_name: string;
	status: 'pending' | 'open' | 'complete';
	round: number;
	winnerId: string | null;
	participantId: string | null;
	challongeParticipantId: string | null;
	score: number;
	bothRunFinished: number;
};

export type GroupedMatches = Match[][][];
