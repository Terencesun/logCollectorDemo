package etcd_test

import (
	"fmt"
	collectorConfig "logCollector/src/config"
	collectorEtcd "logCollector/src/etcd"
	"testing"
	"time"
)

func TestEtcd(t *testing.T)  {
	conf := collectorConfig.ConfigStore{
		EtcdConfig: collectorConfig.Etcd{
			Hosts: []string{"127.0.0.1:32379"},
		},
	}
	// 初始化
	etcdInstance, err := collectorEtcd.InitEtcd(conf.EtcdConfig.Hosts)
	if err != nil {
		t.Error(err)
	}

	etcdInstance.SetWatch("test")

	time.Sleep(time.Second * 5)

	//etcdInstance.SetKey("server/serverLogPath", "./system.log")
	//etcdInstance.SetKey("server/logLevel", "debug")
	etcdInstance.SetKey("server/instances", "[{\"topic\": \"test1\", \"logPath\": \"D:\\\\terence\\\\log\\\\test_log.txt\"},{\"topic\": \"test2\", \"logPath\": \"D:\\\\terence\\\\log\\\\test_.txt1\"}]")
	//etcdInstance.SetKey("kafka/hosts", "[\"127.0.0.1:9092\"]")

	val, err := etcdInstance.GetKey("test")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(val)
	}

	time.Sleep(time.Second * 30)
}
