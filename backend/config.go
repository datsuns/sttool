package backend

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigBody struct {
	ChatTargets                []string `yaml:"CHART_TARGETS"`
	NotifySoundFile            string   `yaml:"NOTIFY_SOUND"`
	DebugMode                  bool     `yaml:"DEBUG"`
	LocalTest                  bool     `yaml:"LOCAL_TEST"`
	LogDest                    string   `yaml:"LOG_DEST"`
	ObsUrl                     string   `yaml:"OBS_URL"`
	ObsPass                    string   `yaml:"OBS_PASS"`
	DelayMinutesFromRaidToStop int      `yaml:"DELAY_TO_STOP"`
	NewClipWatchIntervalSecond int      `yaml:"NEW_CLIP_INTERVAL"`
	LocalServerPortNumber      int      `yaml:"SERVER_PORT"`
	OverlayEnabled             bool     `yaml:"OVERLAY_ENABLE"`
	ClipPlayerWidth            int      `yaml:"CLIP_PLAYER_WIDTH"`
	ClipPlayerHeight           int      `yaml:"CLIP_PLAYER_HEIGHT"`
}

type AuthEntry struct {
	AuthCode     string `yaml:"AUTH_CODE"`
	RefreshToken string `yaml:"REFRESH_TOKEN"`
}

type Config struct {
	Body            ConfigBody
	Auth            AuthEntry
	AppClientId     string
	AppClientSecret string
	TargetUser      string
	TargetUserId    string
	StatsLogPath    string
	RaidLogPath     string
}

var (
	DefaultConfig = ConfigBody{
		ChatTargets:                []string{},
		NotifySoundFile:            NotifySoundDefault,
		DebugMode:                  false,
		LocalTest:                  false,
		LogDest:                    ".",
		ObsUrl:                     "",
		ObsPass:                    "",
		DelayMinutesFromRaidToStop: 3,
		NewClipWatchIntervalSecond: 128,
		LocalServerPortNumber:      8930,
		OverlayEnabled:             false,
		ClipPlayerWidth:            640,
		ClipPlayerHeight:           480,
	}
)

func loadConfigFrom(raw []byte) (*Config, error) {
	ret := &Config{}
	ret.Init()
	if e := yaml.Unmarshal(raw, &ret.Body); e != nil {
		return nil, e
	}
	return ret, nil
}

func setDefaultConfig(path string) (*Config, error) {
	var err error
	raw, err := yaml.Marshal(DefaultConfig)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(path, raw, 0644)
	if err != nil {
		return nil, err
	}
	ret := &Config{}
	ret.Init()
	return ret, nil
}

func LoadConfig() (*Config, error) {
	var e error
	f, e := os.Open(ConfigFilePath)
	if e != nil {
		return setDefaultConfig(ConfigFilePath)
	}
	b, e := io.ReadAll(f)
	if e != nil {
		return nil, e
	}
	return loadConfigFrom(b)
}

func (c *Config) Init() {
	c.Body = DefaultConfig
	c.Auth = AuthEntry{
		AuthCode:     "",
		RefreshToken: "",
	}
	c.AppClientId = AppClientID
	c.AppClientSecret = AppClientSecret
	c.StatsLogPath = StatsLogPath
	c.RaidLogPath = RaidLogPath
}

func (c *Config) Save() error {
	var err error
	body, err := yaml.Marshal(c.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(ConfigFilePath, body, 0644)
	if err != nil {
		return err
	}
	auth, err := yaml.Marshal(c.Auth)
	if err != nil {
		return err
	}
	err = os.WriteFile(AuthInfoFile, auth, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) LoadAuthConfig() error {
	var e error
	f, e := os.Open(AuthInfoFile)
	if e != nil {
		return e
	}
	b, e := io.ReadAll(f)
	if e != nil {
		return e
	}
	if e := yaml.Unmarshal(b, &c.Auth); e != nil {
		return e
	}
	return nil
}

func (c *Config) UpdatAccessToken(auth AuthEntry) {
	c.Auth = auth
}

func (c *Config) UserName() string {
	return c.TargetUser
}

func (c *Config) UserId() string {
	return c.TargetUserId
}

func (c *Config) AuthCode() string {
	return c.Auth.AuthCode
}

func (c *Config) RefreshToken() string {
	return c.Auth.RefreshToken
}

func (c *Config) ClientId() string {
	return c.AppClientId
}

func (c *Config) ClientSecret() string {
	return c.AppClientSecret
}

func (c *Config) IsDebug() bool {
	return c.Body.DebugMode
}

func (c *Config) IsLocalTest() bool {
	return c.Body.LocalTest
}

func (c *Config) LocalPortNum() int {
	return c.Body.LocalServerPortNumber
}

func (c *Config) LogPath() string {
	return c.Body.LogDest
}

func (c *Config) DelayFromRaidToStop() int {
	return c.Body.DelayMinutesFromRaidToStop
}

func (c *Config) OverlayEnabled() bool {
	return c.Body.OverlayEnabled
}

func (c *Config) NotifySoundFilePath() string {
	return c.Body.NotifySoundFile
}

func (c *Config) ClipWatchInterval() int {
	return c.Body.NewClipWatchIntervalSecond
}

func (c *Config) ObsUrl() string {
	return c.Body.ObsUrl
}

func (c *Config) ObsPass() string {
	return c.Body.ObsPass
}

func (c *Config) CipWidth() int {
	return c.Body.ClipPlayerWidth
}

func (c *Config) CipHeight() int {
	return c.Body.ClipPlayerHeight
}
