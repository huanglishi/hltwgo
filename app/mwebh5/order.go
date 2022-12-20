package mwebh5

import (
	"encoding/json"
	"fmt"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取订单列表
func GetOrderList(context *gin.Context) {
	uid := context.DefaultQuery("uid", "0")
	if uid == "0" {
		results.Failed(context, "请传参数uid", nil)
	} else {
		list, err := DB().Table("client_product_order").Where("uid", uid).Fields("id,status,out_trade_no,product_id,title,number,price,total_fee,address,logistics_name,logistics_mobile,note").Order("id desc").Get()
		if err != nil {
			results.Failed(context, "获取订单失败！", err)
		} else {
			for _, val := range list {
				images, _ := DB().Table("client_product_manage").Where("id", val["product_id"]).Value("images")
				if images != nil && images != "" {
					rooturl, _ := DB().Table("merchant_config").Where("keyname", "rooturl").Value("keyvalue")
					//多图
					var parameter []interface{}
					_ = json.Unmarshal([]byte(images.(string)), &parameter)
					var newimg []interface{}
					for _, img := range parameter {
						img = fmt.Sprintf("%s%s", rooturl, img)
						newimg = append(newimg, img)
					}
					val["images"] = newimg
				} else {
					val["images"] = make([]interface{}, 0)
				}
			}
			results.Success(context, "获取订单列表", list, nil)
		}
	}
}

// 获取订单详情
func GetOrderDetail(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	if id == "0" {
		results.Failed(context, "请传参数id", nil)
	} else {
		data, err := DB().Table("client_product_order").Where("id", id).First()
		if err != nil {
			results.Failed(context, "获取订单详情失败！", err)
		} else {
			images, _ := DB().Table("client_product_manage").Where("id", data["product_id"]).Value("images")
			if images != nil && images != "" {
				rooturl, _ := DB().Table("merchant_config").Where("keyname", "rooturl").Value("keyvalue")
				//多图
				var parameter []interface{}
				_ = json.Unmarshal([]byte(images.(string)), &parameter)
				var newimg []interface{}
				for _, img := range parameter {
					img = fmt.Sprintf("%s%s", rooturl, img)
					newimg = append(newimg, img)
				}
				data["images"] = newimg
			} else {
				data["images"] = make([]interface{}, 0)
			}
			//地址
			address, _ := DB().Table("client_member_address").Where("id", data["address_id"]).Fields("name,mobile,province_name,city_name,area_name,address").First()
			data["address"] = address
			results.Success(context, "获取订单详情", data, nil)
		}
	}
}

// 用户下单
func AddOrder(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if _, ok := parameter["product_id"]; !ok || parameter["product_id"] == "" {
		results.Failed(context, "请传参数:product_id（产品id）", nil)
	} else {
		userinfo, _ := DB().Table("client_member").Where("id", parameter["uid"]).Fields("cuid,accountID").First()
		parameter["cuid"] = userinfo["cuid"]
		parameter["accountID"] = userinfo["accountID"]
		parameter["createtime"] = time.Now().Unix()
		parameter["out_trade_no"] = utils.GetSnowflakeId()
		res, err := DB().Table("client_product_order").Data(parameter).Insert()
		if err != nil {
			results.Failed(context, "用户下单失败！", err)
		} else {
			results.Success(context, "用户下单成功1", res, userinfo)
		}
	}
}
