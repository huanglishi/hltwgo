package mwebh5

import (
	"encoding/json"
	"huling/utils/results"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取地址列表
func GetAddressList(context *gin.Context) {
	uid := context.DefaultQuery("uid", "")
	list, err := DB().Table("client_member_address").Where("uid", uid).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "更新失败", err)
	} else {
		results.Success(context, "获取地址列表", list, nil)
	}
}

// 获取地址列表
func GetAddress(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	list, err := DB().Table("client_member_address").Where("id", id).First()
	if err != nil {
		results.Failed(context, "获取地址失败", err)
	} else {
		results.Success(context, "获取地址数据", list, nil)
	}
}

// 添加地址
func SaveAddress(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	var f_id float64 = 0
	if parameter["id"] != nil {
		f_id = parameter["id"].(float64)
	}
	if f_id == 0 {
		userinfo, _ := DB().Table("client_member").Where("id", parameter["uid"]).Fields("cuid,accountID").First()
		parameter["cuid"] = userinfo["cuid"]
		parameter["accountID"] = userinfo["accountID"]
		parameter["createtime"] = time.Now().Unix()
		addId, err := DB().Table("client_member_address").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("client_member_address").
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
func DelAddress(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_member_address").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
}
