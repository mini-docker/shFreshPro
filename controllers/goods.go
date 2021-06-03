package controllers

import (
	"shFreshPro/models"

	"github.com/beego/beego/v2/client/orm"
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
	this.TplName = "index.html"
}
