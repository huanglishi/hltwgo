package cancelapp

import (
	"encoding/json"
	"fmt"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 核销订单操作
func DoCancel(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if _, ok := parameter["cancel_no"]; !ok {
		results.Failed(context, "核销码无效", nil)
	} else {
		order_cancel, err := DB().Table("client_product_order_cancel").Where("cancel_no", parameter["cancel_no"]).First()
		if err != nil {
			results.Failed(context, "订单核销失败！", err)
		} else {
			if order_cancel["build_time"] == nil {
				results.Failed(context, "核销码无效", nil)
				return
			}
			//判断时间
			cancel_valid, _ := DB().Table("merchant_config").Where("keyname", "cancel_valid").Value("keyvalue") //获取有效时间
			nowtime := time.Now().Unix()
			cancel_validint, _ := strconv.ParseInt(cancel_valid.(string), 10, 64)
			if (nowtime - order_cancel["build_time"].(int64)) > cancel_validint*60 { //无效重新生成-已过期
				results.Success(context, "核销码已过期", "codeInvalid", nil)
			} else { //执行核销操作
				//当前用户
				getuser, _ := context.Get("user")
				user := getuser.(*utils.UserClaims)
				_, err := DB().Table("client_product_order").Data(map[string]interface{}{"status": 9, "cancel_time": nowtime, "cancel_cuid": user.ID}).Where("id", order_cancel["order_id"]).Update()
				if err != nil {
					results.Failed(context, "核销码失败", err)
				} else {
					out_trade_no, _ := DB().Table("client_product_order").Where("id", order_cancel["order_id"]).Value("out_trade_no")
					results.Success(context, "订单核销成功", "cancelSuccess", out_trade_no)
				}
			}
		}
	}
}

// 获取核销记录
func CancelRecord(context *gin.Context) {
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	whereMap := DB().Table("client_product_order").Where("cancel_cuid", user.ID).Where("status", 9)
	whereMap2 := DB().Table("client_product_order").Where("cancel_cuid", user.ID).Where("status", 9)
	list, err := whereMap.Limit(pageSize).Page(pageNo).Fields("id,product_id,title,number,price,cancel_time").Order("cancel_time desc , id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			images, _ := DB().Table("client_product_manage").Where("id", val["product_id"]).Value("images")
			if images != nil && images != "" {
				rooturl, _ := DB().Table("merchant_config").Where("keyname", "rooturl").Value("keyvalue")
				//多图
				var parameter []interface{}
				_ = json.Unmarshal([]byte(images.(string)), &parameter)
				// var newimg []interface{}
				// for _, img := range parameter {
				// 	img = fmt.Sprintf("%s%s", rooturl, img)
				// 	newimg = append(newimg, img)
				// }
				val["image"] = fmt.Sprintf("%s%s", rooturl, parameter[0])
			} else {
				val["image"] = ""
			}
		}
		var totalCount int64
		totalCount, _ = whereMap2.Count()
		results.Success(context, "获取核销记录", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}
