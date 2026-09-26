type Result = {
	lead?: number;
	chase?: number;
};

class Judgement {
	results: Record<string, Result> = $state({});

	set(id: string, state: keyof Result, value: number) {
		this.results = {
			...this.results,
			[id]: {
				...this.results[id],
				[state]: value
			}
		};
	}

	get(id: string) {
		return this.results[id];
	}

	finalizedScore() {
		let lead = 0;
		let chase = 0;
		for (const judgement of Object.values(this.results)) {
			lead += judgement.lead ?? 0;
			chase += judgement.chase ?? 0;
		}

		return { lead, chase };
	}

	clear() {
		this.results = {};
	}
}

export { Judgement };
