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
	etcdInstance, err := collectorEtcd.InitEtcd(conf)
	if err != nil {
		t.Error(err)
	}

	etcdInstance.SetWatch("test")

	time.Sleep(time.Second * 5)

	etcdInstance.SetKey("test", "test3")
	etcdInstance.SetKey("test", "test3")

	val, err := etcdInstance.GetKey("test")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(val)
	}

	time.Sleep(time.Second * 30)
}
