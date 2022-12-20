package home

import (
	"huling/utils/results"
	utils "huling/utils/tool"

	"github.com/gin-gonic/gin"
)

// 获取统计数据
func GetModels(context *gin.Context) {
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	//1获取账号总数
	userdata, _ := DB().Table("merchant_user").Where("id", user.ID).Fields("id,sub_account").First()
	pid, _ := DB().Table("merchant_packagedesign").Where("funkey", "microweb").Value("id")
	user_package, _ := DB().Table("merchant_user_package").Where("accountID", user.Accountid).Where("pid", pid).Fields("id,number").First()
	//2获取已经使用
	//2.1已经建子账号
	chridCount, _ := DB().Table("merchant_user").Where("pid", user.ID).Count()
	//2.2获取已使用的微站
	micwebCount, _ := DB().Table("merchant_micweb").Where("accountID", user.Accountid).Count()
	//拼接数组
	menuList := []map[string]interface{}{}
	//轻站统计
	var microwebnum int64 = 0
	var microwebnumvalue int64 = 0
	if user_package["number"] != nil {
		microwebnum = user_package["number"].(int64)
		microwebnumvalue = microwebnum - micwebCount
	}

	menuList = append(menuList, map[string]interface{}{
		"funkey":      "microweb",
		"title":       "轻站数",
		"icon":        "visit-count|svg",
		"value":       microwebnumvalue, //剩余数量
		"valueprefix": "剩余 ",
		"valuesuffix": "",
		"total":       microwebnum, //总数量
		"totalfix":    "",
		"totalsuffix": "",
		"color":       "green",
		"action":      "续费/续量",
	})
	//子账号
	var chridnum int64 = 0
	var chridnumvalue int64 = 0
	if userdata["sub_account"] != nil {
		chridnum = userdata["sub_account"].(int64)
		chridnumvalue = chridnum - chridCount
	}
	menuList = append(menuList, map[string]interface{}{
		"funkey":      "subaccounts",
		"title":       "子账号数",
		"icon":        "total-sales|svg",
		"value":       chridnumvalue, //剩余数量
		"valueprefix": "剩余 ",
		"valuesuffix": "个",
		"total":       chridnum, //总数量
		"totalfix":    "",
		"totalsuffix": "个",
		"color":       "blue",
		"action":      "",
	})
	results.Success(context, "获取统计数据", menuList, nil)
}
