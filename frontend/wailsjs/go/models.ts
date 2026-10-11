export namespace repository {
	
	export class CreateTournamentRequest {
	    name: string;
	    size: number;
	    startingDate: string;
	    track: string;
	    entrants: string[];
	    randomizeQualifyingOrder: boolean;
	
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

