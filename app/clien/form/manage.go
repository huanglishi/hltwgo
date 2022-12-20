package form

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取列表
func GetList(context *gin.Context) {
	uname := context.DefaultQuery("name", "")
	status := context.DefaultQuery("status", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	whereMap := DB().Table("client_form").Where("cuid", user.ClientID)
	whereMap2 := DB().Table("client_form").Where("cuid", user.ClientID)
	if status != "0" {
		whereMap.Where("status", status)
		whereMap2.Where("status", status)
	}
	if uname != "" {
		whereMap.Where("name", "like", "%"+uname+"%")
		whereMap2.Where("name", "like", "%"+uname+"%")
	}
	list, err := whereMap.Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		var totalCount int64
		totalCount, _ = whereMap2.Count()
		results.Success(context, "获取表单列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 添加
func SaveForm(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	var f_id float64 = 0
	if parameter["id"] != nil {
		f_id = parameter["id"].(float64)
	}
	formItem := parameter["formItem"]
	delete(parameter, "formItem")
	if f_id == 0 {
		parameter["createtime"] = time.Now().Unix()
		parameter["accountID"] = user.Accountid
		parameter["cuid"] = user.ClientID
		addId, err := DB().Table("client_form").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			//批量新增表单项
			addformItems(formItem.([]interface{}), addId)
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("client_form").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			//批量新增表单项
			addformItems(formItem.([]interface{}), f_id)
			results.Success(context, "更新成功！", f_id, res)
		}
	}
}

// 批量新增表单项
func addformItems(list []interface{}, form_id interface{}) {
	//批量提交
	// save_arr := []map[string]interface{}{}
	for _, val := range list {
		webb, _ := json.Marshal(val)
		var webjson map[string]interface{}
		_ = json.Unmarshal(webb, &webjson)
		if GetInterfaceToInt(webjson["id"]) == 0 {
			webjson["form_id"] = form_id
			delete(webjson, "id")
			gid, _ := DB().Table("client_form_item").Data(webjson).InsertGetId()
			DB().Table("client_form_item").Data(map[string]interface{}{"weigh": gid}).Where("id", gid).Update()
			// save_arr = append(save_arr, webjson)
		} else {
			DB().Table("client_form_item").Data(webjson).Where("id", webjson["id"]).Update()
		}
	}

}

// 更新状态
func UpLock(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_form").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
func Del(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("client_form").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		//同时删除表单项
		DB().Table("client_form_item").WhereIn("form_id", ids.([]interface{})).Delete()
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
