package system

import (
	"encoding/json"
	"huling/utils/results"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

// 获取分类列表
func GetGroupList(context *gin.Context) {
	types := context.DefaultQuery("type", "")
	title := context.DefaultQuery("title", "")
	status := context.DefaultQuery("status", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//当前用户
	wheremap := DB().Table("common_apitest_group").Where("type", types)
	wheremap2 := DB().Table("common_apitest_group").Where("type", types)
	if status != "0" {
		wheremap.Where("status", status)
		wheremap2.Where("status", status)
	}
	if title != "" {
		wheremap.Where("title", "like", "%"+title+"%")
		wheremap2.Where("title", "like", "%"+title+"%")
	}
	list, err := wheremap.
		Order("id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			//上一级
			pidname, _ := DB().Table("common_apitest_group").Where("id", val["pid"]).Value("name")
			if pidname != nil {
				val["pidname"] = pidname
			} else {
				val["pidname"] = "无"
			}
		}
		var totalCount int64
		totalCount, _ = wheremap2.Count()
		rulenum := GetTreeArray(list, 0, "")
		list_text := GetTreeList_txt(rulenum, "name")
		results.Success(context, "获取分类列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list_text,
		}, nil)
	}
}

// 添加分类
func SaveGroup(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)

	var f_id float64 = 0
	if parameter["id"] != nil {
		f_id = parameter["id"].(float64)
	}
	if f_id == 0 {
		addId, err := DB().Table("common_apitest_group").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("common_apitest_group").
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

// 添加表单使用
func GetFormGroupList(context *gin.Context) {
	keyword := context.DefaultQuery("keyword", "")
	types := context.DefaultQuery("type", "")
	MDB := DB().Table("common_apitest_group").Where("type", types)
	if keyword != "" {
		MDB = MDB.Where("name", "like", "%"+keyword+"%")
	}
	list, _ := MDB.Fields("id,pid,name").Order("id desc").Get()
	if list == nil {
		list = make([]gorose.Data, 0)
		results.Success(context, "分组数据！", list, nil)
	} else {
		rulenum := GetTreeArray(list, 0, "")
		list_text := GetTreeList_txt(rulenum, "name")
		results.Success(context, "分组数据！", list_text, nil)
	}
}

// 删除
func DelGroup(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("common_apitest_group").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 更新状态
func UpStatus(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("common_apitest_group").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
