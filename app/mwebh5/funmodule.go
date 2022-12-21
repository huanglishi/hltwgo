package mwebh5

import (
	"encoding/json"
	"huling/utils/results"

	"github.com/gin-gonic/gin"
)

// 全站搜索
func SearchAll(context *gin.Context) {
	micweb_id := context.DefaultQuery("micweb_id", "")
	keyword := context.DefaultQuery("keyword", "")
	if micweb_id == "" {
		results.Failed(context, "请传站点ID参数：micweb_id", nil)
	} else {
		AMD := DB().Table("client_article_manage")
		PMD := DB().Table("client_product_manage")
		if keyword != "" {
			AMD.Where("title", "like", "%"+keyword+"%")
			PMD.Where("title", "like", "%"+keyword+"%")
		}
		alist, err := AMD.Where("status", 0).Fields("id,type,title,link,views,image,releasetime").Order("id desc").Get()
		plist, perr := PMD.Where("status", 0).Fields("id,type,title,des,views,images,releasetime").Order("id desc").Get()
		if err != nil {
			results.Failed(context, "查找文章数据出错", err)
		} else if perr != nil {
			results.Failed(context, "查找产品数据出错", perr)
		} else {
			for _, val := range plist {
				if val["images"] != "" {
					//多图
					var parameter []interface{}
					_ = json.Unmarshal([]byte(val["images"].(string)), &parameter)
					val["images"] = parameter
				} else {
					val["images"] = make([]interface{}, 0)
				}
			}
			results.Success(context, "全站搜索", map[string]interface{}{"alist": alist, "plist": plist}, nil)
		}
	}

}
