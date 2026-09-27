export type Entrant = {
	id: string;
	participantId: string;
	name: string;
	status: 'qualified' | 'not_qualified' | 'pending' | 'scored';
	score: number;
};
