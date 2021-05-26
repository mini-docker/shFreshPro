package controllers

import beego "github.com/beego/beego/v2/server/web"

type GoodsController struct {
	beego.Controller
}

func (this *GoodsController) ShowIndex() {
	this.TplName = "index.html"
}
