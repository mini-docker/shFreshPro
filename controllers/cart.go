package controllers

import (
	"shFreshPro/models"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gomodule/redigo/redis"
)

type CartController struct {
	beego.Controller
}

func (this *CartController) HandleAddCart() {
	//获取数据
	skuid, err1 := this.GetInt("skuid")
	count, err2 := this.GetInt("count")
	resp := make(map[string]interface{})
	defer this.ServeJSON()

	logs.Info("xxxxxxx", skuid, count)

	//校验数据
	if err1 != nil || err2 != nil {
		resp["code"] = 1
		resp["msg"] = "传递的数据不正确"
		this.Data["json"] = resp
		return
	}
	userName := this.GetSession("userName")
	if userName == nil {
		resp["code"] = 2
		resp["msg"] = "当前用户未登录"
		this.Data["json"] = resp
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")

	//处理数据
	//购物车数据存在redis中，用hash
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		logs.Info("redis数据库链接错误")
		return
	}
	// HGET cart_5 7
	// HGET cart_5 8
	//先获取原来的数量，然后给数量加起来
	preCount, err := redis.Int(conn.Do("hget", "cart_"+strconv.Itoa(user.Id), skuid))
	conn.Do("hset", "cart_"+strconv.Itoa(user.Id), skuid, count+preCount)

	rep, err := conn.Do("hlen", "cart_"+strconv.Itoa(user.Id))
	//回复助手函数 购物车商品种类数
	cartCount, _ := redis.Int(rep, err)

	resp["code"] = 5
	resp["msg"] = "Ok"
	resp["cartCount"] = cartCount

	conn.Do("hset", "cart_"+strconv.Itoa(user.Id), "goodsTypeCount", cartCount)

	this.Data["json"] = resp

	// 由于 ServeJSON 直接返回json对象数据
	// {
	// 	cartCount: 2,
	// 	code: 5,
	// 	msg: "Ok",
	// }

}

//展示购物车页面
func (this *CartController) ShowCart() {
	//用户信息
	userName := GetUser(&this.Controller)

	//从redis中获取数据
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		logs.Info("redis链接失败")
		return
	}
	defer conn.Close()

	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	o.Read(&user, "Name")
	goodsMap, _ := redis.IntMap(conn.Do("hgetall", "cart_"+strconv.Itoa(user.Id))) //map[string]int
	logs.Info("goodsMap: ", goodsMap)                                              //["7","8","16","1"]
	goods := make([]map[string]interface{}, len(goodsMap))                         //[{[key]:{}}]
	i := 0
	totalPrice := 0
	totalCount := 0
	for index, value := range goodsMap {
		skuid, _ := strconv.Atoi(index)
		var goodsSku models.GoodsSKU
		goodsSku.Id = skuid
		o.Read(&goodsSku)

		temp := make(map[string]interface{})
		temp["goodsSku"] = goodsSku
		temp["count"] = value

		totalPrice += goodsSku.Price * value
		totalCount += value

		temp["addPrice"] = goodsSku.Price * value

		goods[i] = temp
		i += 1
	}

	this.Data["totalPrice"] = totalPrice
	this.Data["totalCount"] = totalCount

	this.Data["goods"] = goods

	this.TplName = "cart.html"
}

func (this *CartController) HandleUpdateCart() {
	//获取数据
	skuid, err1 := this.GetInt("skuid")
	count, err2 := this.GetInt("count")
	resp := make(map[string]interface{})
	defer this.ServeJSON()

	//校验数据
	if err1 != nil || err2 != nil {
		resp["code"] = 1
		resp["errmsg"] = "请求数据不正确"
		this.Data["json"] = resp
		return
	}
	userName := this.GetSession("userName")
	if userName == nil {
		resp["code"] = 3
		resp["errmsg"] = "当前用户未登录"
		this.Data["json"] = resp
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		resp["code"] = 2
		resp["errmsg"] = "redis数据库链接失败"
		this.Data["json"] = resp
		return
	}
	defer conn.Close()

	conn.Do("hset", "cart_"+strconv.Itoa(user.Id), skuid, count)

	resp["code"] = 5
	resp["errmsg"] = "OK"

	this.Data["json"] = resp
}

func (this *CartController) DeleteCart() {
	skuid, err := this.GetInt("skuid")

	resp := make(map[string]interface{})
	defer this.ServeJSON()

	//校验数据
	if err != nil {
		resp["code"] = 1
		resp["errmsg"] = "请求数据不正确"

		this.Data["json"] = resp
		return
	}

	//链接数据库
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	defer conn.Close()
	if err != nil {
		resp["code"] = 2
		resp["errmsg"] = "redis数据库链接失败"

		this.Data["json"] = resp
		return
	}

	o := orm.NewOrm()
	userName := this.GetSession("userName")
	var user models.User
	user.Name = userName.(string)

	o.Read(&user, "Name")

	conn.Do("hdel", "cart_"+strconv.Itoa(user.Id), skuid)

	// 处理通用商品类型数量
	goodsTypeCount, err := conn.Do("hget", "cart_"+strconv.Itoa(user.Id), "goodsTypeCount")
	logs.Info("goodsTypeCount", goodsTypeCount)
	if err != nil {
		this.Data["goodsTypeCount"] = 0
		logs.Error("redis数据库链接失败")
	} else if goodsTypeCount == nil {
		this.Data["goodsTypeCount"] = 0
	} else {
		this.Data["goodsTypeCount"] = B2S(goodsTypeCount.([]uint8))
		a, _ := strconv.Atoi((this.Data["goodsTypeCount"]).(string))
		logs.Info(a-1, "aaaaa")
		conn.Do("hset", "cart_"+strconv.Itoa(user.Id), "goodsTypeCount", a-1)
	}

	//返回数据
	resp["code"] = 5
	resp["errmsg"] = "ok"
	this.Data["json"] = resp

}
