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
	Id        string
	Url       string
	Title     string
	Thumbnail string
	ViewCount int
	Duration  float32
	Mp4       string
}

type RaidCallbackParam struct {
	From  UserName
	Clips []UserClip
}
type KeepAliveCallback func()
type ConnectedCallback func()
type RaidCallback func(*RaidCallbackParam)
type CallBack struct {
	KeepAlive   KeepAliveCallback
	OnRaid      RaidCallback
	OnConnected ConnectedCallback
}

type ExitStatus int

const (
	StreamFinished ExitStatus = iota
	ConnectionCanceled
	ConnectionError
)

type BackendContext struct {
	CallBack  *CallBack
	Config    *Config
	Overlay   *OverlayContext
	Stats     *TwitchStats
	Session   *SessionLifecycle
	Refreshed chan struct{}
}

var (
	logger      *slog.Logger
	statsLogger *slog.Logger
)

func buildQuery() string {
	return fmt.Sprintf("keepalive_timeout_seconds=%v", KeepAliveSecond)
}

func connect(localTest bool) (*websocket.Conn, error) {
	var u url.URL
	if localTest {
		u = url.URL{Scheme: LocalTestScheme, Host: LocalTestAddr, Path: ConnectPath, RawQuery: buildQuery()}
	} else {
		u = url.URL{Scheme: GlobalScheme, Host: GlobalHost, Path: ConnectPath, RawQuery: buildQuery()}
	}
	logger.Info("connect", slog.Any("to", u.String()))

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Error("connect::Dial", slog.Any("ERR", err.Error()))
		return nil, err
	}
	return c, nil
}

func receive(cfg *Config, conn *websocket.Conn) (*Responce, []byte, error) {
	r := &Responce{}
	_, message, err := conn.ReadMessage()
	if err != nil {
		logger.Error("receive::ReadMessage", slog.Any("ERR", err.Error()))
		return nil, nil, err
	}
	if cfg.IsDebug() {
		logger.Info("receive", "raw", string(message))
	}
	err = json.Unmarshal(message, &r)
	if err != nil {
		logger.Error("receive::json.Unmarshal", slog.Any("ERR", err.Error()))
		return nil, nil, err
	}
	return r, message, nil
}

// https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#subscription-types
func handleSessionWelcome(ctx *BackendContext, cfg *Config, r *Responce, _ []byte, _ *TwitchStats) {
	if cfg.IsLocalTest() {
		//return
	}
	for k, v := range TwitchEventTable {
		err := createEventSubscription(cfg, r.Payload.Session.Id, k, &v)
		if err != nil {
			logger.Error("handleSessionWelcome::createEventSubscription", slog.Any("ERR", err.Error()))
		}
	}
}

func handleNotification(ctx *BackendContext, cfg *Config, r *Responce, raw []byte, stats *TwitchStats) bool {
	logger.Info("ReceiveNotification", slog.Any("type", r.Payload.Subscription.Type))
	if e, exists := TwitchEventTable[r.Payload.Subscription.Type]; exists {
		e.Handler(ctx, cfg, r, raw, stats)
	} else {
		logger.Error("UNKNOWN notification", slog.Any("Type", r.Payload.Subscription.Type))
	}
	if r.Payload.Subscription.Type == "stream.offline" {
		return false
	}
	return true
}

func progress(ctx *BackendContext, finishChan *chan ExitStatus, cfg *Config, conn *websocket.Conn, stats *TwitchStats) {
	for {
		r, raw, err := receive(cfg, conn)
		if err != nil {
			logger.Error("receive::progress", slog.Any("ERR", err.Error()))
			*finishChan <- ConnectionCanceled
			return
		}
		logger.Info("recv", slog.Any("Type", r.Metadata.MessageType))
		switch r.Metadata.MessageType {
		case "session_welcome":
			logger.Info("progress", slog.Any("event", "connected"))
			handleSessionWelcome(ctx, cfg, r, raw, stats)
			if ctx.CallBack.OnConnected != nil {
				ctx.CallBack.OnConnected()
			}
		case "session_keepalive":
			//logger.Info("progress", slog.Any("event", "keepalive"))
			if ctx.CallBack.KeepAlive != nil {
				ctx.CallBack.KeepAlive()
			}
		case "session_reconnect":
			logger.Info("progress", slog.Any("event", "reconnect"))
		case "notification":
			logger.Info("event: notification")
			if !handleNotification(ctx, cfg, r, raw, stats) {
				*finishChan <- StreamFinished
				return
			}
		case "revocation":
			logger.Info("progress", slog.Any("event", "revocation"))
		default:
			logger.Error("progress::UNKNOWN", slog.Any("Type", r.Metadata.MessageType))
		}
	}
}

func buildLogPath(cfg *Config) string {
	if _, e := os.Stat(cfg.LogPath()); e != nil {
		os.MkdirAll(cfg.LogPath(), 0750)
	}
	if cfg.IsLocalTest() {
		return filepath.Join(cfg.LogPath(), "local.test.txt")
	}
	n := time.Now()
	return filepath.Join(cfg.LogPath(), fmt.Sprintf("%v.txt", n.Format("20060102")))
}

func buildLogger(c *Config, logPath string) (*slog.Logger, *slog.Logger) {
	log, _ := os.OpenFile(logPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	runlog, _ := os.OpenFile(
		filepath.Join(c.LogPath(), "debug.txt"),
		os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666,
	)
	mainLogger := slog.NewTextHandler(runlog, nil)
	// should not pass handler instance after instance of os.Stdout
	// that may depend on spec of slogmulti.Fanout() or slog.NewTextHandler()
	// if pass a handler after os.Stdout, the handler passed after stdout cant write any log
	return slog.New(
			slogmulti.Fanout(
				mainLogger,
				slog.NewTextHandler(os.Stdout, nil),
			),
		),
		slog.New(
			slogmulti.Fanout(
				mainLogger,
				NewTwitchInfoLogger(c, log),
			),
		)
}

func NewBackend(callback *CallBack) *BackendContext {
	ctx := &BackendContext{
		CallBack: callback,
	}
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	ctx.Config = cfg
	path := buildLogPath(cfg)
	logger, statsLogger = buildLogger(cfg, path)
	ctx.Stats = NewTwitchStats()
	ctx.Overlay = NewOverlay(cfg)
	ctx.Refreshed = make(chan struct{})
	ctx.Session = NewSessionLifecycle(&ctx.Refreshed, cfg)
	return ctx
}

func (c *BackendContext) GetOverlayPortNumber() int {
	return c.Config.LocalPortNum()
}

func (c *BackendContext) ServeMain(fin *chan ExitStatus) *websocket.Conn {
	conn, _ := connect(c.Config.IsLocalTest())

	go func() {
		progress(c, fin, c.Config, conn, c.Stats)
	}()
	return conn
}

func (c *BackendContext) Serve() {
	var fin chan ExitStatus
	var conn *websocket.Conn
	statsLogger.Info("ToolVersion", slog.Any(LogFieldName_Type, "ToolVersion"), slog.Any("value", ToolVersion))
	expires, err := ConfirmAccessToken(c.Config)
	if err != nil {
		logger.Error("Serve", slog.Any("msg", "ConfirmAccessToken"), slog.Any("ERR", err.Error()))
		return
	}
	statsLogger.Info("Start", slog.Any(LogFieldName_Type, "TargetUser"), slog.Any("name", c.Config.UserName()), slog.Any("id", c.Config.UserId()))
	c.Session.Serve(expires)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	fin = make(chan ExitStatus)
	conn = c.ServeMain(&fin)

	done := make(chan struct{})
	StartWatcher(c.Config, done)
	if c.Config.OverlayEnabled() {
		c.Overlay.Serve(c.Config)
	}

	for {
		select {
		case status := <-fin:
			//return
			if status == StreamFinished {
				logger.Info("stream finished exit serve")
				done <- struct{}{}
				return
			}
			fin = make(chan ExitStatus)
			conn = c.ServeMain(&fin)
		case <-c.Refreshed:
			logger.Info("Session Refreshed")
			conn.Close()
		case <-interrupt:
			logger.Info("interrupt")
			c.Session.Shutdown()

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Error("write close", slog.Any("ERR", err.Error()))
				return
			}
			conn.Close()
			select {
			case <-fin:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (c *BackendContext) LoadConfig() *ConfigBody {
	return c.Config.LoadRaw()
}

func (c *BackendContext) SaveConfig(cfg *ConfigBody) {
	c.Config.UpdateRaw(cfg)
	if e := c.Config.SaveAll(); e != nil {
		logger.Error("SaveConfig", slog.Any("ERR", e.Error()))
	}
}
