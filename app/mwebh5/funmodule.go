package mwebh5

import (
	"encoding/json"
	"fmt"
	"huling/utils/results"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

// 全站搜索
func SearchAll(context *gin.Context) {
	micweb_id := context.DefaultQuery("micweb_id", "")
	keyword := context.DefaultQuery("keyword", "")
	if micweb_id == "" {
		results.Failed(context, "请传站点ID参数：micweb_id", nil)
	} else {
		micweb, _ := DB().Table("client_micweb").Where("id", micweb_id).Fields("cuid").First()
		AMD := DB().Table("client_article_manage").Where("cuid", micweb["cuid"])
		PMD := DB().Table("client_product_manage").Where("cuid", micweb["cuid"])
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
			//图片
			rooturl, _ := DB().Table("merchant_config").Where("keyname", "rooturl").Value("keyvalue")
			prolist, _ := DB().Table("client_product_manage_pro").Where("cuid", micweb["cuid"]).Where("status", 0).Fields("id,keyname,name,des,weigh,type").Order("weigh asc").Get()
			for _, val := range plist {
				if val["images"] != "" {
					//多图
					var parameter []interface{}
					_ = json.Unmarshal([]byte(val["images"].(string)), &parameter)
					var newimg []interface{}
					for _, img := range parameter {
						img = fmt.Sprintf("%s%s", rooturl, img)
						newimg = append(newimg, img)
					}
					val["images"] = newimg
				} else {
					val["images"] = make([]interface{}, 0)
				}
				var myprolist []gorose.Data
				for _, pro := range prolist {
					pro_val, _ := DB().Table("client_product_manage_pro_val").Where("product_id", val["id"]).Where("pro_id", pro["id"]).Value("val")
					if pro_val != nil {
						pro["val"] = pro_val
					} else {
						pro["val"] = ""
					}
					myprolist = append(myprolist, map[string]interface{}{"id": pro["id"], "keyname": pro["keyname"], "name": pro["name"], "des": pro["des"], "weigh": pro["weigh"], "type": pro["type"], "val": pro["val"]})
				}
				val["prolist"] = myprolist
			}
			results.Success(context, "全站搜索", map[string]interface{}{"alist": alist, "plist": plist}, nil)
		}
	}

}
