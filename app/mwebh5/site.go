package mwebh5

import (
	"encoding/json"
	"huling/utils/results"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

// 获取h5页面数据
func GetHome(context *gin.Context) {
	//当前网站id
	id := context.DefaultQuery("id", "")
	if id == "0" {
		results.Failed(context, "请传网站id", nil)
	} else {
		micweb, _ := DB().Table("client_micweb").Where("id", id).Fields("id,cuid,title,des,footer_tabbar,top_tabbar,side_tabbar,status,updatetime").First()
		var pagedata gorose.Data
		pagedata_f, _ := DB().Table("client_micweb_page").Where("micweb_id", id).Where("ishome", 1).First()
		pagedata = pagedata_f
		if pagedata == nil {
			pagedata_DF, _ := DB().Table("client_micweb_page").Where("micweb_id", id).First()
			pagedata = pagedata_DF
		}
		if pagedata != nil {
			//字符串转JSON
			if pagedata["component"] != nil {
				pagedata["component"] = StingToJSON(pagedata["component"])
			}
			//字符串转JSON
			if pagedata["templateJson"] != nil {
				var parameter interface{}
				_ = json.Unmarshal([]byte(pagedata["templateJson"].(string)), &parameter)
				pagedata["templateJson"] = parameter
			}

			if micweb != nil {
				if micweb["footer_tabbar"] != nil {
					var parameterf interface{}
					_ = json.Unmarshal([]byte(micweb["footer_tabbar"].(string)), &parameterf)
					micweb["footer_tabbar"] = parameterf
				}
			}
			pagedata["micweb"] = micweb
			//判断是否审核
			pagedata["updatetime"] = time.Now().Unix() - micweb["updatetime"].(int64)
			if micweb["status"] == 2 {
				pagedata["audit"] = true
			} else {
				pagedata["audit"] = false
			}
			//判断是否开通支付
			ispay, _ := DB().Table("client_system_paymentconfig").Where("cuid", micweb["cuid"]).Value("mchAPIv3Key")
			if ispay != nil && ispay != "" {
				micweb["ispay"] = true
			} else {
				micweb["ispay"] = false
			}
		}
		results.Success(context, "获取首页数据", pagedata, nil)
	}
}

// 获取单个页面数据
func GetPage(context *gin.Context) {
	//当前页面uuid
	id := context.DefaultQuery("id", "0")
	uuid := context.DefaultQuery("uuid", "0")
	if id == "0" {
		results.Failed(context, "请传网站的：id", nil)
	} else if uuid == "0" {
		results.Failed(context, "请传页面的：uuid", nil)
	} else {
		pagedata, _ := DB().Table("client_micweb_page").Where("micweb_id", id).Where("uuid", uuid).First()
		if pagedata != nil {
			//字符串转JSON
			if pagedata["component"] != nil {
				pagedata["component"] = StingToJSON(pagedata["component"])
			}
			//字符串转JSON
			if pagedata["templateJson"] != nil {
				var parameter interface{}
				_ = json.Unmarshal([]byte(pagedata["templateJson"].(string)), &parameter)
				pagedata["templateJson"] = parameter
			}
			//获取完整信息
			micweb, _ := DB().Table("client_micweb").Where("id", pagedata["micweb_id"]).Fields("id,cuid,title,des,footer_tabbar,top_tabbar,side_tabbar,status,updatetime").First()
			if micweb != nil {
				if micweb["footer_tabbar"] != nil {
					var parameterf interface{}
					_ = json.Unmarshal([]byte(micweb["footer_tabbar"].(string)), &parameterf)
					micweb["footer_tabbar"] = parameterf
				}
			}
			//判断是否开通支付
			ispay, _ := DB().Table("client_system_paymentconfig").Where("cuid", micweb["cuid"]).Value("mchAPIv3Key")
			if ispay != nil && ispay != "" {
				micweb["ispay"] = true
			} else {
				micweb["ispay"] = false
			}
			pagedata["micweb"] = micweb
		}
		results.Success(context, "获取页面数据", pagedata, nil)
	}
}

// 新增访问记录
func AddVisitRecord(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if _, ok := parameter["micweb_id"]; !ok {
		results.Failed(context, "请提交：micweb_id", nil)
	} else {
		micweb, _ := DB().Table("client_micweb").Where("id", parameter["micweb_id"]).Fields("id,cuid,accountID").First()
		if micweb != nil {
			parameter["cuid"] = micweb["cuid"]
			parameter["accountID"] = micweb["accountID"]
		}
		parameter["createtime"] = time.Now().Unix()
		addId, err := DB().Table("client_micweb_visitlog").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加访问失败", err)
		} else {
			results.Success(context, "添加访问成功！", addId, nil)
		}
	}
}
