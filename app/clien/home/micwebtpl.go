package home

/**
轻站模板
*/
import (
	"huling/utils/results"
	"strings"

	"github.com/gin-gonic/gin"
)

// 获取分组
func GetTplGroup(context *gin.Context) {
	dList, _ := DB().Table("client_micweb_tpl_group").Where("status", 0).Where("pid", ">", 0).Limit(12).Fields("id,name,pid,remark").Order("weigh asc").Get()
	results.Success(context, "获取选择模板分类", dList, nil)
}

// 获取模板
func GetTpl(context *gin.Context) {
	key_cid := context.DefaultQuery("cid", "0")
	MDB := DB().Table("client_micweb_tpl_main")
	if key_cid != "0" {
		MDB = MDB.Where("cid", key_cid)
	}
	datalist, _ := MDB.Limit(6).Order("id desc").Get()
	locall_imgurl, _ := DB().Table("merchant_config").Where("keyname", "locall_imgurl").Value("keyvalue")
	for _, val := range datalist {
		if val["image"] != "" && val["image"] != nil {
			val["image"] = strings.Replace(val["image"].(string), "http://192.168.1.118:8098", locall_imgurl.(string), -1)
		}
	}
	results.Success(context, "获取网站模板数据", datalist, nil)
}
