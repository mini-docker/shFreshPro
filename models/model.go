package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct { // 用户表
	Id       int
	Name     string `orm:"size(20);unique"` // 用户名
	Password string `orm:"size(20)"`        // 登录密码
	Email    string `orm:"size(50)"`        // 邮箱
	Active   bool   `orm:"default(false)"`  // 是否激活
	Power    int    `orm:"default(0)"`      // 权限设置  0 表示未激活  1表示激活

}

func init() {

	// 使用全局变量
	GetYaml()

	target_url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		Conf.Mysql.User,
		Conf.Mysql.Password,
		Conf.Mysql.Host,
		Conf.Mysql.Port,
		Conf.Mysql.Name,
	)
	orm.RegisterDataBase("default", "mysql", target_url)

	orm.RegisterModel(new(User))

	orm.RunSyncdb("default", false, true)
}
