package api

import "github.com/astaxie/beego"

func init()  {
	beego.Router("/api/list", &LogApiController{}, "post:List")
	beego.Router("/api/delete", &LogApiController{}, "post:Delete")
	beego.Router("/api/create", &LogApiController{}, "post:Create")
}
