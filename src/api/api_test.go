package api_test

import (
	"fmt"
	collectorApi "logCollector/src/api"
	collectorConfig "logCollector/src/config"
	"testing"
)

func TestApi(t *testing.T)  {
	conf := &collectorConfig.Etcd{
		Hosts: []string{"127.0.0.1:32379"},
	}
	fmt.Println(conf)
	err := collectorApi.Start("127.0.0.1:1234", conf)
	fmt.Println(err)
}
