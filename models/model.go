package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
)

type User struct { // 用户表
	Id        int
	Name      string       `orm:"size(20);unique"` // 用户名
	Password  string       `orm:"size(20)"`        // 登录密码
	Email     string       `orm:"size(50)"`        // 邮箱
	Active    bool         `orm:"default(false)"`  // 是否激活
	Power     int          `orm:"default(0)"`      // 权限设置  0 表示用户  1表示管理员
	Address   []*Address   `orm:"reverse(many)"`
	OrderInfo []*OrderInfo `orm:"reverse(many)"`
}

type Address struct { // 地址表
	Id        int
	Receiver  string       `orm:"size(20)"`      //收件人
	Address   string       `orm:"size(50)"`      //收件地址
	Addr      string       `orm:"size(50)"`      //收件地址
	Zipcode   string       `orm:"size(20)"`      //邮编
	Phone     string       `orm:"size(20)"`      //联系方式
	Isdefault bool         `orm:"defalt(false)"` //是否默认 false 为非默认  true为默认
	User      *User        `orm:"rel(fk)"`       //用户ID
	OrderInfo []*OrderInfo `orm:"reverse(many)"`
}

type Goods struct { //商品SPU表
	Id       int
	Name     string      `orm:"size(20)"`  //商品名称
	Detail   string      `orm:"size(200)"` //详细描述
	GoodsSKU []*GoodsSKU `orm:"reverse(many)"`
}
type GoodsType struct { //商品类型表
	Id                   int
	Name                 string                  //种类名称
	Logo                 string                  //logo
	Image                string                  //图片
	GoodsSKU             []*GoodsSKU             `orm:"reverse(many)"`
	IndexTypeGoodsBanner []*IndexTypeGoodsBanner `orm:"reverse(many)"`
}

// type GoodsSKU struct {
// 	Id        int
// 	Name      string
// 	Introduce string
// 	Num       int
// 	Tip       string
// 	Img       string
// 	Info      string
// 	TypeId    int
// 	SalesNum  int
// 	Counts    int
// 	SpuId     int
// }
type GoodsSKU struct { //商品SKU表
	Id                   int
	Goods                *Goods                  `orm:"rel(fk)"` //商品SPU
	GoodsType            *GoodsType              `orm:"rel(fk)"` //商品所属种类
	Name                 string                  //商品名称
	Desc                 string                  //商品简介
	Price                int                     //商品价格
	Unite                string                  //商品单位
	Image                string                  //商品图片
	Stock                int                     `orm:"default(1)"`   //商品库存
	Sales                int                     `orm:"default(0)"`   //商品销量
	Status               int                     `orm:"default(1)"`   //商品状态
	Time                 time.Time               `orm:"auto_now_add"` //添加时间
	GoodsImage           []*GoodsImage           `orm:"reverse(many)"`
	IndexGoodsBanner     []*IndexGoodsBanner     `orm:"reverse(many)"`
	IndexTypeGoodsBanner []*IndexTypeGoodsBanner `orm:"reverse(many)"`
	OrderGoods           []*OrderGoods           `orm:"reverse(many)"`
}

type GoodsImage struct { //商品图片表
	Id       int
	Image    string    //商品图片
	GoodsSKU *GoodsSKU `orm:"rel(fk)"` //商品SKU
}
type IndexGoodsBanner struct { //首页轮播商品展示表
	Id       int
	GoodsSKU *GoodsSKU `orm:"rel(fk)"` //商品sku
	Image    string    //商品图片
	Index    int       `orm:"default(0)"` //展示顺序
}

type IndexTypeGoodsBanner struct { //首页分类商品展示表
	Id          int
	GoodsType   *GoodsType `orm:"rel(fk)"`    //商品类型
	GoodsSKU    *GoodsSKU  `orm:"rel(fk)"`    //商品sku
	DisplayType int        `orm:"default(1)"` //展示类型 0代表文字，1代表图片
	Index       int        `orm:"default(0)"` //展示顺序
}

type IndexPromotionBanner struct { //首页促销商品展示表
	Id    int
	Name  string `orm:"size(20)"` //活动名称
	Url   string `orm:"size(50)"` //活动链接
	Image string //活动图片
	Index int    `orm:"default(0)"` //展示顺序
}

// type CarouseImg struct {
// 	Id     int
// 	ImgArr []*CarImgInfo
// }

// type CarImgInfo struct {
// 	Url    string
// 	ActUrl string
// 	Sort   string
// }

// type PrompProduct struct {
// 	Id     int
// 	ImgArr []*CarImgInfo
// }

// type IndexProdList struct {
// 	Id   int
// 	List []*IndexProd
// }

// type IndexProd struct {
// 	TypeId      int
// 	SpuId       int
// 	IndexTypeId int
// 	Sort        int
// }

// type IndexTypeList struct {
// 	Id   int
// 	Name string
// }

// 订单表
type OrderInfo struct { //订单表
	Id           int
	OrderId      string    `orm:"unique"`
	User         *User     `orm:"rel(fk)"` //用户
	Address      *Address  `orm:"rel(fk)"` //地址
	PayMethod    int       //付款方式 1已支付 2未支付 3
	TotalCount   int       `orm:"default(1)"` //商品数量
	TotalPrice   int       //商品总价
	TransitPrice int       //运费
	Orderstatus  int       `orm:"default(1)"`   //订单状态 //1已结算 2未结算 3无效
	TradeNo      string    `orm:"default('')"`  //支付编号
	Time         time.Time `orm:"auto_now_add"` //评论时间

	OrderGoods []*OrderGoods `orm:"reverse(many)"`
}
type OrderGoods struct { //订单商品表
	Id        int
	OrderInfo *OrderInfo `orm:"rel(fk)"`    //订单
	GoodsSKU  *GoodsSKU  `orm:"rel(fk)"`    //商品
	Count     int        `orm:"default(1)"` //商品数量
	Price     int        //商品价格
	Comment   string     `orm:"default('')"` //评论
}

func init() {
	l := logs.GetLogger()
	l.Println("this is a message of http")
	logs.GetLogger("ORM").Println("this is a message of orm")
	logs.Debug("my book is bought in the year of ", 2016)
	logs.Info("this %s cat is %v years old", "yellow", 3)
	logs.Warn("json is a type of kv like", map[string]int{"key": 2016})
	logs.Error(1024, "is a very", "good game")
	logs.Critical("oh,crash")

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

	orm.RegisterModel(new(User), new(Address), new(OrderGoods), new(OrderInfo), new(IndexPromotionBanner), new(IndexTypeGoodsBanner), new(IndexGoodsBanner), new(GoodsImage), new(GoodsSKU), new(GoodsType), new(Goods))

	orm.RunSyncdb("default", false, true)
}
