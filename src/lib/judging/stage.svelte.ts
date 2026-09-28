import type { Entrant } from '../../types/tournament';

type Driver = {
	name: string;
	state: 'chase' | 'lead';
	lead_score: number;
	chase_score: number;
};

type State = {
	p1: Driver;
	p2: Driver;
};

const defaultValue: State = {
	p1: {
		name: '',
		state: 'lead',
		lead_score: -1,
		chase_score: -1
	},
	p2: {
		name: '',
		state: 'chase',
		lead_score: -1,
		chase_score: -1
	}
};

class Stage {
	tournamentId = $state('');

	onStage = $state(0);
	orders = $state<number[]>([]);

	// TODO: make queue handle Entrant and the type of whatever to handle the battles
	queue = $state<Entrant[]>([]);
	isQualifying = $state(false);

	state = $state(defaultValue);

	clear(tournamentId = '') {
		this.tournamentId = tournamentId;
		this.orders = [];
		this.queue = [];
		this.isQualifying = false;
		this.state = defaultValue;
	}

	updateDriver(driver: keyof State, field: keyof Driver, value: string | number) {
		this.state[driver][field] = value as never;
	}

	frontQueue() {
		return this.queue[0];
	}

	move(src: number, to: number) {
		const [queueItem] = this.queue.splice(src, 1);
		this.queue.splice(to, 0, queueItem);

		const [orderItem] = this.orders.splice(src, 1);
		this.orders.splice(to, 0, orderItem);
	}

	// <-- QUALIFYING BLOCK -->
	isQualifyingComplete() {
		return !this.queue.some((entrant) => entrant.status === 'pending');
	}

	enqueueEntrants(entrants: Entrant[], orders: string) {
		this.isQualifying = true;
		this.orders = orders.split(', ').map((num) => Number(num));
		this.queue = this.orders.map((index) => entrants[index]);
	}

	updateEntrants(entrant: Partial<Entrant> & { id: string }) {
		const existing = this.queue.find((q) => q.id === entrant.id);

		if (!existing) {
			return;
		}

		Object.assign(existing, entrant);
	}
	// <-- END OF QUALIFYING BLOCK -->

	setScore({ lead, chase }: { lead: number; chase: number }) {
		if (this.state.p1.state === 'lead') {
			this.state.p1.lead_score = lead;
			this.state.p2.chase_score = chase;
		} else {
			this.state.p1.chase_score = chase;
			this.state.p2.lead_score = lead;
		}
	}

	switchLeader() {
		this.state.p1.state = this.state.p1.state === 'lead' ? 'chase' : 'lead';
		this.state.p2.state = this.state.p2.state === 'lead' ? 'chase' : 'lead';
	}

	// TODO: think about this later
	isConcluded() {
		return (
			this.state.p1.lead_score !== -1 &&
			this.state.p1.chase_score !== -1 &&
			this.state.p2.lead_score !== -1 &&
			this.state.p2.chase_score !== -1
		);
	}

	winner() {
		if (!this.isConcluded()) {
			return '';
		}

		const p1 = this.state.p1.lead_score + this.state.p1.chase_score;
		const p2 = this.state.p1.lead_score + this.state.p1.chase_score;
		return p1 > p2 ? this.state.p1.name : this.state.p2.name;
	}
}

export { Stage };
