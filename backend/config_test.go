package backend

import "testing"

func TestConfigLoad(t *testing.T) {
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
		if c.TargetUser != "user1" {
			t.Errorf("load error user1 != %v", c.TargetUser)
		}
	}

}
