package config_test

import (
	"../config"
	"testing"
)

func TestConfigInit(t *testing.T) {
	conf, err := config.InitConfig()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(conf)
}
