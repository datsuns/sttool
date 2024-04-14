package backend

import (
	"fmt"
	"log/slog"
	"net/http"
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

	if e := StartAuthorizationCodeGrantFlow(cfg, AuthRedirectUri, EventSubScope); e != nil {
		return e
	}
	<-fin
	a, r, _ := RequestUserAccessToken(cfg, code, AuthRedirectUri)
	cfg.UpdatAccessToken(AuthEntry{AuthCode: a, RefreshToken: r})
	cfg.Save()
	logger.Info("Issue1stTimeAuthentication", slog.Any("code", code), slog.Any("access", a), slog.Any("refresh", r))
	return nil
}

func ConfirmAccessToken(cfg *Config) error {
	var err error
	if e := cfg.LoadAuthConfig(); e != nil {
		if e := Issue1stTimeAuthentication(cfg); e != nil {
			return e
		}
	}
	valid, name, id, err := confirmUserAccessToken(cfg)
	if err != nil {
		logger.Error("ConfirmUserAccessToken", slog.Any("ERR", err.Error()))
		return err
	}
	cfg.TargetUserId = id
	cfg.TargetUser = name
	if err != nil {
		logger.Error("ConfirmUserAccessToken::ReferTargetUserId", slog.Any("ERR", err.Error()))
		return err
	}

	if valid {
		return nil
	}

	logger.Info("ConfirmUserAccessToken", slog.Any("msg", "start token refresh"))
	a, r, err := RefreshAccessToken(cfg, cfg.RefreshToken())
	if err != nil {
		return err
	}
	cfg.UpdatAccessToken(AuthEntry{AuthCode: a, RefreshToken: r})
	cfg.Save()
	return nil
}

func confirmUserAccessToken(cfg *Config) (bool, string, string, error) {
	valid, name, id, err := ValidateAccessToken(cfg)
	if err != nil {
		logger.Error("ConfirmUserAccessToken", slog.Any("ERR", err.Error()))
		return false, "", "", err
	}
	logger.Info("ConfirmUserAccessToken", slog.Any("valid", valid))
	return valid, name, id, nil
}
