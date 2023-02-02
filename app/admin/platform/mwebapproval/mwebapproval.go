package mwebapproval

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
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
	status := context.DefaultQuery("status", "1")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := DB().Table("client_micweb").Where("status", status)
	if title != "" {
		MDB = MDB.Where("title", "like", "%"+title+"%")
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("publishtime desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			userinfo, _ := DB().Table("merchant_user").Where("id", val["accountID"]).Fields("mobile,name").First()
			val["userinfo"] = userinfo
			if val["des"] == "null" {
				val["des"] = ""
			}
		}
		var totalCount int64
		totalCount, _ = MDB.Count()
		results.Success(context, "获取轻站审批数据列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 获取套餐数据
func GetParentList(context *gin.Context) {
	menuList, _ := DB().Table("merchant_packagedesign").Fields("id,name").Order("id asc").Get()
	results.Success(context, "获取套餐数据", menuList, nil)
}

// 添加套餐
func SaveData(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	var f_id float64 = 0
	if parameter["id"] != nil {
		f_id = parameter["id"].(float64)
	}
	if parameter["menu"] != nil {
		parameter["menu"] = JSONMarshalToString(parameter["menu"])
	}
	parameter["createtime"] = time.Now().Unix()
	if f_id == 0 {
		addId, err := DB().Table("merchant_packagedesign").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("merchant_packagedesign").
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

// 删除
func Del(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("merchant_packagedesign").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 审批轻站
func ApprovalMweb(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_micweb").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"],
		"updatetime": time.Now().Unix(), "approval_err": parameter["approval_err"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		results.Success(context, "提交成功！", res2, nil)
	}
}
