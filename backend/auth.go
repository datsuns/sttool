package backend

import (
	"context"
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
	defer close(fin)
	var code string

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code = redirectHandler(w, r, fin)
	})
	server := &http.Server{
		Addr:    "",
		Handler: mux,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Issue1stTimeAuthentication::start", slog.Any("error", err.Error))
			fin <- struct{}{}
		}
	}()
	time.Sleep(time.Second)

	if e := StartAuthorizationCodeGrantFlow(cfg, AuthRedirectUri, EventSubScope); e != nil {
		return e
	}
	<-fin
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Issue1stTimeAuthentication::shutdown", slog.Any("error", err.Error))
	}

	a, r, _ := RequestUserAccessToken(cfg, code, AuthRedirectUri)
	cfg.UpdatAccessToken(AuthEntry{AuthCode: a, RefreshToken: r})
	cfg.SaveAll()
	logger.Info("Issue1stTimeAuthentication", slog.Any("code", code), slog.Any("access", a), slog.Any("refresh", r))
	statsLogger.Info("Issue1stTimeAuthentication",
		slog.Any(LogFieldName_Type, "ResetToken"),
		slog.Any("code", code),
	)
	return nil
}

func ConfirmAccessToken(cfg *Config) error {
	if e := cfg.LoadAuthConfig(); e != nil {
		logger.Info("ConfirmAccessToken", slog.Any("msg", "LoadAuthConfig error. try to 1st auth"), slog.Any("ERR", e.Error()))
		if e := Issue1stTimeAuthentication(cfg); e != nil {
			return e
		}
	}
	valid, expires, name, id, err := confirmUserAccessToken(cfg)
	if err != nil {
		logger.Error("ConfirmUserAccessToken", slog.Any("ERR", err.Error()))
		return err
	}

	if valid {
		cfg.TargetUserId = id
		cfg.TargetUser = name
		logger.Info("ConfirmUserAccessToken", slog.Any("msg", "ok"), slog.Any("expired", expires))
		return nil
	}

	logger.Info("ConfirmUserAccessToken", slog.Any("msg", "start token refresh"))
	a, r, err := RefreshAccessToken(cfg, cfg.RefreshToken())
	if err != nil {
		return err
	}
	cfg.UpdatAccessToken(AuthEntry{AuthCode: a, RefreshToken: r})
	cfg.SaveAll()
	_, expires, name, id, _ = confirmUserAccessToken(cfg)
	cfg.TargetUserId = id
	cfg.TargetUser = name
	logger.Info("ConfirmUserAccessToken", slog.Any("msg", "token refreshed"), slog.Any("expired", expires))
	return nil
}

func confirmUserAccessToken(cfg *Config) (bool, int, string, string, error) {
	valid, expires, name, id, err := ValidateAccessToken(cfg)
	if err != nil {
		logger.Error("ConfirmUserAccessToken", slog.Any("ERR", err.Error()))
		return false, 0, "", "", err
	}
	logger.Info("ConfirmUserAccessToken", slog.Any("valid", valid))
	return valid, expires, name, id, nil
}
