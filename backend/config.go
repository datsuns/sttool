package backend

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TargetUser                 string   `yaml:"SUBSCRIBE_USER"`
	AuthCode                   string   `yaml:"AUTH_CODE"`
	ClientId                   string   `yaml:"CLIENT_ID"`
	ChatTargets                []string `yaml:"CHART_TARGETS"`
	NotifySoundFile            string   `yaml:"NOTIFY_SOUND"`
	DebugMode                  bool     `yaml:"DEBUG"`
	LocalTest                  bool     `yaml:"LOCAL_TEST"`
	LogDest                    string   `yaml:"LOG_DEST"`
	TargetUserId               string
	StatsLogPath               string
	RaidLogPath                string
	ObsUrl                     string `yaml:"OBS_URL"`
	ObsPass                    string `yaml:"OBS_PASS"`
	DelayMinutesFromRaidToStop int    `yaml:"DELAY_TO_STOP"`
	NewClipWatchIntervalSecond int    `yaml:"NEW_CLIP_INTERVAL"`
}

func loadConfigFrom(raw []byte) (*Config, error) {
	ret := &Config{
		NotifySoundFile:            NotifySoundDefault,
		LogDest:                    ".",
		DelayMinutesFromRaidToStop: 3,
		NewClipWatchIntervalSecond: 128,
	}
	if e := yaml.Unmarshal(raw, ret); e != nil {
		return nil, e
	}
	ret.StatsLogPath = StatsLogPath
	ret.RaidLogPath = RaidLogPath
	return ret, nil
}

func LoadConfig() (*Config, error) {
	var e error
	f, e := os.Open(ConfigFilePath)
	if e != nil {
		return nil, e
	}
	b, e := io.ReadAll(f)
	if e != nil {
		return nil, e
	}
	return loadConfigFrom(b)
}