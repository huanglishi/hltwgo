package mwebh5

import (
	"encoding/json"
	"huling/utils/results"

	"github.com/gin-gonic/gin"
)

// 获取文章详情内容
func GetArticle(context *gin.Context) {
	id := context.DefaultQuery("id", "")
	micweb_id := context.DefaultQuery("micweb_id", "")
	data, err := DB().Table("client_article_manage").Where("id", id).Fields("id,type,title,des,author,image,releasetime,content").First()
	if err != nil {
		results.Failed(context, "获取文章详情内容失败", err)
	} else {
		//获取完整信息
		if micweb_id != "" {
			micweb, _ := DB().Table("client_micweb").Where("id", micweb_id).Fields("id,cuid,title,des,footer_tabbar,status,updatetime").First()
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
			data["micweb"] = micweb
		}
		results.Success(context, "获取文章详情内容", data, nil)
	}
}
