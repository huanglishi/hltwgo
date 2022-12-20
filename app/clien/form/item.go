package form

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// 获取表单项列表
func GetItemList(context *gin.Context) {
	form_id := context.DefaultQuery("form_id", "0")
	MDB := DB().Table("client_form_item").Where("form_id", form_id)
	list, _ := MDB.Order("weigh asc").Get()
	// for _, val := range list {
	// }
	results.Success(context, "获取表单项数据", list, nil)
}

// 保存/修改数据
func SaveItem(context *gin.Context) {
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
		addId, err := DB().Table("client_form_item").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败！", err)
		} else {
			if addId != 0 {
				DB().Table("client_form_item").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		_, err := DB().Table("client_form_item").
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
func DelItem(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("client_form_item").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
}

// 更新状态
func UpItem(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_form_item").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

// 更新必填状态
func UpRequired(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_form_item").Where("id", parameter["id"]).Data(map[string]interface{}{"required": parameter["required"]}).Update()
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
func UpWeigh(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res1, err := DB().Table("client_form_item").Where("id", parameter["id"]).Data(map[string]interface{}{"weigh": parameter["rpweigh"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		DB().Table("client_form_item").Where("id", parameter["rpId"]).Data(map[string]interface{}{"weigh": parameter["weigh"]}).Update()
		msg := "更新成功！"
		if res1 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res1, nil)
	}
}
