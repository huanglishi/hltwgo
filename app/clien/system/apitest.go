package system

import (
	"encoding/json"
	"huling/utils/Toolconf"
	"huling/utils/results"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取文章列表
func GetApiList(context *gin.Context) {
	types := context.DefaultQuery("type", "")
	keyword := context.DefaultQuery("keyword", "")
	url := context.DefaultQuery("url", "")
	cid := context.DefaultQuery("cid", "0")
	status := context.DefaultQuery("status", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	whereMap := DB().Table("common_apitest").Where("type", types)
	whereMap2 := DB().Table("common_apitest").Where("type", types)
	if status != "0" {
		whereMap.Where("status", status)
		whereMap2.Where("status", status)
	}
	if cid != "0" {
		whereMap.Where("cid", cid)
		whereMap2.Where("cid", cid)
	}
	if keyword != "" {
		whereMap.Where("title", "like", "%"+keyword+"%")
		whereMap2.Where("title", "like", "%"+keyword+"%")
	}
	if url != "" {
		whereMap.Where("url", "like", "%"+url+"%")
		whereMap2.Where("url", "like", "%"+url+"%")
	}
	list, err := whereMap.Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			//分组
			group, _ := DB().Table("common_apitest_group").Where("id", val["cid"]).Fields("name,useFrom").First()
			val["catename"] = group["name"]
			val["useFrom"] = group["useFrom"]
		}
		var totalCount int64
		totalCount, _ = whereMap2.Count()
		results.Success(context, "获取列表1", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 获取数据库字段
func GetDBField(context *gin.Context) {
	tablename := context.DefaultQuery("tablename", "")
	tablename_arr := strings.Split(tablename, ",")
	dbname := Toolconf.AppConfig.String("db.name")
	var dielddata_list []map[string]interface{}
	for _, Val := range tablename_arr {
		dielddata, _ := DB().Query("select COLUMN_NAME,COLUMN_COMMENT,COLUMN_TYPE,DATA_TYPE,CHARACTER_MAXIMUM_LENGTH,IS_NULLABLE,COLUMN_DEFAULT,NUMERIC_PRECISION from information_schema.columns where TABLE_SCHEMA='" + dbname + "' AND TABLE_NAME='" + Val + "'")
		for _, data := range dielddata {
			data["tablename"] = Val
			dielddata_list = append(dielddata_list, data)
		}
	}
	results.Success(context, "获取数据库字段", dielddata_list, tablename)
	// if err != nil {
	// 	results.Failed(context, "获取数据库字段失败", err)
	// } else {
	// 	results.Success(context, "获取数据库字段", dielddata_list, tablename)
	// }
}

// 添加
func SaveData(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	var f_id float64 = 0
	if parameter["id"] != nil {
		f_id = parameter["id"].(float64)
	}
	delete(parameter, "catename")
	delete(parameter, "pendingStatus")
	if f_id == 0 {
		parameter["createtime"] = time.Now().Unix()
		addId, err := DB().Table("common_apitest").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("common_apitest").
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
	res2, err := DB().Table("common_apitest").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
func DelData(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("common_apitest").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
