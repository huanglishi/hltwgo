package home

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

// 1 获取客户提交的模板需求
func GetCustomtpl(context *gin.Context) {
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	list, err := DB().Table("client_micweb_customtpl").Where("cuid", user.ClientID).Fields("id,name,tpl_id,backcontent,createtime,backtime").Get()
	if err != nil {
		results.Failed(context, "提交的模板列表失败", err)
	} else {
		if list == nil {
			list = make([]gorose.Data, 0)
		}
		results.Success(context, "获取提交的模板列表", list, nil)
	}
}

// 1 获取客户提交的模板需求-详情
func GetCustomtpldeltail(context *gin.Context) {
	id := context.DefaultQuery("id", "")
	data, err := DB().Table("client_micweb_customtpl").Where("id", id).First()
	if err != nil {
		results.Failed(context, "模板需求详情失败", err)
	} else {
		results.Success(context, "模板需求详情", data, nil)
	}
}

// 2保存客户提交的模板需求
func SaveCustom(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	parameter["cuid"] = user.ClientID
	parameter["accountID"] = user.Accountid
	parameter["createtime"] = time.Now().Unix()
	res2, err := DB().Table("client_micweb_customtpl").Data(parameter).Insert()
	if err != nil {
		results.Failed(context, "需求提交失败！", err)
	} else {
		results.Success(context, "需求提交已提交成功，请等待处理！", res2, nil)
	}
}
