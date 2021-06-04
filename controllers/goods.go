package controllers

import (
	"encoding/json"
	"math"
	"shFreshPro/models"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gomodule/redigo/redis"
)

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

func ShowLayout(this *beego.Controller) {
	o := orm.NewOrm()
	var types []models.GoodsType
	o.QueryTable("GoodsType").All(&types)
	this.Data["types"] = types
	GetUser(this)
	this.Layout = "goodsLayout.html"
}

func PageTool(pageCount int, pageIndex int) []int {

	var pages []int
	if pageCount <= 5 {
		pages = make([]int, pageCount)
		for i, _ := range pages {
			pages[i] = i + 1
		}

		//pages = [1,2,..,pageCount]
	} else if pageIndex <= 3 {
		//pages := make([]int,5)
		pages = []int{1, 2, 3, 4, 5}
	} else if pageIndex > pageCount-3 {
		//pages = [6, 7, 8, 9, 10]
		pages = []int{pageCount - 4, pageCount - 3, pageCount - 2, pageCount - 1, pageCount}
	} else {
		pages = []int{pageIndex - 2, pageIndex - 1, pageIndex, pageIndex + 1, pageIndex + 2}
	}
	return pages

}

func (this *GoodsController) ShowIndex() {
	GetUser(&this.Controller)

	o := orm.NewOrm()
	//获取类型数据
	var goodsTypes []models.GoodsType
	o.QueryTable("GoodsType").All(&goodsTypes)
	this.Data["goodsTypes"] = goodsTypes

	//获取轮播图数据
	var indexGoodsBanner []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&indexGoodsBanner)
	this.Data["indexGoodsBanner"] = indexGoodsBanner

	//获取促销商品数据
	var promotionGoods []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&promotionGoods)
	this.Data["promotionsGoods"] = promotionGoods

	//首页展示商品数据
	goods := make([]map[string]interface{}, len(goodsTypes))

	//向切片interface中插入类型数据
	for index, value := range goodsTypes {
		//获取对应类型的首页展示商品
		temp := make(map[string]interface{})
		temp["type"] = value
		goods[index] = temp
	}

	//商品数据
	for _, value := range goods {
		var textGoods []models.IndexTypeGoodsBanner
		var imgGoods []models.IndexTypeGoodsBanner
		//获取文字商品数据
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").OrderBy("Index").Filter("GoodsType", value["type"]).Filter("DisplayType", 0).All(&textGoods)
		//获取图片商品数据
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").OrderBy("Index").Filter("GoodsType", value["type"]).Filter("DisplayType", 1).All(&imgGoods)

		value["textGoods"] = textGoods
		value["imgGoods"] = imgGoods
	}
	this.Data["goods"] = goods
	a, _ := json.Marshal(goods)

	var b []map[string]interface{}
	err := json.Unmarshal(a, &b)
	if err != nil {
		logs.Info("UnMarshal error", err)
	}

	// logs.Info("goods: ", b)
	this.TplName = "index.html"
}

// goods
// [imgGoods:
//     [map
//         [DisplayType:1
//             GoodsSKU:map[Desc:草莓简介
//                             Goods:map[Detail: GoodsSKU:<nil> Id:1 Name:]
//                             GoodsImage:<nil>
//                             GoodsType:map[GoodsSKU:<nil> Id:1 Image: IndexTypeGoodsBanner:<nil> Logo: Name:]
//                             Id:1
//                             Image:group1/M00/00/00/rBCzg1oKqFGAR2tjAAAljHPuXJg4272079
//                             IndexGoodsBanner:<nil>
//                             IndexTypeGoodsBanner:<nil>
//                             Name:草莓 500g
//                             OrderGoods:<nil>
//                             Price:10
//                             Sales:0
//                             Status:1
//                             Stock:98
//                             Time:2017-11-15T04:10:14+04:00
//                             Unite:500g
//                         ]
//             GoodsType:map[GoodsSKU:<nil>
//                             Id:1
//                             Image:group1/M00/00/00/rBCzg1oKeNKAEl87AAAmv27pX4k4942898
//                             IndexTypeGoodsBanner:<nil>
//                             Logo:fruit
//                             Name:新鲜水果
//                         ]
//             Id:1
//             Index:0
//         ]
//     ]
// ]

//展示商品详情
func (this *GoodsController) ShowGoodsDetail() {
	id, err := this.GetInt("id")
	if err != nil {
		logs.Error("浏览器请求错误")
		this.Redirect("/", 302)
	}
	o := orm.NewOrm()
	var goodsSku models.GoodsSKU
	goodsSku.Id = id
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType", "Goods").Filter("Id", id).One(&goodsSku)

	var goodsNew []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType", goodsSku.GoodsType).OrderBy("Time").Limit(2, 0).All(&goodsNew)
	this.Data["goodsNew"] = goodsNew

	this.Data["goodsSku"] = goodsSku

	//添加历史浏览记录
	//判断用户是否登录
	userName := this.GetSession("userName")
	if userName != nil {
		//查询用户信息
		o := orm.NewOrm()
		var user models.User
		user.Name = userName.(string)
		o.Read(&user, "Name")
		//添加历史记录,用redis存储
		conn, err := redis.Dial("tcp", "127.0.0.1:6379")
		defer conn.Close()
		if err != nil {
			logs.Info("redis链接错误")
		}
		//把以前相同商品的历史浏览记录删除
		conn.Do("lrem", "history_"+strconv.Itoa(user.Id), 0, id)
		//添加新的商品浏览记录
		conn.Do("lpush", "history_"+strconv.Itoa(user.Id), id)

	}

	ShowLayout(&this.Controller)
	this.TplName = "detail.html"
}

//展示商品列表页
func (this *GoodsController) ShowList() {
	id, err := this.GetInt("typeId")
	// id := this.GetString("TypeId")
	// logs.Info("id: ", id)
	if err != nil {
		logs.Info("请求路径错误", err)
		this.Redirect("/", 302)
		return
	}
	//处理数据
	ShowLayout(&this.Controller)

	o := orm.NewOrm()
	var goodsNew []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).OrderBy("Time").Limit(2, 0).All(&goodsNew)
	this.Data["goodsNew"] = goodsNew

	//获取商品
	var goods []models.GoodsSKU

	//分页实现
	//获取pageCount
	count, _ := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).Count()
	pageSize := 3
	pageCount := math.Ceil(float64(count) / float64(pageSize))

	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	pages := PageTool(int(pageCount), pageIndex)
	this.Data["pages"] = pages
	this.Data["typeId"] = id
	this.Data["pageIndex"] = pageIndex

	start := (pageIndex - 1) * pageSize
	//获取上一页页码
	prePage := pageIndex - 1
	if prePage <= 1 {
		prePage = 1
	}
	this.Data["prePage"] = prePage

	//获取下一页页码
	nextPage := pageIndex + 1
	if nextPage > int(pageCount) {
		nextPage = int(pageCount)
	}
	this.Data["nextPage"] = nextPage
	logs.Info("start: ", start)

	//按照一定顺序获取商品
	sort := this.GetString("sort")
	if sort == "" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).Limit(pageSize, start).All(&goods)
		this.Data["sort"] = ""
		this.Data["goods"] = goods

	} else if sort == "price" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).OrderBy("Price").Limit(pageSize, start).All(&goods)
		this.Data["sort"] = "price"
		this.Data["goods"] = goods
	} else {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).OrderBy("Sales").Limit(pageSize, start).All(&goods)
		this.Data["sort"] = "sale"
		this.Data["goods"] = goods
	}

	this.TplName = "list.html"
}
