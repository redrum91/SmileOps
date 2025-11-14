export namespace excel {
	
	export class Operation {
	    dates: string[];
	    numbers: string[];
	    comment: string;
	
	    static createFrom(source: any = {}) {
	        return new Operation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dates = source["dates"];
	        this.numbers = source["numbers"];
	        this.comment = source["comment"];
	    }
	}
	export class Patient {
	    id: string;
	    fio: string;
	    implantNumber: number;
	    operations: Record<string, Operation>;
	    controlHalfYear: string;
	    controlYear: string;
	    occupationalHygiene: string;
	
	    static createFrom(source: any = {}) {
	        return new Patient(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.fio = source["fio"];
	        this.implantNumber = source["implantNumber"];
	        this.operations = this.convertValues(source["operations"], Operation, true);
	        this.controlHalfYear = source["controlHalfYear"];
	        this.controlYear = source["controlYear"];
	        this.occupationalHygiene = source["occupationalHygiene"];
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

