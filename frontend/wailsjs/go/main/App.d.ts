// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';
import {backend} from '../models';

export function DebugAppendEntry():Promise<void>;

export function DebugRaidTest(arg1:string):Promise<void>;

export function LoadConfig():Promise<main.AppConfig>;

export function OnConnectedCallback():Promise<void>;

export function OnKeepAliveCallback():Promise<void>;

export function OnRaidCallback(arg1:backend.RaidCallbackParam):Promise<void>;

export function OpenDiectoryDialog(arg1:string):Promise<string>;

export function OpenFileDialog(arg1:string,arg2:string):Promise<string>;

export function OpenURL(arg1:string):Promise<void>;

export function SaveConfig(arg1:main.AppConfig):Promise<void>;

export function StartClip(arg1:string,arg2:number):Promise<void>;

export function StopClip():Promise<void>;

export function StopObsStream():Promise<void>;

export function TestObsConnection():Promise<string>;
