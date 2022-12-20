package teacharticle

import (
	"encoding/json"
	"huling/utils/results"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取文章列表
func GetList(context *gin.Context) {
	title := context.DefaultQuery("title", "")
	status := context.DefaultQuery("status", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	whereMap := DB().Table("admin_article")
	whereMap2 := DB().Table("admin_article")
	if status != "0" {
		whereMap.Where("status", status)
		whereMap2.Where("status", status)
	}
	if title != "" {
		whereMap.Where("title", "like", "%"+title+"%")
		whereMap2.Where("title", "like", "%"+title+"%")
	}
	list, err := whereMap.Limit(pageSize).Page(pageNo).Fields("id,type,title,des,status,createtime").Order("id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		var totalCount int64
		totalCount, _ = whereMap2.Count()
		results.Success(context, "获取文章列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 添加文章
func SaveArticle(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	var f_id float64 = 0
	if parameter["id"] != nil {
		f_id = parameter["id"].(float64)
	}
	delete(parameter, "pendingStatus")
	if f_id == 0 {
		parameter["createtime"] = time.Now().Unix()
		addId, err := DB().Table("admin_article").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 {
				DB().Table("admin_article").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("admin_article").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			results.Success(context, "更新成功！", res, nil)
		}
	}
}

// 更新状态
func UpLock(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("admin_article").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res2, nil)
	}
}

// 删除
func DelArticle(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("admin_article").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 获取文章大编辑内容
func GetArticle(context *gin.Context) {
	id := context.DefaultQuery("id", "")
	data, err := DB().Table("admin_article").Where("id", id).First()
	if err != nil {
		results.Failed(context, "获取文章大内容字段失败", err)
	} else {
		results.Success(context, "文章大内容字段", data, nil)
	}
}
