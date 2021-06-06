package controllers

import (
	"shFreshPro/models"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gomodule/redigo/redis"
	// alipay "github.com/smartwalle/alipay/v3"
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

func (this *OrderController) AddOrder() {
	addrid, _ := this.GetInt("addrid")
	payId, _ := this.GetInt("payId")
	skuid := this.GetString("skuids")
	ids := skuid[1 : len(skuid)-1]
	// 多个id组成的数组 通过空格拼接成字符串
	skuids := strings.Split(ids, " ")

	logs.Error(skuids)
	totalCount, _ := this.GetInt("totalCount")
	transferPrice, _ := this.GetInt("transferPrice")
	realyPrice, _ := this.GetInt("realyPrice")

	resp := make(map[string]interface{})
	defer this.ServeJSON()
	//校验数据
	if len(skuids) == 0 {
		resp["code"] = 1
		resp["errmsg"] = "数据库链接错误"
		this.Data["json"] = resp
		return
	}

	o := orm.NewOrm()
	tx, err := o.Begin() // 标识事务的开始
	if err != nil {
		logs.Info("事务开始出错")
		tx.Rollback()
	}

	userName := this.GetSession("userName")
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")

	var order models.OrderInfo
	order.OrderId = time.Now().Format("2006010215030405") + strconv.Itoa(user.Id)
	order.User = &user
	order.Orderstatus = 1
	order.PayMethod = payId
	order.TotalCount = totalCount
	order.TotalPrice = realyPrice
	order.TransitPrice = transferPrice
	//查询地址
	var addr models.Address
	addr.Id = addrid
	o.Read(&addr)

	order.Address = &addr

	//执行插入操作
	o.Insert(&order)

	//想订单商品表中插入数据
	conn, _ := redis.Dial("tcp", "127.0.0.1:6379")

	for _, skuid := range skuids {
		id, _ := strconv.Atoi(skuid)

		var goods models.GoodsSKU
		goods.Id = id

		i := 3
		for i > 0 {
			o.Read(&goods)
			var orderGoods models.OrderGoods
			orderGoods.GoodsSKU = &goods
			orderGoods.OrderInfo = &order

			count, _ := redis.Int(conn.Do("hget", "cart_"+strconv.Itoa(user.Id), id))

			if count > goods.Stock {
				resp["code"] = 2
				resp["errmsg"] = "商品库存不足"
				this.Data["json"] = resp
				tx.Rollback() //标识事务的回滚
				return
			}

			preCount := goods.Stock
			time.Sleep(time.Second * 5)
			logs.Info(preCount, user.Id)

			orderGoods.Count = count
			orderGoods.Price = count * goods.Price

			o.Insert(&orderGoods)
			goods.Stock -= count
			goods.Sales += count

			updateCount, _ := o.QueryTable("GoodsSKU").Filter("Id", goods.Id).Filter("Stock", preCount).Update(orm.Params{"Stock": goods.Stock, "Sales": goods.Sales})
			if updateCount == 0 {
				if i > 0 {
					i -= 1
					continue
				}
				resp["code"] = 3
				resp["errmsg"] = "商品库存改变,订单提交失败"
				this.Data["json"] = resp
				tx.Rollback() //标识事务的回滚
				return
			} else {
				conn.Do("hdel", "cart_"+strconv.Itoa(user.Id), goods.Id)
				break
			}
		}
	}

	//返回数据
	tx.Commit() //提交事务
	resp["code"] = 5
	resp["errmsg"] = "ok"
	this.Data["json"] = resp

}
func (this *OrderController) HandlePay() {
	// var appId = ""
	// var privateKey = ""
	// var client, err = alipay.New(appId, privateKey, true)
	// // 公钥证书见 alipay 说明
	// if err!=nil{
	// 	logs.Info("支付初始化错误！")
	// }
	//获取数据
	orderId := this.GetString("orderId")
	totalPrice := this.GetString("totalPrice")
	logs.Info(orderId, totalPrice)

	// var p = alipay.AliPayTradePagePay{}
	// p.NotifyURL = "http://xxx"
	// p.ReturnURL = "http://127.0.0.1:8080/user/payok"
	// p.Subject = "天天生鲜购物平台"
	// p.OutTradeNo = orderId
	// p.TotalAmount = totalPrice
	// p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	// var url, errs = client.TradePagePay(p)
	// if errs != nil {
	// 	logs.Error(err)
	// }

	// var payURL = url.String()
	// this.Redirect(payURL,302)
	this.Redirect("/", 302)
}

func (this *OrderController) PayOk() {
	orderId := this.GetString("out_trade_no")
	logs.Info(orderId, "orderId")

	//操作数据

	o := orm.NewOrm()
	count, _ := o.QueryTable("OrderInfo").Filter("OrderId", orderId).Update(orm.Params{"Orderstatus": 2})
	if count == 0 {
		logs.Info("更新数据失败")
		this.Redirect("/user/userCenterOrder", 302)
		return
	}
	this.Redirect("/user/userCenterOrder", 302)
}
