package backend

const (
	ToolVersion = "1.6.0"

	KeepAliveSecond        = "30"
	GlobalScheme           = "wss"
	GlobalHost             = "eventsub.wss.twitch.tv"
	LocalTestAddr          = "127.0.0.1:8080"
	LocalTestScheme        = "ws"
	ConnectPath            = "/ws"
	ConfigFilePath         = "config.yaml"
	AuthInfoFile           = ".auth.yaml"
	AuthRedirectUri        = "http://localhost"
	LogFieldName_Type      = "type"
	LogFieldName_UserName  = "user"
	LogFieldName_LoginName = "login"
	LogTextSplit           = "   "
	StatsLogPath           = "配信履歴.txt"
	RaidLogPath            = "レイド.txt"
	NotifySoundDefault     = "C:\\Windows\\Media\\chimes.wav"
)

var (
	EventSubScope = []string{
		"bits:read",
		"channel:read:subscriptions",
		"channel:read:redemptions",
		"channel:manage:raids",
		"moderator:read:followers",
		"user:read:chat",
	}
)
