package backend

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	slogmulti "github.com/samber/slog-multi"
)

type UserClip struct {
	Url       string
	Title     string
	Thumbnail string
	ViewCount int
}

type RaidCallbackParam struct {
	From  UserName
	Clips []UserClip
}
type KeepAliveCallback func()
type RaidCallback func(*RaidCallbackParam)
type CallBack struct {
	KeepAlive KeepAliveCallback
	OnRaid    RaidCallback
}
type BackgroundContext struct {
	CallBack *CallBack
}

var (
	Debug = false
	Test  = false

	logger      *slog.Logger
	infoLogger  *slog.Logger
	statsLogger *slog.Logger
	logSplit    = "   "

	scheme = "wss"
	addr   = "eventsub.wss.twitch.tv"

	path      = "/ws"
	keepalive = "30"

	stats *TwitchStats
)

func buildQuery() string {
	return fmt.Sprintf("keepalive_timeout_seconds=%v", keepalive)
}

func connect() (*websocket.Conn, error) {
	var u url.URL
	if Test {
		u = url.URL{Scheme: LocalTestScheme, Host: LocalTestAddr, Path: path, RawQuery: buildQuery()}
	} else {
		u = url.URL{Scheme: scheme, Host: addr, Path: path, RawQuery: buildQuery()}
	}
	logger.Info("connecting to " + u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Error("dial error" + err.Error())
		return nil, err
	}
	return c, nil
}

func receive(conn *websocket.Conn) (*Responce, []byte, error) {
	r := &Responce{}
	_, message, err := conn.ReadMessage()
	if err != nil {
		logger.Error("ReadMessage " + err.Error())
		return nil, nil, err
	}
	if Debug {
		logger.Info("receive", "raw", string(message))
	}
	err = json.Unmarshal(message, &r)
	if err != nil {
		logger.Error("json.Unmarshal", "ERR", err.Error())
		return nil, nil, err
	}
	return r, message, nil
}

// https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#subscription-types
func handleSessionWelcome(cfg *Config, r *Responce, raw []byte, _ *TwitchStats) {
	if Test {
		return
	}
	for k, v := range TwitchEventTable {
		err := createEventSubscription(cfg, r, k, &v)
		if err != nil {
			logger.Error("Eventsub Request", "ERROR", err.Error())
		}
	}
}

func handleNotification(ctx *BackgroundContext, cfg *Config, r *Responce, raw []byte, stats *TwitchStats) bool {
	logger.Info("ReceiveNotification", "type", r.Payload.Subscription.Type)
	if e, exists := TwitchEventTable[r.Payload.Subscription.Type]; exists {
		e.Handler(ctx, cfg, r, raw, stats)
	} else {
		logger.Error("UNKNOWN notification", "Type", r.Payload.Subscription.Type)
	}
	if r.Payload.Subscription.Type == "stream.offline" {
		return false
	}
	return true
}

func progress(ctx *BackgroundContext, done *chan struct{}, cfg *Config, conn *websocket.Conn, stats *TwitchStats) {
	for {
		r, raw, err := receive(conn)
		if err != nil {
			break
		}
		logger.Info("recv", "Type", r.Metadata.MessageType)
		switch r.Metadata.MessageType {
		case "session_welcome":
			logger.Info("event: connected")
			handleSessionWelcome(cfg, r, raw, stats)
		case "session_keepalive":
			//logger.Info("event: keepalive")
			ctx.CallBack.KeepAlive()
		case "session_reconnect":
			logger.Info("event: reconnect")
		case "notification":
			logger.Info("event: notification")
			if !handleNotification(ctx, cfg, r, raw, stats) {
				return
			}
		case "revocation":
			logger.Info("event: revocation")
		default:
			logger.Error("UNKNOWN Event", "Type", r.Metadata.MessageType)
		}
	}
}

func buildLogPath(cfg *Config) string {
	if Test {
		return filepath.Join(cfg.LogDest, "local.test.txt")
	}
	n := time.Now()
	return filepath.Join(cfg.LogDest, fmt.Sprintf("%v.txt", n.Format("20060102")))
}

func buildLogger(c *Config, logPath string, debug bool) (*slog.Logger, *slog.Logger, *slog.Logger) {
	log, _ := os.OpenFile(logPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	runlog, _ := os.OpenFile(
		filepath.Join(c.LogDest, "debug.txt"),
		os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if debug {
		return slog.New(
				slogmulti.Fanout(
					slog.NewTextHandler(os.Stdout, nil),
					slog.NewTextHandler(runlog, nil),
					NewTwitchInfoLogger(c, os.Stdout),
				),
			),
			slog.New(NewTwitchInfoLogger(c, log)),
			slog.New(NewTwitchInfoLogger(c, os.Stdout))
	} else {
		return slog.New(
				slogmulti.Fanout(
					slog.NewTextHandler(runlog, nil),
				),
			),
			slog.New(NewTwitchInfoLogger(c, log)),
			slog.New(NewTwitchInfoLogger(c, os.Stdout))
	}
}

func Serve(callback *CallBack) {
	ctx := &BackgroundContext{CallBack: callback}
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	Debug = cfg.DebugMode
	Test = cfg.LocalTest

	path := buildLogPath(cfg)
	logger, statsLogger, infoLogger = buildLogger(cfg, path, Debug)
	cfg.TargetUserId = ReferTargetUserId(cfg)
	statsLogger.Info("ToolVersion", slog.Any(LogFieldName_Type, "ToolVersion"), slog.Any("value", ToolVersion))

	stats = NewTwitchStats()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _ := connect()
	defer c.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		progress(ctx, &done, cfg, c, stats)
	}()
	StartWatcher(cfg, done)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			logger.Info("done")
			return
		//case t := <-ticker.C:
		//	err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		//	log.Println("write:", t)
		//	if err != nil {
		//		log.Println("writeERR:", err)
		//		return
		//	}
		case <-interrupt:
			logger.Info("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Error("write close " + err.Error())
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
