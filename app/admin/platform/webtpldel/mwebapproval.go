package webtpldel

import (
	"encoding/json"
	"huling/utils/results"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取数据列表-分页先用Limit在page分则无小-type=1是平台
func Getlist(context *gin.Context) {
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	title := context.DefaultQuery("title", "")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := DB().Table("client_micweb_tpl_main").Where("isdel", 1)
	FDB := DB().Table("client_micweb_tpl_main").Where("isdel", 1)
	if title != "" {
		MDB = MDB.Where("title", "like", "%"+title+"%")
		FDB = FDB.Where("title", "like", "%"+title+"%")
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("delTime desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			userinfo, _ := DB().Table("client_user").Where("id", val["cuid"]).Fields("mobile,name").First()
			val["userinfo"] = userinfo
			del_cuiduser, _ := DB().Table("client_user").Where("id", val["del_cuid"]).Fields("mobile,name").First()
			val["deluser"] = del_cuiduser
		}
		var totalCount int64
		totalCount, _ = FDB.Count()
		results.Success(context, "获取轻模板列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 真删除
func Del(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("client_micweb_tpl_main").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 恢复模板
func RestoreTpl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_micweb_tpl_main").Where("id", parameter["id"]).Data(map[string]interface{}{"isdel": 0}).Update()
	if err != nil {
		results.Failed(context, "恢复失败！", err)
	} else {
		results.Success(context, "恢复成功！", res2, nil)
	}
}
