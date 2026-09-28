export type Tournament = {
	id: string;
	name: string;
	startingDate: Date | string;
	track: string | null;
	qualifyingOrder: string;
	status: 'qualifying' | 'elimination';
};

export type Entrant = {
	id: string;
	participantId: string;
	name: string;
	status: 'qualified' | 'not_qualified' | 'pending' | 'scored';
	score: number;
};
