package product

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// 获取参数列表
func GetProlist(context *gin.Context) {
	pro_id := context.DefaultQuery("pro_id", "0")
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("client_product_manage_pro_list").Where("cuid", user.ClientID).Where("pro_id", pro_id).Order("weigh asc").Get()
	for _, val := range list {
		val["edit"] = false
	}
	results.Success(context, "产品参数数据！", list, nil)
}

// 保存/修改数据
func SaveProlist(context *gin.Context) {
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
	delete(parameter, "id")
	delete(parameter, "edit")
	if f_id == 0 {
		parameter["cuid"] = user.ClientID
		parameter["accountID"] = user.Accountid
		addId, err := DB().Table("client_product_manage_pro_list").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败！", err)
		} else {
			if addId != 0 {
				DB().Table("client_product_manage_pro_list").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		_, err := DB().Table("client_product_manage_pro_list").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			results.Success(context, "更新成功！", f_id, user)
		}
	}
}

// 删除
func DelProlist(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("client_product_manage_pro_list").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
}

// 更新状态
func UpProlist(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_product_manage_pro_list").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

// 更新排序
func UpWeighlist(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res1, err := DB().Table("client_product_manage_pro_list").Where("id", parameter["id"]).Data(map[string]interface{}{"weigh": parameter["rpweigh"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		DB().Table("client_product_manage_pro_list").Where("id", parameter["rpId"]).Data(map[string]interface{}{"weigh": parameter["weigh"]}).Update()
		msg := "更新成功！"
		if res1 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res1, nil)
	}
}
