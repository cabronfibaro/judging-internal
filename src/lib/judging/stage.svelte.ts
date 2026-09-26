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
	state = $state(defaultValue);

	updateDriver(driver: keyof State, field: keyof Driver, value: string | number) {
		this.state[driver][field] = value as never;
	}

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
