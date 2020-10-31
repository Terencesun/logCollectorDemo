package api

import (
	"github.com/astaxie/beego"
	collectorConfig "logCollector/src/config"
)

func Start(host string, conf *collectorConfig.Etcd) (err error) {
	beego.BConfig.CopyRequestBody = true
	err = InitController(conf)
	if err != nil {
		return
	}
	beego.Run(host)
	return
}
