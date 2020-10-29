package etcd

import (
	"context"
	etcd "github.com/coreos/etcd/clientv3"
	collectorConfig "logCollector/src/config"
	"time"
)

type Etcd struct {
	Client *etcd.Client
}

func (p *Etcd) GetConfig(key string) (val string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	kv := etcd.NewKV(p.Client)
	res, err := kv.Get(ctx, key)
	cancel()
	if err != nil {
		return
	}
	if len(res.Kvs) != 0 {
		val = string(res.Kvs[0].Value)
	}
	return
}

func (p *Etcd) SetConfig(key string, val string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	kv := etcd.NewKV(p.Client)
	_, err = kv.Put(ctx, key, val)
	cancel()
	return
}

func (p *Etcd) SetWatch(key string) (watchChan *chan interface{}) {
	// 必须初始化，要不为nil channel，无论给不给值进去都会阻塞
	ch := make(chan interface{})
	go func() {
		watchCh := p.Client.Watch(context.TODO(), key)
		for res := range watchCh {
			ch <- res
		}
	}()
	return &ch
}

func InitEtcd(conf collectorConfig.ConfigStore) (etcdInstance *Etcd, err error) {
	etcdInstance = &Etcd{}
	etcdInstance.Client, err = etcd.New(etcd.Config{
		Endpoints: conf.EtcdConfig.Hosts,
		DialTimeout: 30 * time.Second,
	})
	if err != nil {
		return
	}
	return
}
