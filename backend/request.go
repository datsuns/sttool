package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func issueEventSubRequest(cfg *Config, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if cfg.IsDebug() {
		logger.Info("rest auth", "Auth", cfg.AuthCode(), "ClientID", cfg.ClientId())
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.AuthCode()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", cfg.ClientId())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byteArray, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if cfg.IsDebug() {
		logger.Info("request", "Status", resp.Status, "URL", url, "RawRet", string(byteArray))
	}
	switch resp.StatusCode {
	case 200:
	case 202:
	default:
		return nil, fmt.Errorf("error responce. status[%v] msg[%v]", resp.StatusCode, string(byteArray))
	}
	return byteArray, nil
}

func createEventSubscription(cfg *Config, r *Responce, event string, e *EventTableEntry) error {
	bin := e.Builder(cfg, r, event, e.Version)
	logger.Info("create EventSub", "SessionID", r.Payload.Session.Id, "User", cfg.TargetUserId, "Type", event, "Raw", string(bin))
	_, err := issueEventSubRequest(cfg, "POST", "https://api.twitch.tv/helix/eventsub/subscriptions", bytes.NewReader(bin))
	return err
}

func referTargetUserIdWith(cfg *Config, username string) (string, string, string, error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/users?login=%v", username)
	ret, err := issueEventSubRequest(cfg, "GET", url, nil)
	if err != nil {
		logger.Error("Eventsub Request", slog.Any("ERR", err.Error()))
		return "", "", "", err
	}
	r := &GetUsersApiResponce{}
	err = json.Unmarshal(ret, &r)
	if err != nil {
		logger.Error("json.Unmarshal", slog.Any("ERR", err.Error()))
		return "", "", "", err
	}
	logger.Info("referUserId", "id", r.Data[0].Id, "name", r.Data[0].DisplayName)
	return r.Data[0].Id, r.Data[0].Login, r.Data[0].DisplayName, nil
}

func ConfirmUserAccessToken(cfg *Config) error {
	return nil
}

func ReferTargetUserId(cfg *Config) (string, error) {
	id, _, _, err := referTargetUserIdWith(cfg, cfg.TargetUser())
	if err != nil {
		return "", err
	}
	return id, nil
}

func ReferTargetUser(cfg *Config) (string, string, string, error) {
	return referTargetUserIdWith(cfg, cfg.TargetUser())
}

func ReferUserClips(cfg *Config, userId string) (string, *GetClipsApiResponce) {
	return referUserClipsByDate(cfg, userId, true, nil)
}

func issueGetClipRequest(cfg *Config, url string) (string, *GetClipsApiResponce) {
	raw, err := issueEventSubRequest(cfg, "GET", url, nil)
	if err != nil {
		logger.Error("Eventsub Request", "ERROR", err.Error())
		return "", nil
	}

	r := &GetClipsApiResponce{}
	err = json.Unmarshal(raw, &r)
	if err != nil {
		logger.Error("json.Unmarshal", "ERR", err.Error())
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
