package backend

import (
	"os"
	"testing"
)

func TestConfig_Load(t *testing.T) {
	raw := []byte(`SUBSCRIBE_USER: user1
CHART_TARGETS:
  - user2
  - user3
AUTH_CODE: auth
CLIENT_ID: client
`)
	if c, e := loadConfigFrom(raw); e != nil {
		t.Errorf("load error [%v]", e.Error())
	} else {
		if c.TargetUser != "" {
			t.Errorf("load error user1 != %v", c.TargetUser)
		}
	}
}

func TestConfig_Save(t *testing.T) {
	dest := "test.yaml"
	raw := []byte(`SERVER_PORT: 1234
`)
	c, e := loadConfigFrom(raw)
	if e != nil {
		t.Errorf("load error [%v]", e.Error())
	}
	if c.LocalPortNum() != 1234 {
		t.Errorf("invalid overlay port 1234 != %v", c.LocalPortNum())
	}
	if e := c.SaveTo(dest); e != nil {
		t.Errorf("save error [%v]", e.Error())
	}
	if _, e := os.Stat(dest); e != nil {
		t.Errorf("file cant write [%v]", e.Error())
	}
	b, err := os.ReadFile(dest)
	if err != nil {
		t.Errorf("file load error [%v]", e.Error())
	}

	c2, err := loadConfigFrom(b)
	if e != nil {
		t.Errorf("load2 error [%v]", e.Error())
	}
	if c2.LocalPortNum() != 1234 {
		t.Errorf("invalid overlay port(2) 1234 != %v", c.LocalPortNum())
	}
	os.Remove(dest)
}

func TestConfig_DontUpdateSavedData(t *testing.T) {
	dest := "test.yaml"
	dest2 := "test2.yaml"
	raw := []byte(`
  DEBUG: true
  CLIP_PLAYER_WIDTH: 300
`)
	os.WriteFile(dest, raw, 0666)
	cfg, err := LoadConfigFromFile(dest)
	if err != nil {
		t.Errorf("file load error %v", err.Error())
	}

	if !cfg.IsDebug() {
		t.Errorf("invalid debug mode flag true != %v", cfg.IsDebug())
	}
	if cfg.ClipWidth() != 300 {
		t.Errorf("invalid clip window width 300 != %v", cfg.ClipWidth())
	}
	if cfg.ClipHeight() != 480 {
		t.Errorf("invalid clip window height 480 != %v", cfg.ClipHeight())
	}

	if err := cfg.SaveTo(dest2); err != nil {
		t.Errorf("save error %v", err)
	}

	cfg, err = LoadConfigFromFile(dest2)
	if err != nil {
		t.Errorf("file load error %v", err.Error())
	}
	if cfg.ClipHeight() != 480 {
		t.Errorf("invalid clip window height 480 != %v", cfg.ClipHeight())
	}
	os.Remove(dest)
	os.Remove(dest2)
}
