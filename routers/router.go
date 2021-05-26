package routers

import (
	"shFreshPro/controllers"

	beego "github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func init() {
	beego.InsertFilter("/user/*", beego.BeforeExec, filterFunc)
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.UserController{}, "get:ShowReg;post:HandleReg")
	//激活用户
	beego.Router("/active", &controllers.UserController{}, "get:ActiveUser")
	//用户登录
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
}

var filterFunc = func(ctx *beecontext.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
