package dashboard

import (
	"encoding/json"
	"huling/utils/helper"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取统计数据
func GetNumList(context *gin.Context) {
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
	//附件
	// menuList = append(menuList, map[string]interface{}{
	// 	"funkey":      "attachment",
	// 	"title":       "附件",
	// 	"icon":        "download-count|svg",
	// 	"value":       200, //剩余数量
	// 	"valueprefix": "剩余 ",
	// 	"valuesuffix": " M",
	// 	"total":       120000, //总数量
	// 	"totalfix":    "",
	// 	"totalsuffix": " M",
	// 	"color":       "blue",
	// 	"action":      "",
	// })
	results.Success(context, "获取统计数据", menuList, nil)
}

// 获取轻站预置支付套餐
func GetWebPayList(context *gin.Context) {
	funkey := context.DefaultQuery("funkey", "")
	packagedata, _ := DB().Table("merchant_packagedesign").Where("funkey", funkey).Fields("id,name,price,count").First()
	paylist, _ := DB().Table("merchant_packagedesign_paylist").Where("packg_id", packagedata["id"]).Get()
	results.Success(context, "获取统计数据", map[string]interface{}{"paylist": paylist, "packagedata": packagedata}, nil)
}

// 轻站续费购买订单
func SavaPayOrder(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	parameter["accountID"] = user.Accountid
	parameter["createtime"] = time.Now().Unix()
	parameter["updatetime"] = time.Now().Unix()
	parameter["orderid"] = helper.GetSnowflakeId1()
	addId, err := DB().Table("merchant_order").Data(parameter).InsertGetId()
	if err != nil {
		results.Failed(context, "添加失败", err)
	} else {
		results.Success(context, "添加成功！", addId, nil)
	}
}

// 获取剩余轻站数
func IsHaseMicroweb(context *gin.Context) {
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	//1获总数
	pid, _ := DB().Table("merchant_packagedesign").Where("funkey", "microweb").Value("id")
	user_package, _ := DB().Table("merchant_user_package").Where("accountID", user.Accountid).Where("pid", pid).Fields("id,number").First()
	//2获取已经使用
	//2.2获取已使用的微站
	micwebCount, _ := DB().Table("merchant_micweb").Where("accountID", user.Accountid).Count()
	results.Success(context, "获取剩余轻站数", (user_package["number"].(int64) - micwebCount), nil)
}

// 获取剩余子账号数
func IsHaseSubAccount(context *gin.Context) {
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	//1获取账号总数
	userdata, _ := DB().Table("merchant_user").Where("id", user.ID).Fields("id,sub_account").First()
	//2已经建子账号
	chridCount, _ := DB().Table("merchant_user").Where("pid", user.ID).Count()
	results.Success(context, "获取剩余轻站数", (userdata["sub_account"].(int64) - chridCount), nil)
}
