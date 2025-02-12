package backend

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func redirectHandler(w http.ResponseWriter, r *http.Request, fin chan struct{}) string {
	ret := ""
	v := r.URL.Query()
	//logger.Info("redirectHandler", slog.Any("msg", "occur"), slog.Any("code", v["code"]))
	if v.Has("error") {
		fmt.Fprintf(w, "再度認証しなおしてください")
	} else {
		fmt.Fprintf(w, "ブラウザを閉じてください")
		ret = v.Get("code")
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	fin <- struct{}{}
	return ret
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
	if code == "" {
		logger.Error("Issue1stTimeAuthentication::rejected", slog.Any("error", "1st time Auth rejected"))
		statsLogger.Error("rejected", slog.Any(LogFieldName_Type, "Error:1stTimeAuthentication"),
			slog.Any("reason", "1st time Auth rejected"),
		)
		return errors.New("1stAuthRejected")
	}

	a, r, _ := RequestUserAccessToken(cfg, code, AuthRedirectUri)
	UpdateSavedRefreshToken(cfg, a, r)
	logger.Info("Issue1stTimeAuthentication", slog.Any("code", code), slog.Any("access", a), slog.Any("refresh", r))
	statsLogger.Info("Issue1stTimeAuthentication",
		slog.Any(LogFieldName_Type, "ResetToken"),
		slog.Any("code", code),
	)
	return nil
}

func ConfirmAccessToken(cfg *Config) (int, error) {
	if e := cfg.LoadAuthConfig(); e != nil {
		logger.Info("ConfirmAccessToken", slog.Any("msg", "LoadAuthConfig error. try to 1st auth"), slog.Any("ERR", e.Error()))
		if e := Issue1stTimeAuthentication(cfg); e != nil {
			return 0, e
		}
	}
	valid, expires, name, id, err := ValidateAccessToken(cfg)
	if err != nil {
		logger.Error("ConfirmUserAccessToken", slog.Any("ERR", err.Error()))
		return 0, err
	}

	if valid {
		cfg.TargetUserId = id
		cfg.TargetUser = name
		logger.Info("ConfirmUserAccessToken", slog.Any("msg", "ok"), slog.Any("expired", expires))
		return expires, nil
	}

	logger.Info("ConfirmUserAccessToken", slog.Any("msg", "start token refresh"))
	a, r, err := RefreshAccessToken(cfg, cfg.RefreshToken())
	if err != nil {
		logger.Error("ConfirmUserAccessToken : RefreshAccessToken", slog.Any("ERR", err.Error()))
		return 0, err
	}

	expires, err = UpdateSavedRefreshToken(cfg, a, r)
	if err != nil {
		logger.Error("ConfirmUserAccessToken : UpdateSavedRefreshToken", slog.Any("ERR", err.Error()))
		return 0, err
	}
	return expires, nil
}

func UpdateSavedRefreshToken(cfg *Config, authCode string, refreshToken string) (int, error) {
	cfg.UpdatAccessToken(AuthEntry{AuthCode: authCode, RefreshToken: refreshToken})
	valid, expires, name, id, err := ValidateAccessToken(cfg)
	if err != nil {
		logger.Error("UpdateSavedRefreshToken", slog.Any("ERR", err.Error()))
		return 0, err
	}
	if !valid {
		logger.Error("UpdateSavedRefreshToken", slog.Any("msg", "token still invalid"))
		return 0, fmt.Errorf("Token still invalid")
	}
	cfg.TargetUserId = id
	cfg.TargetUser = name
	cfg.SaveAll()
	return expires, err
}
