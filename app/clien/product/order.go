package product

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取订单
func GetOrderList(context *gin.Context) {
	title := context.DefaultQuery("title", "")
	status := context.DefaultQuery("status", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	whereMap := DB().Table("client_product_order").Where("cuid", user.ClientID)
	whereMap2 := DB().Table("client_product_order").Where("cuid", user.ClientID)
	if status != "0" {
		whereMap.Where("status", status)
		whereMap2.Where("status", status)
	}
	if title != "" {
		whereMap.Where("title", "like", "%"+title+"%")
		whereMap2.Where("title", "like", "%"+title+"%")
	}
	list, err := whereMap.Limit(pageSize).Page(pageNo).Fields("id,title,price,number,product_id,total_fee,out_trade_no,note,address_id,address,logistics_name,logistics_mobile,createtime,paytime,status").Order("createtime desc , id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		var totalCount int64
		totalCount, _ = whereMap2.Count()
		results.Success(context, "获取产品订单列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 获取订单
func GetOrder(context *gin.Context) {
	id := context.DefaultQuery("id", "")
	data, err := DB().Table("client_product_order").Where("id", id).First()
	if err != nil {
		results.Failed(context, "获取订单失败", err)
	} else {
		//发货地址
		address, _ := DB().Table("client_member_address").Where("id", data["address_id"]).Fields("name,mobile,province_name,city_name,area_name,address").First()
		data["address"] = address
		results.Success(context, "获取订单", data, nil)
	}
}

// 更新订单状态
func UpOrder(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_product_order").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res2, nil)
	}
}
