package backend

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type CreateRequestBuilder func(*Config, string, string, string) []byte
type NotificationHandler func(*BackendContext, *Config, *Responce, []byte, *TwitchStats)

type EventTableEntry struct {
	LogTitle string
	Version  string
	Builder  CreateRequestBuilder
	Handler  NotificationHandler
}

var (
	// https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#subscription-types
	TwitchEventTable = map[string]EventTableEntry{
		"channel.subscribe":            {"サブスク", "1", buildRequest, handleNotificationChannelSubscribe}, // channel:read:subscriptions
		"channel.cheer":                {"cheer", "1", buildRequest, handleNotificationChannelCheer},    // bits:read
		"stream.online":                {"配信開始", "1", buildRequest, handleNotificationStreamOnline},
		"stream.offline":               {"配信終了", "1", buildRequest, handleNotificationStreamOffline},
		"channel.subscription.gift":    {"サブギフ", "1", buildRequest, handleNotificationChannelSubscriptionGift},       // channel:read:subscriptions
		"channel.subscription.message": {"再サブスク", "1", buildRequest, handleNotificationChannelSubscriptionMessage},   // channel:read:subscriptionsg",
		"channel.chat.notification":    {"通知", "1", buildRequestWithUser, handleNotificationChannelChatNotification}, // user:read:chat
		"channel.chat.message":         {"チャット", "1", buildRequestWithUser, handleNotificationChannelChatMessage},    // user:read:chat
		"channel.raid":                 {"レイド開始", "1", buildRequestWithFromUser, handleNotificationRaidStarted},      // none
		"channel.follow":               {"フォロー", "2", buildRequestWithModerator, handleNotificationChannelFollow},    // moderator:read:followers
		"channel.channel_points_custom_reward_redemption.add": {"チャネポ", "1", buildRequest, handleNotificationChannelPointsCustomRewardRedemptionAdd}, // channel:read:redemptions
	}
)

func TypeToLogTitle(t string) string {
	if s, exists := TwitchEventTable[t]; exists {
		return s.LogTitle + LogTextSplit
	} else {
		return fmt.Sprintf("%v%v", t, LogTextSplit)
	}
}

func buildRequestWithModerator(cfg *Config, sessionID, subscType, version string) []byte {
	c := RequestConditionWithModerator{
		BroadcasterUserId: cfg.TargetUserId,
		ModeratorUserId:   cfg.TargetUserId,
	}
	t := SubscriptionTransport{
		Method:    "websocket",
		SessionId: sessionID,
	}
	body := CreateSubscriptionBodyWithModerator{
		Type:      subscType,
		Version:   version,
		Condition: c,
		Transport: t,
	}
	bin, _ := json.Marshal(&body)
	return bin
}

func buildRequest(cfg *Config, sessionID, subscType, version string) []byte {
	c := RequestCondition{
		BroadcasterUserId: cfg.TargetUserId,
	}
	t := SubscriptionTransport{
		Method:    "websocket",
		SessionId: sessionID,
	}
	body := CreateSubscriptionBody{
		Type:      subscType,
		Version:   version,
		Condition: c,
		Transport: t,
	}
	bin, _ := json.Marshal(&body)
	return bin
}

func buildRequestWithUser(cfg *Config, sessionID, subscType, version string) []byte {
	c := RequestConditionWithUser{
		BroadcasterUserId: cfg.TargetUserId,
		UserId:            cfg.TargetUserId,
	}
	t := SubscriptionTransport{
		Method:    "websocket",
		SessionId: sessionID,
	}
	body := CreateSubscriptionBodyWithUser{
		Type:      subscType,
		Version:   version,
		Condition: c,
		Transport: t,
	}
	bin, _ := json.Marshal(&body)
	return bin
}

func buildRequestWithFromUser(cfg *Config, sessionID, subscType, version string) []byte {
	c := RequestConditionWithFromUser{
		FromBroadcasterUserId: cfg.TargetUserId,
	}
	t := SubscriptionTransport{
		Method:    "websocket",
		SessionId: sessionID,
	}
	body := CreateSubscriptionBodyWithFromUser{
		Type:      subscType,
		Version:   version,
		Condition: c,
		Transport: t,
	}
	bin, _ := json.Marshal(&body)
	return bin
}

func handleNotificationDefault(_ *Config, r *Responce, _ []byte, _ *TwitchStats) {
	statsLogger.Info("event(no handler)",
		slog.Any(LogFieldName_Type, r.Payload.Subscription.Type),
	)
}

func handleNotificationChannelSubscribe(_ *BackendContext, _ *Config, r *Responce, raw []byte, s *TwitchStats) {
	v := &ResponceChannelSubscribe{}
	err := json.Unmarshal(raw, &v)
	if err != nil {
		logger.Error("handleNotificationChannelSubscribe::Unmarshal", slog.Any("ERR", err.Error()), slog.Any("raw", string(raw)))
	}
	e := &v.Payload.Event
	if v.Payload.Event.IsGift {
		// ギフトを受け取った人の分は無理に出さなくてよい
		//infoLogger.Info("event(Subscribed<Gift>)",
		//	slog.Any(LogFieldName_Type, r.Payload.Subscription.Type),
		//	slog.Any(LogFieldName_UserName, e.UserName),
		//	slog.Any("tear", e.Tier),
		//	slog.Any("gift", e.IsGift),
		//)
	} else {
		statsLogger.Info("event(Subscribed)",
			slog.Any(LogFieldName_Type, r.Payload.Subscription.Type),
			slog.Any(LogFieldName_UserName, e.UserName),
			slog.Any("tear", e.Tier),
			slog.Any("gift", e.IsGift),
		)
		s.SubScribe(UserName(e.UserName), e.Tier)
	}
}

func handleNotificationChannelCheer(_ *BackendContext, _ *Config, r *Responce, raw []byte, s *TwitchStats) {
	v := &ResponceChannelCheer{}
	err := json.Unmarshal(raw, &v)
	if err != nil {
		logger.Error("handleNotificationChannelCheer::Unmarshal", slog.Any("ERR", err.Error()), slog.Any("raw", string(raw)))
	}
	e := &v.Payload.Event
	statsLogger.Info("event(Cheer)",
		slog.Any(LogFieldName_Type, r.Payload.Subscription.Type),
		slog.Any(LogFieldName_UserName, e.UserName),
		slog.Any("anonymous", e.IsAnonymous),
		slog.Any("bits", e.Bits),
		slog.Any("msg", e.Message),
	)
	s.Cheer(UserName(e.UserName), e.Bits)
}

func handleNotificationStreamOnline(_ *BackendContext, cfg *Config, r *Responce, raw []byte, s *TwitchStats) {
	path := buildLogPath(cfg)
	_, statsLogger = buildLogger(cfg, path)

	v := &ResponceStreamOnline{}
	err := json.Unmarshal(raw, &v)
	if err != nil {
		logger.Error("handleNotificationStreamOnline::Unmarshal", slog.Any("ERR", err.Error()), slog.Any("raw", string(raw)))
	}
	s.StreamStarted()
	e := &v.Payload.Event
	statsLogger.Info("event(Online)",
		slog.Any(LogFieldName_Type, r.Payload.Subscription.Type),
		slog.Any(LogFieldName_UserName, e.BroadcasterUserName),
		slog.Any("at", e.StartedAt),
	)
	os.Remove(cfg.RaidLogPath)
}

func handleNotificationStreamOffline(_ *BackendContext, cfg *Config, r *Responce, raw []byte, s *TwitchStats) {
	v := &ResponceStreamOffline{}
	err := json.Unmarshal(raw, &v)
	if err != nil {
		logger.Error("handleNotificationStreamOffline::Unmarshal", slog.Any("ERR", err.Error()), slog.Any("raw", string(raw)))
	}
	e := &v.Payload.Event
	statsLogger.Info("event(Offline)",
		slog.Any(LogFieldName_Type, r.Payload.Subscription.Type),
		slog.Any(LogFieldName_UserName, e.BroadcasterUserName),
	)
	s.StreamFinished()
	log, _ := os.OpenFile(cfg.StatsLogPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	defer log.Close()
	log.WriteString(s.String())
}

// サブギフした
func handleNotificationChannelSubscriptionGift(_ *BackendContext, _ *Config, r *Responce, raw []byte, s *TwitchStats) {
	v := &ResponceChannelSubscriptionGift{}
	err := json.Unmarshal(raw, &v)
	if err != nil {
		logger.Error("handleNotificationChannelSubscriptionGift::Unmarshal", slog.Any("ERR", err.Error()), slog.Any("raw", string(raw)))
	}
	e := &v.Payload.Event
	statsLogger.Info("event(Gift)",
		slog.Any(LogFieldName_Type, r.Payload.Subscription.Type),
		slog.Any(LogFieldName_UserName, e.UserName),
		slog.Any("tear", e.Tier),
		slog.Any("num", e.Total),
		slog.Any("cumulative", e.CumulativeTotal),
		slog.Any("anonymous", e.IsAnonymous),
	)

	s.SubGift(UserName(e.UserName), e.Total)
}

// 継続サブスクをチャットでシェアした
func handleNotificationChannelSubscriptionMessage(_ *BackendContext, _ *Config, r *Responce, raw []byte, s *TwitchStats) {
	v := &ResponceChannelSubscriptionMessage{}
	err := json.Unmarshal(raw, &v)
	if err != nil {
		logger.Error("handleNotificationChannelSubscriptionMessage::Unmarshal", slog.Any("ERR", err.Error()), slog.Any("raw", string(raw)))
	}
	e := &v.Payload.Event
	statsLogger.Info("event(ReSubscribed)",
		slog.Any(LogFieldName_Type, r.Payload.Subscription.Type),
		slog.Any(LogFieldName_UserName, e.UserName),
		slog.Any("tear", e.Tier),
		slog.Any("duration", e.DurationMonths),
		slog.Any("streak", e.StreakMonths),
		slog.Any("cumlative", e.CumulativeMonths),
	)
	s.SubScribe(UserName(e.UserName), e.Tier)
}

func handleNotificationChannelPointsCustomRewardRedemptionAdd(_ *BackendContext, _ *Config, r *Responce, raw []byte, s *TwitchStats) {
	v := &ResponceChannelPointsCustomRewardRedemptionAdd{}
	err := json.Unmarshal(raw, &v)
	if err != nil {
		logger.Error("handleNotificationChannelPointsCustomRewardRedemptionAdd::Unmarshal", slog.Any("ERR", err.Error()), slog.Any("raw", string(raw)))
	}
	e := &v.Payload.Event
	statsLogger.Info("event(Channel Points)",
		slog.Any(LogFieldName_Type, r.Payload.Subscription.Type),
		slog.Any(LogFieldName_UserName, e.UserName),
		slog.Any("login", e.UserLogin),
		slog.Any("title", e.Reward.Title),
	)
	s.ChannelPoint(UserName(e.UserName), ChannelPointTitle(e.Reward.Title))
}

func handleNotificationChannelChatNotificationSubGifted(_ *BackendContext, _ *Config, r *Responce, e *EventFormatChannelChatNotification, s *TwitchStats) {
	statsLogger.Info("event(SubGiftReceived)",
		slog.Any(LogFieldName_Type, r.Payload.Subscription.Type),
		slog.Any("category", "サブギフ受信"),
		slog.Any("from", e.ChatterUserName),
		slog.Any("to", e.SubGift.RecipientUserName),
	)
	s.SubGifted(UserName(e.SubGift.RecipientUserName), e.SubGift.Sub_Tier)
}

func handleNotificationChannelChatNotificationRaid(ctx *BackendContext, cfg *Config, r *Responce, e *EventFormatChannelChatNotification, s *TwitchStats) {
	statsLogger.Info("event(Raid)",
		slog.Any(LogFieldName_Type, r.Payload.Subscription.Type),
		slog.Any("category", "レイド"),
		slog.Any("from", e.RaId.UserName),
		slog.Any("viewers", e.RaId.ViewerCount),
	)
	s.Raid(UserName(e.RaId.UserName), e.RaId.ViewerCount)
	clipText, clips := ReferUserClips(cfg, e.RaId.UserId)
	log, _ := os.OpenFile(cfg.RaidLogPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	defer log.Close()
	fmt.Fprintf(log, "-- %v さんのクリップ -- \n", e.RaId.UserName)
	log.WriteString(clipText)
	p := &RaidCallbackParam{From: UserName(e.RaId.UserName), Clips: []UserClip{}}
	for _, c := range clips.Data {
		p.Clips = append(p.Clips, UserClip{
			Id:        c.Id,
			Url:       c.Url,
			Thumbnail: c.ThumbnailUrl,
			Title:     c.Title,
			ViewCount: c.ViewCount,
			Duration:  c.Duration,
			Mp4:       ConvertThumbnailToMp4Url(c.ThumbnailUrl),
		})
	}
	if ctx.CallBack.OnRaid != nil {
		ctx.CallBack.OnRaid(p)
	}
}

func handleNotificationChannelChatNotification(ctx *BackendContext, cfg *Config, r *Responce, raw []byte, s *TwitchStats) {
	v := &ResponceChannelChatNotification{}
	err := json.Unmarshal(raw, &v)
	if err != nil {
		logger.Error("handleNotificationChannelChatNotification::Unmarshal", slog.Any("ERR", err.Error()), slog.Any("raw", string(raw)))
	}
	e := &v.Payload.Event
	switch e.NoticeType {
	case "sub":
	case "resub":
		// サブスク継続をチャットで宣言したイベント
		// channel.subscription.message も来るはずなのでそっちでハンドリングする
	case "sub_gift":
		handleNotificationChannelChatNotificationSubGifted(ctx, cfg, r, e, s)
	case "community_sub_gift":
	case "gift_paid_upgrade":
	case "prime_paid_upgrade":
	case "raid":
		handleNotificationChannelChatNotificationRaid(ctx, cfg, r, e, s)
	case "unraid":
	case "pay_it_forward":
	case "announcement":
	case "bits_badge_tier":
	case "charity_donation":
	default:
		logger.Error("event(NotParsed)", slog.Any("raw", string(raw)))
	}
}

func handleNotificationChannelChatMessage(_ *BackendContext, _ *Config, r *Responce, raw []byte, s *TwitchStats) {
	v := &ResponceChatMessage{}
	err := json.Unmarshal(raw, &v)
	if err != nil {
		logger.Error("handleNotificationChannelChatMessage::Unmarshal", slog.Any("ERR", err.Error()), slog.Any("raw", string(raw)))
	}
	e := &v.Payload.Event
	logType := r.Payload.Subscription.Type
	switch e.MessageType {
	case "power_ups_gigantified_emote":
		logType = "巨大化スタンプ"
		s.GigantifiedEmote(UserName(e.ChatterUserName))
	case "power_ups_message_effect":
		logType = "メッセージエフェクト"
		s.MessageEffect(UserName(e.ChatterUserName))
	}
	statsLogger.Info("event(ChatMsg)",
		slog.Any(LogFieldName_Type, logType),
		slog.Any(LogFieldName_UserName, e.ChatterUserName),
		slog.Any(LogFieldName_LoginName, e.ChatterUserLogin),
		slog.Any("text", e.Message.Text),
	)
}

func handleNotificationChannelFollow(_ *BackendContext, _ *Config, r *Responce, raw []byte, s *TwitchStats) {
	v := &ResponceChannelFollow{}
	err := json.Unmarshal(raw, &v)
	if err != nil {
		logger.Error("handleNotificationChannelFollow::Unmarshal", slog.Any("ERR", err.Error()), slog.Any("raw", string(raw)))
	}
	e := &v.Payload.Event
	statsLogger.Info("event(Channel Follow)",
		slog.Any(LogFieldName_Type, r.Payload.Subscription.Type),
		slog.Any(LogFieldName_UserName, e.UserName),
		slog.Any(LogFieldName_LoginName, e.UserLogin),
	)
	s.Follow(UserName(e.UserName))
}

func handleNotificationRaidStarted(_ *BackendContext, cfg *Config, r *Responce, raw []byte, _ *TwitchStats) {
	statsLogger.Info("event(Raid Started)",
		slog.Any(LogFieldName_Type, r.Payload.Subscription.Type),
	)
	if cfg.StopStreamAfterRaided() {
		go func() {
			logger.Info("StopStream Start")
			ticker := time.NewTicker(time.Second * time.Duration(cfg.DelayFromRaidToStop()))
			<-ticker.C
			StopObsStream(cfg)
			logger.Info("StopStream End")
		}()
	}
}
