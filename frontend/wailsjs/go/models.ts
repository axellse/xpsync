export namespace main {
	
	export class Action {
	    Type: string;
	    FileName: string;
	    Data: number[];
	
	    static createFrom(source: any = {}) {
	        return new Action(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Type = source["Type"];
	        this.FileName = source["FileName"];
	        this.Data = source["Data"];
	    }
	}
	export class Device {
	    Ip: string;
	    Name: string;
	    LastPing: number;
	    PendingAction: Action;
	
	    static createFrom(source: any = {}) {
	        return new Device(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Ip = source["Ip"];
	        this.Name = source["Name"];
	        this.LastPing = source["LastPing"];
	        this.PendingAction = this.convertValues(source["PendingAction"], Action);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

