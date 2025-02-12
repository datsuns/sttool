// session_reconnectの制御や401のハンドリングを増やしたので不要になってるはず(しばし残しておく)
package backend

import (
	"log/slog"
	"time"
)

const (
	RefreshRetryInterval = 60
	MinumumExpiredSecond = 120
)

type SessionLifecycle struct {
	Expired           int
	Config            *Config
	Refreshed         *chan struct{}
	Cancel            chan struct{}
	localTestInterval int
}

func NewSessionLifecycle(refreshd *chan struct{}, cfg *Config) *SessionLifecycle {
	return &SessionLifecycle{
		Expired:           0,
		Config:            cfg,
		Refreshed:         refreshd,
		Cancel:            make(chan struct{}),
		localTestInterval: 15,
	}
}

func (sl *SessionLifecycle) Serve(firstExpires int) {
	if firstExpires > MinumumExpiredSecond {
		sl.Expired = firstExpires - MinumumExpiredSecond
	} else {
		sl.Expired, _ = sl.refresh()
	}
	//sl.Expired = sl.localTestInterval
	go func() { sl.watch() }()
}

func (sl *SessionLifecycle) Shutdown() {
	sl.Cancel <- struct{}{}
	// TODO syncronize w/ watch() finished
}

func (sl *SessionLifecycle) watch() {
	var err error
	for {
		select {
		case <-time.After(time.Duration(sl.Expired) * time.Second):
			sl.Expired, err = sl.refresh()
			if err != nil {
				logger.Error("SessionLifecycle::Serve", slog.Any("msg", "token refresh error"), slog.Any("ERR", err.Error()))
				sl.Expired = RefreshRetryInterval
			} else {
				*(sl.Refreshed) <- struct{}{}
			}
		case <-sl.Cancel:
			return
		}
	}
}

func (sl *SessionLifecycle) refresh() (int, error) {
	a, r, err := RefreshAccessToken(sl.Config, sl.Config.RefreshToken())
	if err != nil {
		logger.Error("SessionLifecycle::refresh", slog.Any("msg", "RefreshAccessToken"), slog.Any("ERR", err.Error()))
		return 0, err
	}
	expires, err := UpdateSavedRefreshToken(sl.Config, a, r)
	logger.Info("SessionLifecycle::refresh", slog.Any("msg", "token refreshed"), slog.Any("expired", expires))
	if expires > MinumumExpiredSecond {
		expires -= MinumumExpiredSecond
	}
	return expires, nil
	//return sl.localTestInterval, nil
}
