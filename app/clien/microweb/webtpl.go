package microweb

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 8.1 删除整站模板-假删除
func DelWebTpl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	res2, err := DB().Table("client_micweb_tpl_main").Where("id", parameter["id"]).Data(map[string]interface{}{"isdel": 1, "del_cuid": user.ID, "delTime": time.Now().Unix()}).Update()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 8.1 删除整站模板-真删
func DelWebTpl_real(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_micweb_tpl_main").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		//删除子页面
		DB().Table("client_micweb_tpl_main_page").Where("main_id", parameter["id"]).Delete()
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 8.2 获取选择模板分类
func GetSelectTplGroup(context *gin.Context) {
	dList, _ := DB().Table("client_micweb_tpl_group").Where("status", 0).Fields("id,name,pid,remark").Order("weigh asc").Get()
	dList = utils.GetTreeArray(dList, 0, "")
	results.Success(context, "获取选择模板分类", dList, nil)
}

// 8.3 获取选择模板
func GetSelectTplList(context *gin.Context) {
	key_pcid := context.DefaultQuery("pcid", "0")
	key_cid := context.DefaultQuery("cid", "0")
	MDB := DB().Table("client_micweb_tpl_main").Where("isdel", 0)
	if key_cid != "0" && key_pcid != key_cid {
		MDB = MDB.Where("cid", key_cid)
	} else {
		cids, _ := DB().Table("client_micweb_tpl_group").Where("status", 0).Where("pid", key_pcid).Pluck("id")
		MDB = MDB.WhereIn("cid", cids.([]interface{}))
	}
	datalist, _ := MDB.Order("id desc").Get()
	locall_imgurl, _ := DB().Table("merchant_config").Where("keyname", "locall_imgurl").Value("keyvalue")
	for _, val := range datalist {
		if val["image"] != "" && val["image"] != nil {
			val["image"] = strings.Replace(val["image"].(string), "http://192.168.1.118:8098", locall_imgurl.(string), -1)
		}
	}
	results.Success(context, "获取网站模板数据！", datalist, nil)
}
