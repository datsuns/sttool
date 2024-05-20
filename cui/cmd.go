package main

import "sttool/backend"

func OnKeepAliveCallback() {
}

func OnRaidCallback(_ *backend.RaidCallbackParam) {
}

func OnConnectedCallback() {
}

func main() {
	callback := &backend.CallBack{
		KeepAlive:   OnKeepAliveCallback,
		OnRaid:      OnRaidCallback,
		OnConnected: OnConnectedCallback,
	}
	b := backend.NewBackend(callback)
	b.Serve()
}
