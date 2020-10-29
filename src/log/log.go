package log

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	collectionConfig "logCollector/src/config"
)

func InitLogger(conf collectionConfig.ConfigStore) (err error) {
	config := make(map[string]interface{})
	config["filename"] = conf.SystemLogPath
	switch conf.SystemLogLevel {
	case "debug":
		config["level"] = logs.LevelDebug
	case "warn":
		config["level"] = logs.LevelWarn
	case "alert":
		config["level"] = logs.LevelAlert
	case "info":
		config["level"] = logs.LevelInfo
	default:
		config["level"] = logs.LevelInfo
	}

	configByte, _ := json.Marshal(config)

	err = logs.SetLogger(logs.AdapterFile, string(configByte))

	if err != nil {
		return
	}

	return
}
