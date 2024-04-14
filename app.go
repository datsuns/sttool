package main

import (
	"context"
	"fmt"
	"sttool/backend"

	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	Items   []backend.UserClip
	Backend *backend.BackendContext
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.Items = []backend.UserClip{}
	callback := &backend.CallBack{
		KeepAlive: a.OnKeepAliveCallback,
		OnRaid:    a.OnRaidCallback,
	}
	a.Backend = backend.NewBackend(callback)
	go a.Backend.Serve()
}

func (a *App) StartClip(id string, duration float32) {
	a.Backend.Overlay.StartClip(id, duration)
}

func (a *App) StopClip() {
	a.Backend.Overlay.StopClip()
}

func (a *App) GetServerPort() int {
	return a.Backend.GetOverlayPortNumber()
}

func (a *App) OpenURL(url string) {
	if err := browser.OpenURL(url); err != nil {
		runtime.LogDebug(a.ctx, fmt.Sprintf("URL[%v] open error: %v", url, err))
	}
}

func (a *App) OnKeepAliveCallback() {
	//runtime.LogDebug(a.ctx, "KeepAlive")
	//runtime.EventsEmit(a.ctx, "testevent", "event from backend", a.Items)
}

func (a *App) OnRaidCallback(param *backend.RaidCallbackParam) {
	runtime.LogDebug(a.ctx,
		fmt.Sprintf("OnRaidCallback from[%v]", param.From),
	)
	runtime.EventsEmit(a.ctx, "OnRaid", "raided users clip", param.From, param.Clips)
	//runtime.EventsEmit(a.ctx, "OnRaid", "raided users clip", param.From, param.Clips)
}

func (a *App) DebugAppendEntry() {
	a.Items = append(a.Items, backend.UserClip{Url: "https://example2.com", Thumbnail: "https://example2.com/thumbnail.jpg", ViewCount: 300, Title: "Example Video 3"})
	runtime.EventsEmit(a.ctx, "testevent", "event from backend", a.Items)
	runtime.LogDebug(a.ctx, "test")
}

func (a *App) DebugRaidTest(userName string) {
	runtime.LogDebug(a.ctx, fmt.Sprintf("start DebugRaid w/ [%v]", userName))
	cfg, _ := backend.LoadConfig()
	cfg.TargetUser = userName
	id, _, diplayName, _, _ := backend.ReferTargetUser(cfg)
	runtime.LogDebug(a.ctx, fmt.Sprintf("user [%v] is [%v]", userName, id))
	_, ret := backend.ReferUserClips(cfg, id)
	data := []backend.UserClip{}
	for _, c := range ret.Data {
		data = append(data, backend.UserClip{
			Id:        c.Id,
			Url:       c.Url,
			Thumbnail: c.ThumbnailUrl,
			ViewCount: c.ViewCount,
			Title:     c.Title,
			Duration:  c.Duration,
		})
		runtime.LogDebug(a.ctx, fmt.Sprintf("found clip [%v]", c.Title))
	}
	runtime.EventsEmit(a.ctx, "OnRaid", "raided users clip", diplayName, data)
}
