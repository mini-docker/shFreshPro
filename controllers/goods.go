package controllers

import (
	"encoding/json"
	"shFreshPro/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
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

	logs.Info("goods: ", b)
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
