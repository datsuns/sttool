package backend

import (
	"bytes"
	"encoding/json"
	"errors"
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
	case 401:
		logger.Error("issueRequest", slog.Any("msg", "401 error"), slog.Any("Status", resp.Status), slog.Any("URL", r.URL), slog.Any("RawRet", string(byteArray)))
		return nil, resp.StatusCode, errors.New(RequestErrorBy401)
	default:
		logger.Error("issueRequest", slog.Any("msg", "unexpected status"), slog.Any("Status", resp.Status), slog.Any("URL", r.URL), slog.Any("RawRet", string(byteArray)))
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

func issueGetClipRequest(cfg *Config, url string) (string, *GetClipsApiResponce, error) {
	raw, _, err := issueEventSubRequest(cfg, "GET", url, nil)
	if err != nil {
		logger.Error("Eventsub Request", slog.Any("ERR", err.Error()))
		return "", nil, err
	}

	r := &GetClipsApiResponce{}
	err = json.Unmarshal(raw, &r)
	if err != nil {
		logger.Error("json.Unmarshal", slog.Any("ERR", err.Error()))
		return "", nil, err
	}
	ret := ""
	for _, v := range r.Data {
		//infoLogger.Info("UserClip", slog.Any("タイトル", v.Title), slog.Any("URL", v.Url), slog.Any("視聴回数", v.ViewCount))
		ret += fmt.Sprintf("   再生回数[%v] / タイトル[%v] / URL[ %v ] / Id[ %v ]\n", v.ViewCount, v.Title, v.Url, v.Id)
	}
	return ret, r, nil
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

func CreateEventSubscription(cfg *Config, sessionID, event string, e *EventTableEntry) error {
	bin := e.Builder(cfg, sessionID, event, e.Version)
	logger.Info("create EventSub",
		slog.Any("SessionID", sessionID),
		slog.Any("User", cfg.TargetUserId),
		slog.Any("Type", event),
		slog.Any("Raw", string(bin)),
	)
	endpoint := "https://api.twitch.tv/helix/eventsub/subscriptions"
	if cfg.IsLocalTest() {
		endpoint = fmt.Sprintf("http://%v/eventsub/subscriptions", LocalTestAddr)
	}
	_, _, err := issueEventSubRequest(cfg, "POST", endpoint, bytes.NewReader(bin))
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

// https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/#use-the-authorization-code-to-get-a-token
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

// https://dev.twitch.tv/docs/authentication/refresh-tokens/
// "expires_in"を見てタイミングを測るのもいいけどほかの理由でもinvalidになる可能性あるので
// 401応答をハンドリングするほうがいいよ、とのこと
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
		logger.Error("RefreshAccessToken::issueRequest", slog.Any("ERR", err.Error()))
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

// https://dev.twitch.tv/docs/authentication/validate-tokens/
func ValidateAccessToken(cfg *Config) (bool, int, string, string, error) {
	req, err := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)
	if err != nil {
		logger.Error("ValidateAccessToken::http.NewRequest", slog.Any("ERR", err.Error()))
		return false, 0, "", "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", cfg.AuthCode()))

	byteArray, statusCode, err := issueRequest(req, cfg.IsDebug())
	if statusCode == 401 {
		logger.Info("ValidateAccessToken::401", slog.Any("raw", string(byteArray)))
		return false, 0, "", "", nil
	}
	if err != nil {
		logger.Error("ValidateAccessToken::Error", slog.Any("ERR", err.Error()))
		return false, 0, "", "", err
	}

	r := &ValidateTokenResponce{}
	err = json.Unmarshal(byteArray, &r)
	if err != nil {
		logger.Error("json.Unmarshal", slog.Any("ERR", err.Error()))
		return false, 0, "", "", err
	}
	if cfg.IsDebug() {
		logger.Info("ValidateAccessToken", slog.Any("raw", r))
	}
	return statusCode == 200, r.ExpiresIn, r.Login, r.UserId, nil
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

func ReferUserClips(cfg *Config, userId string) (string, *GetClipsApiResponce, error) {
	return ReferUserClipsByDate(cfg, userId, true, nil)
}

func ReferUserClipsByDate(cfg *Config, userId string, featured bool, date *time.Time) (text string, ret *GetClipsApiResponce, err error) {
	maxN := 4
	url := fmt.Sprintf("https://api.twitch.tv/helix/clips?broadcaster_id=%v&is_featured=%v&first=%v", userId, featured, maxN)
	if date != nil {
		url += fmt.Sprintf("&started_at=%v", date.UTC().Format(time.RFC3339))
	}

	text, ret, err = issueGetClipRequest(cfg, url)
	if err != nil {
		return "", nil, err
	}
	if len(ret.Data) > 0 {
		return text, ret, nil
	}
	url = fmt.Sprintf("https://api.twitch.tv/helix/clips?broadcaster_id=%v&is_featured=%v&first=%v", userId, false, maxN)
	if date != nil {
		url += fmt.Sprintf("&started_at=%v", date.UTC().Format(time.RFC3339))
	}
	return issueGetClipRequest(cfg, url)
}

func ReferUserChannelRewards(cfg *Config, userId string) (*GetCustomRewardResponce, error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/channel_points/custom_rewards?broadcaster_id=%v", userId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("issueEventSubRequest::http.NewRequest", slog.Any("ERR", err.Error()))
		return nil, err
	}
	//logger.Info("rest auth", slog.Any("Auth", cfg.AuthCode()), slog.Any("ClientID", cfg.ClientId()))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.AuthCode()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", cfg.ClientId())

	raw, _, err := issueRequest(req, cfg.IsDebug())
	if err != nil {
		return nil, err
	}

	r := &GetCustomRewardResponce{}
	err = json.Unmarshal(raw, &r)
	if err != nil {
		logger.Error("json.Unmarshal", slog.Any("ERR", err.Error()))
		return nil, nil
	}
	return r, nil
}
