package tplaplication

import (
	"encoding/json"
	"huling/utils/results"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取数据列表-分页先用Limit在page分则无小-type=1是平台
func Getlist(context *gin.Context) {
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	title := context.DefaultQuery("title", "")
	status := context.DefaultQuery("status", "0")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := DB().Table("client_micweb_customtpl")
	CMDB := DB().Table("client_micweb_customtpl")
	if status != "0" {
		MDB = MDB.Where("status", status)
		CMDB = CMDB.Where("status", status)
	}
	if title != "" {
		MDB = MDB.Where("name", "like", "%"+title+"%")
		CMDB = CMDB.Where("name", "like", "%"+title+"%")
	}

	list, err := MDB.Limit(pageSize).Page(pageNo).Fields("id,name,des,status,createtime").Order("id asc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			userinfo, _ := DB().Table("merchant_user").Where("id", val["uid"]).Fields("mobile,name").First()
			val["userinfo"] = userinfo
		}
		var totalCount int64
		totalCount, _ = CMDB.Count()
		results.Success(context, "获取申请模板列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 获取单条数据详情
func GetDetail(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	data, err := DB().Table("client_micweb_customtpl").Where("id", id).First()
	if err != nil {
		results.Failed(context, "获取单条数据详情败", err)
	} else {
		results.Success(context, "获取单条数据详情", data, nil)
	}
}

// 3获取模板
func GetTplList(context *gin.Context) {
	keyword := context.DefaultQuery("keyword", "")
	MDB := DB().Table("client_micweb_tpl_main")
	if keyword != "" {
		MDB = MDB.Where("title", "like", "%"+keyword+"%")
	}
	data, err := MDB.Fields("id,title").Get()
	if err != nil {
		results.Failed(context, "获取模板情败", err)
	} else {
		results.Success(context, "获取模板列表", data, nil)
	}
}

// 处理申请
func DoResult(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_micweb_customtpl").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"],
		"backtime": time.Now().Unix(), "backcontent": parameter["backcontent"], "tpl_id": parameter["tpl_id"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		results.Success(context, "提交成功！", res2, nil)
	}
}
