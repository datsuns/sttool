package backend

import "testing"

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
	raw := []byte(`SERVER_PORT: 1234
`)
	c, e := loadConfigFrom(raw)
	if e != nil {
		t.Errorf("load error [%v]", e.Error())
	}
	if c.LocalPortNum() != 1234 {
		t.Errorf("invalid overlay port 1234 != %v", c.LocalPortNum())
	}
}
