package backend

const (
	ToolVersion = "1.2.0"

	LocalTestAddr          = "127.0.0.1:8080"
	LocalTestScheme        = "ws"
	ConfigFilePath         = "config.yaml"
	AuthInfoFile           = ".auth.yaml"
	AuthRedirectUri        = "http://localhost"
	LogFieldName_Type      = "type"
	LogFieldName_UserName  = "user"
	LogFieldName_LoginName = "login"
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
