export namespace backend {
	
	export class UserClip {
	    Id: string;
	    Url: string;
	    Title: string;
	    Thumbnail: string;
	    ViewCount: number;
	    Duration: number;
	    Mp4: string;
	
	    static createFrom(source: any = {}) {
	        return new UserClip(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Url = source["Url"];
	        this.Title = source["Title"];
	        this.Thumbnail = source["Thumbnail"];
	        this.ViewCount = source["ViewCount"];
	        this.Duration = source["Duration"];
	        this.Mp4 = source["Mp4"];
	    }
	}
	export class RaidCallbackParam {
	    From: string;
	    Clips: UserClip[];
	
	    static createFrom(source: any = {}) {
	        return new RaidCallbackParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.From = source["From"];
	        this.Clips = this.convertValues(source["Clips"], UserClip);
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

export namespace main {
	
	export class AppConfig {
	    ChatTargets: string[];
	    NotifySoundFile: string;
	    DebugMode: boolean;
	    LocalTest: boolean;
	    LogDest: string;
	    ObsIp: string;
	    ObsPort: number;
	    ObsPass: string;
	    StopStreamAfterRaided: boolean;
	    DelaySecondsFromRaidToStop: number;
	    NewClipWatchIntervalSecond: number;
	    LocalServerPortNumber: number;
	    OverlayEnabled: boolean;
	    ClipPlayerWidth: number;
	    ClipPlayerHeight: number;
	    LogTopIndent: string;
	    LogUserNamePrefix: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ChatTargets = source["ChatTargets"];
	        this.NotifySoundFile = source["NotifySoundFile"];
	        this.DebugMode = source["DebugMode"];
	        this.LocalTest = source["LocalTest"];
	        this.LogDest = source["LogDest"];
	        this.ObsIp = source["ObsIp"];
	        this.ObsPort = source["ObsPort"];
	        this.ObsPass = source["ObsPass"];
	        this.StopStreamAfterRaided = source["StopStreamAfterRaided"];
	        this.DelaySecondsFromRaidToStop = source["DelaySecondsFromRaidToStop"];
	        this.NewClipWatchIntervalSecond = source["NewClipWatchIntervalSecond"];
	        this.LocalServerPortNumber = source["LocalServerPortNumber"];
	        this.OverlayEnabled = source["OverlayEnabled"];
	        this.ClipPlayerWidth = source["ClipPlayerWidth"];
	        this.ClipPlayerHeight = source["ClipPlayerHeight"];
	        this.LogTopIndent = source["LogTopIndent"];
	        this.LogUserNamePrefix = source["LogUserNamePrefix"];
	    }
	}

}

