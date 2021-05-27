package routers

import (
	"shFreshPro/controllers"

	beego "github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func init() {
	beego.InsertFilter("/user/*", beego.BeforeExec, filterFunc)
	// beego.Router("/", &controllers.MainController{})
	beego.Router("/", &controllers.GoodsController{}, "get:ShowIndex")
	beego.Router("/register", &controllers.UserController{}, "get:ShowReg;post:HandleReg")
	//激活用户
	beego.Router("/active", &controllers.UserController{}, "get:ActiveUser")
	//用户登录
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	//退出登录
	beego.Router("/user/logout", &controllers.UserController{}, "get:Logout")
	//用户中心信息页
	beego.Router("/user/userCenterInfo", &controllers.UserController{}, "get:ShowUserCenterInfo")
	//用户中心订单页
	beego.Router("/user/userCenterOrder", &controllers.UserController{}, "get:ShowUserCenterOrder")
	//用户中心地址页  命名语义化 命名即注释
	beego.Router("/user/userCenterSite", &controllers.UserController{}, "get:ShowUserCenterSite;post:HandleUserCenterSite")

}

var filterFunc = func(ctx *beecontext.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
