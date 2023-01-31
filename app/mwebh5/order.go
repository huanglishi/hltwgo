package mwebh5

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

// 获取订单列表
func GetOrderList(context *gin.Context) {
	uid := context.DefaultQuery("uid", "0")
	if uid == "0" {
		results.Failed(context, "请传参数uid", nil)
	} else {
		list, err := DB().Table("client_product_order").Where("uid", uid).Fields("id,status,prepay_id,prepay_time,out_trade_no,product_id,title,number,price,total_fee,address,logistics_name,logistics_mobile,note").Order("id desc").Get()
		if err != nil {
			results.Failed(context, "获取订单失败！", err)
		} else {
			for _, val := range list {
				product, _ := DB().Table("client_product_manage").Where("id", val["product_id"]).Fields("images,type").First()
				val["type"] = product["type"]
				images := product["images"]
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
				if val["prepay_id"] != "" && val["prepay_id"] != nil {
					//判断过期
					timestamp := strconv.FormatInt(time.Now().Unix(), 10) //10位时间戳
					//当前时间戳转int
					intNum, _ := strconv.Atoi(timestamp)
					timestampint := int64(intNum)
					prepay_time_int := val["prepay_time"].(int64) //数据库的时间传戳
					diff := timestampint - prepay_time_int        //
					gethour := diff / 3600
					var towhur int64 = 2
					if gethour > towhur { //支付已过期
						val["status"] = 12
					}
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
			product, _ := DB().Table("client_product_manage").Where("id", data["product_id"]).Fields("images,type,cancel_des").First()
			data["type"] = product["type"]
			data["cancel_des"] = product["cancel_des"]
			images := product["images"]
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
			//判断订单过期
			if data["prepay_id"] != "" && data["prepay_id"] != nil {
				//判断过期
				timestamp := strconv.FormatInt(time.Now().Unix(), 10) //10位时间戳
				//当前时间戳转int
				intNum, _ := strconv.Atoi(timestamp)
				timestampint := int64(intNum)
				prepay_time_int := data["prepay_time"].(int64) //数据库的时间传戳
				diff := timestampint - prepay_time_int         //
				gethour := diff / 3600
				var towhur int64 = 2
				if gethour > towhur { //支付已过期
					data["status"] = 12
				}
			}
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
		getid, err := DB().Table("client_product_order").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "用户下单失败！", err)
		} else {
			results.Success(context, "用户下单成功1", map[string]interface{}{"order_id": getid, "out_trade_no": parameter["out_trade_no"]}, nil)
		}
	}
}

// 获取核销码
func GetCancelNo(context *gin.Context) {
	order_id := context.DefaultQuery("order_id", "0")
	if order_id == "0" {
		results.Failed(context, "请传参数id", nil)
	} else {
		cancel_no := utils.GetSnowflakeId()
		adddata := map[string]interface{}{"cancel_no": cancel_no, "order_id": order_id, "build_time": time.Now().Unix()}
		_, err := DB().Table("client_product_order_cancel").Data(adddata).Insert()
		if err != nil {
			results.Failed(context, "用户下单失败！", err)
		} else {
			cancel_valid, _ := DB().Table("merchant_config").Where("keyname", "cancel_valid").Value("keyvalue")
			adddata["cancel_valid"] = cancel_valid
			results.Success(context, "获取订单核新销码1", adddata, nil)
		}
	}
}

// 检查是否已经扫码
func GetIsCancel(context *gin.Context) {
	order_id := context.DefaultQuery("order_id", "0")
	cancel_no := context.DefaultQuery("cancel_no", "0")
	if order_id == "0" {
		results.Failed(context, "请传参数订单id：order_id", nil)
	} else if cancel_no == "0" {
		results.Failed(context, "请传参数核销码：cancel_no", nil)
	} else {
		status, gerr := DB().Table("client_product_order").Where("id", order_id).Value("status")
		if gerr != nil || status == nil {
			results.Failed(context, "获取核销状态失败", gerr)
		} else {
			var nisnum int64 = 9
			if status.(int64) == nisnum {
				results.Success(context, "订已核销", "cancel", nil)
			} else {
				//判断时间
				cancel_valid, _ := DB().Table("merchant_config").Where("keyname", "cancel_valid").Value("keyvalue") //获取有效时间
				nowtime := time.Now().Unix()
				cancel_validint, _ := strconv.ParseInt(cancel_valid.(string), 10, 64)
				order_cancel, _ := DB().Table("client_product_order_cancel").Where("cancel_no", cancel_no).Order("build_time desc").Fields("cancel_no,build_time").First()
				if (nowtime - order_cancel["build_time"].(int64)) > cancel_validint*60 { //无效重新生成-已过期
					results.Success(context, "订单未核销:核销码已过期", "codeInvalid", nil)
				} else {
					results.Success(context, "订单未核销:核销码有效", "codeValid", status)
				}

			}
		}
	}
}
