package api

import (
	"encoding/json"
	"github.com/astaxie/beego"
	collectorConfig "logCollector/src/config"
	collectorEtcd "logCollector/src/etcd"
)

var EtcdClient *collectorEtcd.Etcd

func InitController(conf *collectorConfig.Etcd) (err error) {
	// 初始化etcd连接
	EtcdClient, err = collectorEtcd.InitEtcd(conf.Hosts)
	if err != nil {
		return
	}
	return
}

type ResTemp struct {
	Code int
	Data interface{}
}

type DeleteModel struct {
	LogPath string `json:"logPath"`
}

type LogApiController struct {
	beego.Controller
}

func (p *LogApiController) List() {
	CODE := make(map[string]int)
	CODE["OK"] = 10
	CODE["ERROR"] = 11
	ret := ResTemp{}
	tmp := make([]collectorConfig.Instance, 0)
	val, err := EtcdClient.GetKey("server/instances")
	if err != nil {
		ret.Code = CODE["ERROR"]
		ret.Data = "error"
		p.Data["json"] = ret
	} else {
		err = json.Unmarshal([]byte(val), &tmp)
		if err != nil {
			ret.Code = CODE["ERROR"]
			ret.Data = "unmarshal error"
			p.Data["json"] = ret
		} else {
			ret.Code = CODE["OK"]
			ret.Data = tmp
			p.Data["json"] = ret
		}
	}
	p.ServeJSON()
}

func (p *LogApiController) Delete() {
	CODE := make(map[string]int)
	CODE["OK"] = 10
	CODE["INSTANCE_NO_FOUND"] = 11
	CODE["ERROR"] = 12
	params := DeleteModel{}
	// 删除操作
	ret := ResTemp{}
	tmp := make([]collectorConfig.Instance, 0)
	tmp2 := make([]collectorConfig.Instance, 0)

	val, err := EtcdClient.GetKey("server/instances")
	if err != nil {
		ret.Code = CODE["ERROR"]
		ret.Data = "error"
		p.Data["json"] = ret
		p.ServeJSON()
		return
	}

	err = json.Unmarshal([]byte(val), &tmp)
	if err != nil {
		ret.Code = CODE["ERROR"]
		ret.Data = "unmarshal error"
		p.Data["json"] = ret
		p.ServeJSON()
		return
	}

	// 查看删除的东西在不在
	err = json.Unmarshal(p.Ctx.Input.RequestBody, &params)
	if err != nil {
		ret.Code = CODE["ERROR"]
		ret.Data = "unmarshal error"
		p.Data["json"] = ret
		p.ServeJSON()
		return
	}
	var l = false
	for _, v := range tmp {
		if v.LogFilePath == params.LogPath {
			l = true
		} else {
			tmp2 = append(tmp2, v)
		}
	}

	if l {
		newVal, err := json.Marshal(tmp2)
		if err != nil {
			ret.Code = CODE["ERROR"]
			ret.Data = "marshal error"
			p.Data["json"] = ret
			p.ServeJSON()
			return
		}
		err = EtcdClient.SetKey("server/instances", string(newVal))
		if err != nil {
			ret.Code = CODE["ERROR"]
			ret.Data = "set key error"
			p.Data["json"] = ret
			p.ServeJSON()
			return
		}
		val, err := EtcdClient.GetKey("server/instances")
		if err != nil {
			ret.Code = CODE["ERROR"]
			ret.Data = "error"
			p.Data["json"] = ret
			p.ServeJSON()
			return
		}
		err = json.Unmarshal([]byte(val), &tmp)
		if err != nil {
			ret.Code = CODE["ERROR"]
			ret.Data = "unmarshal error"
			p.Data["json"] = ret
			p.ServeJSON()
			return
		}

		ret.Code = CODE["OK"]
		ret.Data = tmp
		p.Data["json"] = ret
		p.ServeJSON()
		return
	} else {
		ret.Code = CODE["INSTANCE_NO_FOUND"]
		ret.Data = "path no found"
		p.Data["json"] = ret
		p.ServeJSON()
		return
	}
}

func (p *LogApiController) Create() {
	CODE := make(map[string]int)
	CODE["OK"] = 10
	CODE["PATH_EXIST"] = 11
	CODE["ERROR"] = 12
	params := collectorConfig.Instance{}
	ret := ResTemp{}
	tmp := make([]collectorConfig.Instance, 0)

	val, err := EtcdClient.GetKey("server/instances")
	if err != nil {
		ret.Code = CODE["ERROR"]
		ret.Data = "error"
		p.Data["json"] = ret
		p.ServeJSON()
		return
	}

	err = json.Unmarshal([]byte(val), &tmp)
	if err != nil {
		ret.Code = CODE["ERROR"]
		ret.Data = "unmarshal error"
		p.Data["json"] = ret
		p.ServeJSON()
		return
	}

	err = json.Unmarshal(p.Ctx.Input.RequestBody, &params)
	if err != nil {
		ret.Code = CODE["ERROR"]
		ret.Data = "unmarshal error"
		p.Data["json"] = ret
		p.ServeJSON()
		return
	}
	var l = false
	loop: for _, v := range tmp {
		if v.LogFilePath == params.LogFilePath {
			l = true
			break loop
		}
	}
	if l {
		ret.Code = CODE["PATH_EXIST"]
		ret.Data = "the path is exist"
		p.Data["json"] = ret
		p.ServeJSON()
		return
	}

	// 不存在
	tmp = append(tmp, params)
	newVal, err := json.Marshal(tmp)
	if err != nil {
		ret.Code = CODE["ERROR"]
		ret.Data = "marshal error"
		p.Data["json"] = ret
		p.ServeJSON()
		return
	}
	err = EtcdClient.SetKey("server/instances", string(newVal))
	if err != nil {
		ret.Code = CODE["ERROR"]
		ret.Data = "set key error"
		p.Data["json"] = ret
		p.ServeJSON()
		return
	}
	val, err = EtcdClient.GetKey("server/instances")
	if err != nil {
		ret.Code = CODE["ERROR"]
		ret.Data = "error"
		p.Data["json"] = ret
		p.ServeJSON()
		return
	}
	err = json.Unmarshal([]byte(val), &tmp)
	if err != nil {
		ret.Code = CODE["ERROR"]
		ret.Data = "unmarshal error"
		p.Data["json"] = ret
		p.ServeJSON()
		return
	}

	ret.Code = CODE["OK"]
	ret.Data = tmp
	p.Data["json"] = ret
	p.ServeJSON()
	return
}
