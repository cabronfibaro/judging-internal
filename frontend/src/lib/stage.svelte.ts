type Driver = {
  participantId: string;
  state: "chase" | "lead";
  lead_score: number;
  chase_score: number;
};

type State = {
  p1: Driver;
  p2: Driver;
};

const defaultValue: State = {
  p1: {
    participantId: "",
    state: "lead",
    lead_score: -1,
    chase_score: -1,
  },
  p2: {
    participantId: "",
    state: "chase",
    lead_score: -1,
    chase_score: -1,
  },
};

class Stage {
  state = $state(defaultValue);
  matchId = $state<string | undefined>(undefined);

  clear() {
    this.matchId = undefined;
    this.state = {
      p1: { ...defaultValue.p1 },
      p2: { ...defaultValue.p2 },
    };
  }

  setDrivers(p1: string, p2: string, lead: "p1" | "p2") {
    this.state.p1.participantId = p1;
    this.state.p2.participantId = p2;

    this.state.p1.state = lead === "p1" ? "lead" : "chase";
    this.state.p2.state = lead === "p2" ? "lead" : "chase";
  }

  updateDriver(
    driver: keyof State,
    field: keyof Driver,
    value: string | number,
  ) {
    this.state[driver][field] = value as never;
  }

  setScore({ lead, chase }: { lead: number; chase: number }) {
    if (this.state.p1.state === "lead") {
      this.state.p1.lead_score = lead;
      this.state.p2.chase_score = chase;
    } else {
      this.state.p1.chase_score = chase;
      this.state.p2.lead_score = lead;
    }
  }

  getScores() {
    return [
      {
        participantId: this.state.p1.participantId,
        score:
          this.state.p1.state === "lead"
            ? this.state.p1.lead_score
            : this.state.p1.chase_score,
      },
      {
        participantId: this.state.p2.participantId,
        score:
          this.state.p2.state === "lead"
            ? this.state.p2.lead_score
            : this.state.p2.chase_score,
      },
    ];
  }

  resetScores() {
    this.state.p1.lead_score = -1;
    this.state.p1.chase_score = -1;
    this.state.p2.lead_score = -1;
    this.state.p2.chase_score = -1;
  }

  switchLeader() {
    this.state.p1.state = this.state.p1.state === "lead" ? "chase" : "lead";
    this.state.p2.state = this.state.p2.state === "lead" ? "chase" : "lead";
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
      return "";
    }

    const p1 = this.state.p1.lead_score + this.state.p1.chase_score;
    const p2 = this.state.p2.lead_score + this.state.p2.chase_score;

    return p1 > p2 ? this.state.p1.participantId : this.state.p2.participantId;
  }
}

export { Stage };
