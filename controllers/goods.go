package controllers

import beego "github.com/beego/beego/v2/server/web"

type GoodsController struct {
	beego.Controller
}

func GetUser(this *beego.Controller) string {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Data["userName"] = ""
	} else {
		this.Data["userName"] = userName.(string)
		return userName.(string)
	}
	return ""
}

func (this *GoodsController) ShowIndex() {
	GetUser(&this.Controller)

	this.TplName = "index.html"
}
