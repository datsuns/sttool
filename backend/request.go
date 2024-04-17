package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/browser"
)

func issueRequest(r *http.Request, debug bool) ([]byte, int, error) {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		logger.Error("issueRequest::http.DefaultClient.Do", slog.Any("ERR", err.Error()))
		return nil, 0, err
	}
	defer resp.Body.Close()

	byteArray, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("issueRequest::io.ReadAll", slog.Any("ERR", err.Error()))
		return nil, 0, err
	}
	if debug {
		logger.Info("issueRequest", slog.Any("Status", resp.Status), slog.Any("URL", r.URL), slog.Any("RawRet", string(byteArray)))
	}
	switch resp.StatusCode {
	case 200:
	case 202:
	default:
		return nil, resp.StatusCode, fmt.Errorf("error responce. status[%v] msg[%v]", resp.StatusCode, string(byteArray))
	}
	return byteArray, resp.StatusCode, nil
}

func issueEventSubRequest(cfg *Config, method, url string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logger.Error("issueEventSubRequest::http.NewRequest", slog.Any("ERR", err.Error()))
		return nil, 0, err
	}
	if cfg.IsDebug() {
		logger.Info("rest auth", slog.Any("Auth", cfg.AuthCode()), slog.Any("ClientID", cfg.ClientId()))
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.AuthCode()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", cfg.ClientId())

	return issueRequest(req, cfg.IsDebug())
}

// https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/#authorization-code-grant-flow
func StartAuthorizationCodeGrantFlow(cfg *Config, redirectUrl string, scope []string) error {
	url := fmt.Sprintf(
		"https://id.twitch.tv/oauth2/authorize?client_id=%v&force_verify=true&redirect_uri=%v&response_type=code&scope=%v",
		cfg.ClientId(),
		redirectUrl,
		strings.Join(scope, " "),
	)
	return browser.OpenURL(url)
}

func createEventSubscription(cfg *Config, r *Responce, event string, e *EventTableEntry) error {
	bin := e.Builder(cfg, r, event, e.Version)
	logger.Info("create EventSub",
		slog.Any("SessionID", r.Payload.Session.Id),
		slog.Any("User", cfg.TargetUserId),
		slog.Any("Type", event),
		slog.Any("Raw", string(bin)),
	)
	_, _, err := issueEventSubRequest(cfg, "POST", "https://api.twitch.tv/helix/eventsub/subscriptions", bytes.NewReader(bin))
	return err
}

func referTargetUserIdWith(cfg *Config, username string) (string, string, string, int, error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/users?login=%v", username)
	ret, status, err := issueEventSubRequest(cfg, "GET", url, nil)
	if err != nil {
		logger.Error("Eventsub Request", slog.Any("ERR", err.Error()))
		return "", "", "", status, err
	}
	r := &GetUsersApiResponce{}
	err = json.Unmarshal(ret, &r)
	if err != nil {
		logger.Error("json.Unmarshal", slog.Any("ERR", err.Error()))
		return "", "", "", status, err
	}
	logger.Info("referUserId", slog.Any("id", r.Data[0].Id), slog.Any("name", r.Data[0].DisplayName))
	return r.Data[0].Id, r.Data[0].Login, r.Data[0].DisplayName, status, nil
}

func RequestUserAccessToken(cfg *Config, code, redirectUri string) (string, string, error) {
	var err error
	params := url.Values{}
	params.Add("Content-Type", "application/x-www-form-urlencoded")
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", cfg.ClientId())
	params.Add("client_secret", cfg.ClientSecret())
	params.Add("code", code)
	params.Add("redirect_uri", redirectUri)

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", bytes.NewBufferString(params.Encode()))
	if err != nil {
		logger.Error("RequestUserAccessToken::http.NewRequest", slog.Any("ERR", err.Error()))
		return "", "", err
	}

	byteArray, _, err := issueRequest(req, cfg.IsDebug())
	if err != nil {
		return "", "", err
	}

	r := &RequestTokenByCodeResponce{}
	err = json.Unmarshal(byteArray, &r)
	if err != nil {
		logger.Error("json.Unmarshal", slog.Any("ERR", err.Error()))
		return "", "", nil
	}
	return r.AccessToken, r.RefreshToken, nil
}

func RefreshAccessToken(cfg *Config, refreshToken string) (string, string, error) {
	var err error
	params := url.Values{}
	params.Add("Content-Type", "application/x-www-form-urlencoded")
	params.Add("client_id", cfg.ClientId())
	params.Add("client_secret", cfg.ClientSecret())
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", bytes.NewBufferString(params.Encode()))
	if err != nil {
		logger.Error("RefreshAccessToken::http.NewRequest", slog.Any("ERR", err.Error()))
		return "", "", err
	}

	byteArray, _, err := issueRequest(req, cfg.IsDebug())
	if err != nil {
		return "", "", err
	}

	r := &RefreshTokenResponce{}
	err = json.Unmarshal(byteArray, &r)
	if err != nil {
		logger.Error("json.Unmarshal", slog.Any("ERR", err.Error()))
		return "", "", nil
	}
	return r.AccessToken, r.RefreshToken, nil
}

func ValidateAccessToken(cfg *Config) (bool, string, string, error) {
	req, err := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)
	if err != nil {
		logger.Error("ValidateAccessToken::http.NewRequest", slog.Any("ERR", err.Error()))
		return false, "", "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", cfg.AuthCode()))

	byteArray, statusCode, err := issueRequest(req, cfg.IsDebug())
	if statusCode == 401 {
		return false, "", "", nil
	}
	if err != nil {
		return false, "", "", err
	}

	r := &ValidateTokenResponce{}
	err = json.Unmarshal(byteArray, &r)
	if err != nil {
		logger.Error("json.Unmarshal", slog.Any("ERR", err.Error()))
		return false, "", "", err
	}
	return statusCode == 200, r.Login, r.UserId, nil
}

func ReferTargetUserId(cfg *Config) (string, int, error) {
	id, _, _, status, err := referTargetUserIdWith(cfg, cfg.UserName())
	if err != nil {
		return "", status, err
	}
	return id, status, nil
}

func ReferTargetUser(cfg *Config) (string, string, string, int, error) {
	return referTargetUserIdWith(cfg, cfg.UserName())
}

func ReferUserClips(cfg *Config, userId string) (string, *GetClipsApiResponce) {
	return referUserClipsByDate(cfg, userId, true, nil)
}

func issueGetClipRequest(cfg *Config, url string) (string, *GetClipsApiResponce) {
	raw, _, err := issueEventSubRequest(cfg, "GET", url, nil)
	if err != nil {
		logger.Error("Eventsub Request", slog.Any("ERR", err.Error()))
		return "", nil
	}

	r := &GetClipsApiResponce{}
	err = json.Unmarshal(raw, &r)
	if err != nil {
		logger.Error("json.Unmarshal", slog.Any("ERR", err.Error()))
		return "", nil
	}
	ret := ""
	for _, v := range r.Data {
		//infoLogger.Info("UserClip", slog.Any("タイトル", v.Title), slog.Any("URL", v.Url), slog.Any("視聴回数", v.ViewCount))
		ret += fmt.Sprintf("   再生回数[%v] / タイトル[%v] / URL[ %v ] / Id[ %v ]\n", v.ViewCount, v.Title, v.Url, v.Id)
	}
	return ret, r
}

func referUserClipsByDate(cfg *Config, userId string, featured bool, date *time.Time) (text string, ret *GetClipsApiResponce) {
	maxN := 4
	url := fmt.Sprintf("https://api.twitch.tv/helix/clips?broadcaster_id=%v&is_featured=%v&first=%v", userId, featured, maxN)
	if date != nil {
		url += fmt.Sprintf("&started_at=%v", date.UTC().Format(time.RFC3339))
	}

	text, ret = issueGetClipRequest(cfg, url)
	if len(ret.Data) > 0 {
		return text, ret
	}
	url = fmt.Sprintf("https://api.twitch.tv/helix/clips?broadcaster_id=%v&is_featured=%v&first=%v", userId, false, maxN)
	if date != nil {
		url += fmt.Sprintf("&started_at=%v", date.UTC().Format(time.RFC3339))
	}
	return issueGetClipRequest(cfg, url)
}
