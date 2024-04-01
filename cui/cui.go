package main

import "sttool/backend"

func OnKeepAliveCallback() {
}

func OnRaidCallback(param *backend.RaidCallbackParam) {
}

func main() {
	callback := &backend.CallBack{
		KeepAlive: OnKeepAliveCallback,
		OnRaid:    OnRaidCallback,
	}
	backend.Serve(callback)
}
