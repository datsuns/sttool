package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sttool/backend"

	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	Items   []backend.UserClip
	Backend *backend.BackendContext
}

type AppConfig struct {
	backend.ConfigBody
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.Items = []backend.UserClip{}
	callback := &backend.CallBack{
		KeepAlive:   a.OnKeepAliveCallback,
		OnRaid:      a.OnRaidCallback,
		OnConnected: a.OnConnectedCallback,
	}
	a.Backend = backend.NewBackend(callback)
	go a.Backend.Serve()
}

func (a *App) StartClip(url string, duration float32) {
	a.Backend.Overlay.StartClip(url, duration)
}

func (a *App) StopClip() {
	a.Backend.Overlay.StopClip()
}

func (a *App) LoadConfig() *AppConfig {
	ret := &AppConfig{ConfigBody: *a.Backend.LoadConfig()}
	return ret
}

func (a *App) SaveConfig(appCfg *AppConfig) {
	a.Backend.SaveConfig(&appCfg.ConfigBody)
}

func (a *App) OpenURL(url string) {
	if err := browser.OpenURL(url); err != nil {
		runtime.LogDebug(a.ctx, fmt.Sprintf("URL[%v] open error: %v", url, err))
	}
}

func (a *App) StopObsStream() {
	a.Backend.StopObsStream()
}

func (a *App) OpenFileDialog(prev, filter string) string {
	filters := []runtime.FileFilter{}
	for _, f := range strings.Split(filter, ",") {
		filters = append(filters, runtime.FileFilter{Pattern: f})
	}
	opt := runtime.OpenDialogOptions{
		DefaultDirectory: filepath.Dir(prev),
		Title:            "ファイル選択",
		Filters:          filters,
	}
	selected, err := runtime.OpenFileDialog(a.ctx, opt)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("OpenFileDialog:ERR %v", err.Error()))
		return ""
	}
	return selected
}

func (a *App) OpenDiectoryDialog(prev string) string {
	opt := runtime.OpenDialogOptions{
		Title:            "ファルダ選択",
		DefaultDirectory: filepath.Dir(prev),
	}
	selected, err := runtime.OpenDirectoryDialog(a.ctx, opt)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("OpenDiectoryDialog:ERR %v", err.Error()))
		return ""
	}
	return selected
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

func (a *App) OnConnectedCallback() {
	runtime.LogDebug(a.ctx, "Connected")
	runtime.EventsEmit(a.ctx, "OnConnected", "connected")
}

func (a *App) DebugAppendEntry() {
	a.Items = append(a.Items, backend.UserClip{Url: "https://example2.com", Thumbnail: "https://example2.com/thumbnail.jpg", ViewCount: 300, Title: "Example Video 3"})
	runtime.EventsEmit(a.ctx, "testevent", "event from backend", a.Items)
	runtime.LogDebug(a.ctx, "test")
}

func (a *App) DebugRaidTest(userName string) {
	runtime.LogDebug(a.ctx, fmt.Sprintf("start DebugRaid w/ [%v]", userName))
	cfg, _ := backend.LoadConfig()
	cfg.LoadAuthConfig()
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
			Mp4:       backend.ConvertThumbnailToMp4Url(c.ThumbnailUrl),
		})
		runtime.LogDebug(a.ctx, fmt.Sprintf("found clip [%v]", c.Title))
	}
	runtime.EventsEmit(a.ctx, "OnRaid", "raided users clip", diplayName, data)
}
