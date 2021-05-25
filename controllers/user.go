package controllers

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"shFreshPro/models"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
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
		this.Ctx.SetCookie("userName", temp, 3*60)
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}

	this.Ctx.WriteString("登录成功")
}