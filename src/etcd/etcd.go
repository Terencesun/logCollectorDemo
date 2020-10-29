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

func (p *Etcd) GetKey(key string) (val string, err error) {
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

func (p *Etcd) SetKey(key string, val string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	kv := etcd.NewKV(p.Client)
	_, err = kv.Put(ctx, key, val)
	cancel()
	return
}

func (p *Etcd) SetWatch(key string) (watchChan *chan string) {
	// 必须初始化，要不为nil channel，无论给不给值进去都会阻塞
	ch := make(chan string)
	go func() {
		var lastVal string = ""
		watchCh := p.Client.Watch(context.TODO(), key)
		for res := range watchCh {
			//ch <- res
			if lastVal != string(res.Events[0].Kv.Value) {
				ch <- string(res.Events[0].Kv.Value)
			} else {
				continue
			}
		}
	}()
	return &ch
}

func (p *Etcd) GetConfig()  {

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
