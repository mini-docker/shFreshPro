package controllers

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"shFreshPro/models"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/utils"
	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	beego.Controller
}

// 显示注册页面
func (this *UserController) ShowReg() {
	this.TplName = "register.html"
}

func (this *UserController) HandleReg() {
	// 1.获取数据
	userName := this.GetString("user_name")
	pwd := this.GetString("pwd")
	cpwd := this.GetString("cpwd")
	email := this.GetString("email")
	fmt.Println("register page!!!!!", userName, pwd, cpwd, email)
	// 2.校验数据
	if userName == "" || pwd == "" || cpwd == "" || email == "" {
		this.Data["errmsg"] = "数据不完整，请重新注册"
		this.TplName = "register.html"
		return
	}
	if pwd != cpwd {
		this.Data["errmsg"] = "两次输入密码不一致，请重新注册！"
		this.TplName = "register.html"
		return
	}
	reg, _ := regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	res := reg.FindString(email)
	if res == "" {
		this.Data["errmsg"] = "邮箱格式不正确"
		this.TplName = "register.html"
		return
	}
	// 3.处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	user.Password = pwd
	user.Email = email

	_, err := o.Insert(&user)
	if err != nil {
		this.Data["errmsg"] = "注册失败,请更换数据注册"
		this.TplName = "register.html"
		return
	}
	//发送邮件 短信和邮件信息保存到隐藏文件以外
	emailConfig := models.Conf.User.EmailConfig
	emailConn := utils.NewEMail(emailConfig)
	emailConn.From = models.Conf.User.From
	emailConn.To = []string{email}
	emailConn.Subject = "天天生鲜用户注册"
	//注意这里我们发送给用户的是激活请求地址
	emailConn.Text = "192.168.1.221:8080/active?id=" + strconv.Itoa(user.Id)

	emailConn.Send()
	this.Ctx.WriteString("注册成功，请去相应邮箱激活用户！")
}

// 激活处理
func (this *UserController) ActiveUser() {
	// 获取数据
	id, err := this.GetInt("id")
	// 校验数据
	if err != nil {
		this.Data["errmsg"] = "要激活的用户不存在"
		this.TplName = "register.html"
		return
	}
	//处理数据
	//更新操作
	o := orm.NewOrm()
	var user models.User
	user.Id = id
	err = o.Read(&user)
	if err != nil {
		this.Data["errmsg"] = "要激活的用户不存在"
		this.TplName = "register.html"
		return
	}
	user.Active = true
	o.Update(&user)

	//返回视图
	this.Redirect("/login", 302)
}

//展示登录页面
func (this *UserController) ShowLogin() {
	userName := this.Ctx.GetCookie("userName")
	//解码
	temp, _ := base64.StdEncoding.DecodeString(userName)
	if string(temp) == "" {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	} else {
		this.Data["userName"] = string(temp)
		this.Data["checked"] = "checked"
	}
	this.TplName = "login.html"
}

//处理登录业务
//处理登录业务
func (this *UserController) HandleLogin() {

	//1.获取数据
	userName := this.GetString("username")
	pwd := this.GetString("pwd")

	//2.校验数据
	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "登录数据不完整，请重新输入！"
		this.TplName = "login.html"
		return
	}
	//3.处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName

	err := o.Read(&user, "Name")
	if err != nil {
		this.Data["errmsg"] = "用户名或密码错误，请重新输入！"
		this.TplName = "login.html"
		return
	}
	if user.Password != pwd {
		this.Data["errmsg"] = "用户名或密码错误，请重新输入！"
		this.TplName = "login.html"
		return
	}
	if user.Active != true {
		this.Data["errmsg"] = "用户未激活，请先往邮箱激活！"
		this.TplName = "login.html"
		return
	}

	//4.返回视图1‘
	remember := this.GetString("remember")

	//base64加密
	if remember == "on" {
		temp := base64.StdEncoding.EncodeToString([]byte(userName))
		// fmt.Println("temp: ", temp)
		this.Ctx.SetCookie("userName", temp, 24*3600*30)
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}

	//跳转到首页,
	/*
		1.首页的简单显示实现
		2.登录判断（路由过滤器）
		3.首页显示
		4.三个页面
			视图布局
			添加地址页（如何让页面只显示一个地址）
			用户中心信息页显示
	*/
	this.SetSession("userName", userName)
	// this.Ctx.WriteString("登录成功")
	this.Redirect("/", 302)
}

func (this *UserController) Logout() {
	this.DelSession("userName")
	//跳转视图
	this.Redirect("/login", 302)
}

//用户中心信息页
func (this *UserController) ShowUserCenterInfo() {

	userName := GetUser(&this.Controller)
	this.Data["userName"] = userName

	o := orm.NewOrm()
	var addr models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Name", userName).Filter("Isdefault", true).One(&addr)

	if addr.Id == 0 {
		this.Data["addr"] = ""
	} else {
		this.Data["addr"] = addr
	}

	this.Layout = "userCenterLayout.html"
	this.TplName = "user_center_info.html"
}

//用户中心订单页
func (this *UserController) ShowUserCenterOrder() {

	userName := GetUser(&this.Controller)
	this.Data["userName"] = userName

	this.Layout = "userCenterLayout.html"
	this.TplName = "user_center_order.html"
}

//
func (this *UserController) ShowUserCenterSite() {
	userName := GetUser(&this.Controller)
	this.Data["userName"] = userName

	o := orm.NewOrm()
	var addr models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Name", userName).Filter("Isdefault", true).One(&addr)

	if addr.Id == 0 {
		this.Data["addr"] = ""
		this.Data["address"] = ""
		this.Data["isaddr"] = 0
	} else {
		this.Data["addr"] = addr
		this.Data["address"] = addr
		this.Data["isaddr"] = 1
	}

	logs.Info(this.Data["isaddr"], "llllllll")
	this.Layout = "userCenterLayout.html"
	this.TplName = "user_center_site.html"
}

//
func (this *UserController) HandleUserCenterSite() {
	receiver := this.GetString("receiver")
	addr := this.GetString("addr")
	zipCode := this.GetString("zipCode")
	phone := this.GetString("phone")

	if receiver == "" || addr == "" || zipCode == "" || phone == "" {
		logs.Info("添加数据不完整")
		this.Redirect("/user/userCenterSite", 302)
		return
	}

	o := orm.NewOrm()
	var addrUser models.Address
	addrUser.Isdefault = true
	err := o.Read(&addrUser, "Isdefault")
	// 添加默认地址之前需要把原来的默认地址更新成非默认地址
	if err == nil {
		addrUser.Isdefault = false
		o.Update(&addrUser)
	}
	// 更新默认地址时，给原来的地址对象的ID赋值了，这时候用原来的地址对象插入，意思是用原来的ID做插入操作,会报错
	// 关联
	userName := this.GetSession("userName")
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")

	var addUserNew models.Address
	addUserNew.Receiver = receiver
	addUserNew.Zipcode = zipCode
	addUserNew.Addr = addr
	addUserNew.Phone = phone
	addUserNew.Isdefault = true
	addUserNew.User = &user
	o.Insert(&addUserNew)

	this.Redirect("/user/userCenterSite", 302)
}
