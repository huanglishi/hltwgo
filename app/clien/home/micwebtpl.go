package home

/**
轻站模板
*/
import (
	"encoding/json"
	"huling/utils/results"
	"io/ioutil"
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
	MDB := DB().Table("client_micweb_tpl_main").Where("isdel", 0).Where("flag", "home")
	if key_cid != "0" {
		MDB = MDB.Where("cid", key_cid)
	}
	datalist, _ := MDB.Limit(100).Order("id desc").Get()
	locall_imgurl, _ := DB().Table("merchant_config").Where("keyname", "locall_imgurl").Value("keyvalue")
	for _, val := range datalist {
		if val["image"] != "" && val["image"] != nil {
			val["image"] = strings.Replace(val["image"].(string), "http://192.168.1.118:8098", locall_imgurl.(string), -1)
		}
		//分类
		catename, _ := DB().Table("client_micweb_tpl_group").Where("id", val["cid"]).Value("name")
		val["catename"] = catename
	}
	results.Success(context, "获取网站模板数据", datalist, nil)
}

// 设为首页推广-首页展示
func SetHomeview(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if _, ok := parameter["tplid"]; !ok {
		results.Failed(context, "请传参数模板id:tplid", nil)
		return
	}
	if _, ok := parameter["flag"]; !ok {
		parameter["flag"] = ""
	}
	res, err := DB().Table("client_micweb_tpl_main").Data(map[string]interface{}{"flag": parameter["flag"]}).Where("id", parameter["tplid"]).Update()
	if err != nil {
		results.Failed(context, "设置失败", err)
	} else {
		results.Success(context, "设置成功！", res, nil)
	}
}
