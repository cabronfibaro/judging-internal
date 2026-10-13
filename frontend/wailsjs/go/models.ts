export namespace entrant {
	
	export class Entrant {
	    id: string;
	    participantId: string;
	    tournamentId: string;
	    name: string;
	    status: string;
	    score?: number;
	
	    static createFrom(source: any = {}) {
	        return new Entrant(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.participantId = source["participantId"];
	        this.tournamentId = source["tournamentId"];
	        this.name = source["name"];
	        this.status = source["status"];
	        this.score = source["score"];
	    }
	}

}

export namespace repository {
	
	export class ActiveBattle {
	    id: string;
	    matchId: string;
	    participantId: string;
	    leadStatus: number;
	    rerun: number;
	
	    static createFrom(source: any = {}) {
	        return new ActiveBattle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.matchId = source["matchId"];
	        this.participantId = source["participantId"];
	        this.leadStatus = source["leadStatus"];
	        this.rerun = source["rerun"];
	    }
	}
	export class CreateTournamentRequest {
	    name: string;
	    size: number;
	    startingDate: string;
	    track: string;
	    entrants: string[];
	    randomizeQualifyingOrder: boolean;
	    qualifyingOrder: string[];
	
	    static createFrom(source: any = {}) {
	        return new CreateTournamentRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.size = source["size"];
	        this.startingDate = source["startingDate"];
	        this.track = source["track"];
	        this.entrants = source["entrants"];
	        this.randomizeQualifyingOrder = source["randomizeQualifyingOrder"];
	        this.qualifyingOrder = source["qualifyingOrder"];
	    }
	}
	export class Match {
	    matchId: string;
	    participant_name?: string;
	    status: string;
	    round: number;
	    winnerId?: string;
	    challongeParticipantId?: string;
	    participantId?: string;
	    score?: number;
	    bothRunFinished: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Match(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.matchId = source["matchId"];
	        this.participant_name = source["participant_name"];
	        this.status = source["status"];
	        this.round = source["round"];
	        this.winnerId = source["winnerId"];
	        this.challongeParticipantId = source["challongeParticipantId"];
	        this.participantId = source["participantId"];
	        this.score = source["score"];
	        this.bothRunFinished = source["bothRunFinished"];
	    }
	}
	export class UpdateEntrantParams {
	    TournamentID: string;
	    EntrantID: string;
	    Status?: string;
	    Score?: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateEntrantParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TournamentID = source["TournamentID"];
	        this.EntrantID = source["EntrantID"];
	        this.Status = source["Status"];
	        this.Score = source["Score"];
	    }
	}
	export class UpdateMatchSingleParticipantParams {
	    TournamentID: string;
	    MatchID: string;
	    ParticipantID: string;
	    Score?: number;
	    WinnerID?: string;
	    Status?: string;
	    AppendScore: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UpdateMatchSingleParticipantParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TournamentID = source["TournamentID"];
	        this.MatchID = source["MatchID"];
	        this.ParticipantID = source["ParticipantID"];
	        this.Score = source["Score"];
	        this.WinnerID = source["WinnerID"];
	        this.Status = source["Status"];
	        this.AppendScore = source["AppendScore"];
	    }
	}
	export class UpdateTournamentRequest {
	    name?: string;
	    size?: number;
	    startingDate?: string;
	    track?: string;
	    qualifyingOrder?: string[];
	
	    static createFrom(source: any = {}) {
	        return new UpdateTournamentRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.size = source["size"];
	        this.startingDate = source["startingDate"];
	        this.track = source["track"];
	        this.qualifyingOrder = source["qualifyingOrder"];
	    }
	}

}

export namespace tournament {
	
	export class Tournament {
	    id: string;
	    name: string;
	    starting_date: string;
	    track?: string;
	    qualifying_order: string;
	    status: string;
	    url: string;
	    size: number;
	    progress: number;
	
	    static createFrom(source: any = {}) {
	        return new Tournament(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.starting_date = source["starting_date"];
	        this.track = source["track"];
	        this.qualifying_order = source["qualifying_order"];
	        this.status = source["status"];
	        this.url = source["url"];
	        this.size = source["size"];
	        this.progress = source["progress"];
	    }
	}

}

