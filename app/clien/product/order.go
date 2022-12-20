package product

import (
	"huling/utils/results"
	utils "huling/utils/tool"
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
	list, err := whereMap.Limit(pageSize).Page(pageNo).Fields("id,price,number,product_id,total_fee,out_trade_no,note,address_id,address,logistics_name,logistics_mobile,createtime,paytime").Order("createtime desc , id desc").Get()
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
