package mwebh5

import (
	"huling/utils/results"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取h5页面数据
func GetData(context *gin.Context) {
	//当前用户
	id := context.DefaultQuery("id", "0")
	data, _ := DB().Table("merchant_micweb_item").Where("id", id).First()
	if data != nil {
		//获取首页id
		if data["ishome"] == 1 {
			data["homeId"] = data["id"]
		} else {
			home_id, _ := DB().Table("merchant_micweb_item").Where("micweb_id", data["micweb_id"]).Where("ishome", 1).Value("id")
			data["homeId"] = home_id
		}
		//判断是否审核
		micweb, _ := DB().Table("merchant_micweb").Where("id", data["micweb_id"]).Fields("status,updatetime").First()
		data["updatetime"] = time.Now().Unix() - micweb["updatetime"].(int64)
		if micweb["status"] == 2 {
			data["audit"] = true
		} else {
			data["audit"] = false
		}
	}
	results.Success(context, "获取页面数据", data, nil)
}

// 获取h5页面预览数据
func GetPreviewData(context *gin.Context) {
	//当前用户
	id := context.DefaultQuery("id", "0")
	main_id := context.DefaultQuery("main_id", "0")
	home_id, _ := DB().Table("merchant_micweb_tpl_main").Where("id", main_id).Value("home_id")
	var item_id interface{}
	if main_id == id {
		item_id = home_id
	} else {
		item_id = id
	}
	data, _ := DB().Table("merchant_micweb_tpl_main_page").Where("main_id", main_id).Where("item_id", item_id).First()
	results.Success(context, "获取页面数据", data, nil)
}
