package mwebh5

//产品
import (
	"encoding/json"
	"fmt"
	"huling/utils/results"

	"github.com/gin-gonic/gin"
)

// 获取详情内容
func GetProduct(context *gin.Context) {
	id := context.DefaultQuery("id", "")
	micweb_id := context.DefaultQuery("micweb_id", "")
	data, err := DB().Table("client_product_manage").Where("id", id).Fields("id,cuid,type,title,des,images,releasetime,content,createtime").First()
	if err != nil {
		results.Failed(context, "获取产品详情内容失败", err)
	} else {
		if data == nil {
			results.Success(context, "产品内容不存在！", nil, nil)
		} else {
			prolist, _ := DB().Table("client_product_manage_pro").Where("cuid", data["cuid"]).Where("status", 0).Fields("id,keyname,name,des,weigh,type").Order("weigh asc").Get()
			for _, pro := range prolist {
				pro_val, _ := DB().Table("client_product_manage_pro_val").Where("product_id", id).Where("pro_id", pro["id"]).Value("val")
				if pro_val != nil {
					pro["val"] = pro_val
				} else {
					pro["val"] = ""
				}
			}
			data["prolist"] = prolist
			//图片
			rooturl, _ := DB().Table("merchant_config").Where("keyname", "rooturl").Value("keyvalue")
			if data["images"] != "" {
				//多图
				var parameter []interface{}
				_ = json.Unmarshal([]byte(data["images"].(string)), &parameter)
				var newimg []interface{}
				for _, img := range parameter {
					img = fmt.Sprintf("%s%s", rooturl, img)
					newimg = append(newimg, img)
				}
				data["images"] = newimg
			} else {
				data["images"] = make([]interface{}, 0)
			}
			//标签
			lids, _ := DB().Table("client_product_lid").Where("product_id", id).Pluck("lid")
			labels, _ := DB().Table("client_product_label").WhereIn("id", lids.([]interface{})).Pluck("name")
			data["labels"] = labels
			//获取客服
			service, _ := DB().Table("client_product_service").Where("product_id", id).Where("status", 0).Order("id desc").Get()
			data["service"] = service
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
			results.Success(context, "获取产品详情内容", data, nil)
		}
	}
}
