package repayment

import (
	"encoding/json"
	"huling/utils/results"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取数据列表-分页先用Limit在page分则无小-type=1是平台
func Getlist(context *gin.Context) {
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	title := context.DefaultQuery("title", "")
	status := context.DefaultQuery("status", "0")
	orderid := context.DefaultQuery("orderid", "")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := DB().Table("merchant_order")
	CMDB := DB().Table("merchant_order")
	if status != "0" {
		MDB = MDB.Where("status", status)
		CMDB = CMDB.Where("status", status)
	}
	if orderid != "" {
		MDB = MDB.Where("orderid", orderid)
		CMDB = CMDB.Where("orderid", orderid)
	}
	if title != "" {
		MDB = MDB.Where("title", "like", "%"+title+"%")
		CMDB = CMDB.Where("title", "like", "%"+title+"%")
	}

	list, err := MDB.Limit(pageSize).Page(pageNo).Order("id asc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			userinfo, _ := DB().Table("merchant_user").Where("id", val["accountID"]).Fields("mobile,name").First()
			val["userinfo"] = userinfo
		}
		var totalCount int64
		totalCount, _ = CMDB.Count()
		results.Success(context, "获取续费支付列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 处理订单
func DoResult(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("merchant_order").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"],
		"updatetime": time.Now().Unix(), "remark": parameter["remark"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		results.Success(context, "提交成功！", res2, nil)
	}
}
