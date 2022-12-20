package mwebh5

import (
	"encoding/json"
	"huling/utils/results"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

// 1 获取预览模板-入口
func GetWebtpl(context *gin.Context) {
	tpl_id := context.DefaultQuery("tpl_id", "0")
	if tpl_id == "0" {
		results.Failed(context, "请传页面主表：tpl_id(模板信息的id就是保存封面那个表)", nil)
	} else {
		var pagedata gorose.Data
		pagedata_home, _ := DB().Table("client_micweb_tpl_main_page").Where("main_id", tpl_id).Where("ishome", 1).First()
		if pagedata_home == nil {
			pagedata_def, _ := DB().Table("client_micweb_tpl_main_page").Where("main_id", tpl_id).First()
			pagedata = pagedata_def
		} else {
			pagedata = pagedata_home
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
			//获取完整信息
			micweb, tpl_main_err := DB().Table("client_micweb_tpl_main").Where("id", tpl_id).Fields("id,title,footer_tabbar").First()
			if tpl_main_err != nil {
				results.Failed(context, "获取主表信息错误", tpl_main_err)
			} else {
				if micweb != nil {
					if micweb["footer_tabbar"] != nil {
						var parameterf interface{}
						_ = json.Unmarshal([]byte(micweb["footer_tabbar"].(string)), &parameterf)
						micweb["footer_tabbar"] = parameterf
					}
				}
				pagedata["micweb"] = micweb
			}
		}
		results.Success(context, "获取页面数据", pagedata, nil)
	}

}

// 2 获取单页模板-操作跳转
func GetWebtplPage(context *gin.Context) {
	tpl_id := context.DefaultQuery("tpl_id", "0")
	uuid := context.DefaultQuery("uuid", "0")
	if tpl_id == "0" {
		results.Failed(context, "请传页面主表：tpl_id(模板信息的id就是保存封面那个表)", nil)
	} else if uuid == "0" {
		results.Failed(context, "请传页面的：uuid", nil)
	} else {
		pagedata, _ := DB().Table("client_micweb_tpl_main_page").Where("main_id", tpl_id).Where("uuid", uuid).First()
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
			micweb, tpl_main_err := DB().Table("client_micweb_tpl_main").Where("id", tpl_id).Fields("id,title,footer_tabbar").First()
			if tpl_main_err != nil {
				results.Failed(context, "获取主表信息错误", tpl_main_err)
			} else {
				if micweb != nil {
					if micweb["footer_tabbar"] != nil {
						var parameterf interface{}
						_ = json.Unmarshal([]byte(micweb["footer_tabbar"].(string)), &parameterf)
						micweb["footer_tabbar"] = parameterf
					}
				}
				pagedata["micweb"] = micweb
			}
		}
		results.Success(context, "获取页面数据", pagedata, nil)
	}

}
