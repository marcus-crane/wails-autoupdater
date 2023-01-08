export namespace main {
	
	export class UpdateStatus {
	    update_available: boolean;
	    remote_version: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.update_available = source["update_available"];
	        this.remote_version = source["remote_version"];
	    }
	}

}

