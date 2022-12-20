package ordermanag

import (
	"huling/utils/results"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取套餐购买数据
func GetOrderList(context *gin.Context) {
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := DB().Table("merchant_order")
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("id asc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		var totalCount int64
		totalCount, _ = MDB.Count()
		results.Success(context, "获取套餐购买数据列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}
