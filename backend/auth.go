package backend

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func redirectHandler(w http.ResponseWriter, r *http.Request, fin chan struct{}) string {
	v := r.URL.Query()
	//logger.Info("redirectHandler", slog.Any("msg", "occur"), slog.Any("code", v["code"]))
	fmt.Fprintf(w, "ブラウザを閉じてください")
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	fin <- struct{}{}
	return v.Get("code")
}

func Issue1stTimeAuthentication(cfg *Config) error {
	redirectUri := "http://localhost"
	fin := make(chan struct{})
	var code string
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code = redirectHandler(w, r, fin)
	})

	go func() {
		if err := http.ListenAndServe(
			"",
			nil,
		); err != nil {
			logger.Error("Ovelay:ERROR", slog.Any("error", err.Error))
		}
	}()
	time.Sleep(time.Second)

	if e := StartAuthorizationCodeGrantFlow(cfg, redirectUri, []string{
		"bits:read",
		"channel:read:subscriptions",
		"channel:read:redemptions",
		"channel:manage:raids",
		"moderator:read:followers",
		"user:read:chat",
	}); e != nil {
		return e
	}
	<-fin
	a, r, _ := RequestUserAccessToken(cfg, code, redirectUri)
	cfg.UpdatUserAccessToken(a)
	cfg.UpdatRefreshToken(r)
	cfg.Save()
	logger.Info("Issue1stTimeAuthentication", slog.Any("code", code), slog.Any("access", a), slog.Any("refresh", r))
	return fmt.Errorf("Issue1stTimeAuthentication")
}

func ConfirmAccessToken(cfg *Config) error {
	var err error
	if e := cfg.LoadAuthConfig(); e != nil {
		return Issue1stTimeAuthentication(cfg)
	}
	expired, err := ConfirmUserAccessToken(cfg)
	if err != nil {
		logger.Error("ConfirmUserAccessToken", slog.Any("ERR", err.Error()))
		return err
	}
	if expired == false {
		return nil
	}
	logger.Warn("ConfirmUserAccessToken", slog.Any("msg", "need to refresh token"))
	token, _, err := RequestUserAccessToken(cfg, "", "")
	if err != nil {
		return err
	}
	logger.Warn("ConfirmUserAccessToken", slog.Any("msg", "token refreshed"))
	cfg.UpdatUserAccessToken(token)
	cfg.Save()
	return nil
}

func ConfirmUserAccessToken(cfg *Config) (bool, error) {
	if _, err := os.Stat(AuthInfoFile); err != nil {
		logger.Warn("ConfirmUserAccessToken", slog.Any("msg", "not authlized. start to connect"))
		return true, nil
	}
	_, status, err := ReferTargetUserId(cfg)
	if err != nil {
		if status == 401 {
			logger.Warn("ConfirmUserAccessToken", slog.Any("msg", "token expired"))
			return true, nil
		}
		logger.Error("ConfirmUserAccessToken", slog.Any("status", status))
		return true, err
	}
	return false, nil
}
