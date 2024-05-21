package main

import (
	"fmt"
	"sttool/backend"
)

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
	backend.ConfirmAccessToken(b.Config)
	ret := b.ListChannelRewards()
	for _, p := range ret {
		fmt.Printf("title[%v] id[%v] enable[%v] paused[%v]\n", p.Title, p.Id, p.Enabled, p.Paused)
	}
}
