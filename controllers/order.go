package controllers

import (
	"shFreshPro/models"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gomodule/redigo/redis"
)

type OrderController struct {
	beego.Controller
}

func (this *OrderController) ShowOrder() {

	skuids := this.GetStrings("skuid")
	logs.Info(skuids)

	// 校验数据
	if len(skuids) == 0 {
		logs.Info("请求数据错误")
		this.Redirect("/user/cart", 302)
		return
	}

	//处理数据
	o := orm.NewOrm()
	conn, _ := redis.Dial("tcp", "127.0.0.1:6379")
	defer conn.Close()
	//获取用户数据
	var user models.User
	userName := this.GetSession("userName")
	user.Name = userName.(string)
	o.Read(&user, "Name")

	goodsBuffer := make([]map[string]interface{}, len(skuids))

	totalPrice := 0
	totalCount := 0
	for index, skuid := range skuids {
		temp := make(map[string]interface{})

		id, _ := strconv.Atoi(skuid)
		//查询商品数据
		var goodsSku models.GoodsSKU
		goodsSku.Id = id
		o.Read(&goodsSku)

		temp["goods"] = goodsSku
		//获取商品数量
		count, _ := redis.Int(conn.Do("hget", "cart_"+strconv.Itoa(user.Id), id))
		temp["count"] = count
		//计算小计
		amount := goodsSku.Price * count
		temp["amount"] = amount

		//计算总金额和总件数
		totalCount += count
		totalPrice += amount

		goodsBuffer[index] = temp
	}

	this.Data["goodsBuffer"] = goodsBuffer

	//获取地址数据
	var addrs []models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Id", user.Id).All(&addrs)
	this.Data["addrs"] = addrs

	//传递总金额和总件数
	this.Data["totalPrice"] = totalPrice
	this.Data["totalCount"] = totalCount
	transferPrice := 10
	this.Data["transferPrice"] = transferPrice
	this.Data["realyPrice"] = totalPrice + transferPrice

	//传递所有商品的id
	this.Data["skuids"] = skuids

	//返回视图
	this.TplName = "place_order.html"
}
