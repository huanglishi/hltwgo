package form

import (
	"encoding/json"
	"huling/utils/results"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// 获取表单规则和选项数据
func GetRuleAndselectData(context *gin.Context) {
	form_id := context.DefaultQuery("form_id", "0")
	//获取规则
	rulelist, _ := DB().Table("client_form_rule").Where("form_id", form_id).Get()
	formitemlist, ierr := DB().Table("client_form_item").Where("form_id", form_id).Get()
	if ierr != nil {
		results.Failed(context, "获取表单规则和选项数据失败！", ierr)
	} else {
		for _, val := range rulelist {
			val["showadd"] = true
			val["isdow"] = true
			if val["show_item_ids"] != nil {
				val["show_item_ids"] = StingToJSON(val["show_item_ids"])
			}
			if val["show_item_text"] != nil {
				val["show_item_text"] = StingToJSON(val["show_item_text"])
			}
			if val["selectval"] != nil {
				val["selectval"] = StingToJSON(val["selectval"])
			}
		}
		results.Success(context, "获取表单规则和选项数据", map[string]interface{}{"rulelist": rulelist, "formitemlist": formitemlist}, nil)
	}
}

// 保存
func SaveRule(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter []map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	var gerro error
	var resdata []int
	for _, val := range parameter {
		webb, _ := json.Marshal(val)
		var webjson map[string]interface{}
		_ = json.Unmarshal(webb, &webjson)
		delete(webjson, "showadd")
		delete(webjson, "isdow")
		if webjson["show_item_ids"] != nil {
			webjson["show_item_ids"] = JSONMarshalToString(webjson["show_item_ids"])
		}
		if webjson["show_item_text"] != nil {
			webjson["show_item_text"] = JSONMarshalToString(webjson["show_item_text"])
		}
		if webjson["selectval"] != nil {
			webjson["selectval"] = JSONMarshalToString(webjson["selectval"])
		}
		if GetInterfaceToInt(webjson["id"]) == 0 {
			delete(webjson, "id")
			res, err := DB().Table("client_form_rule").Data(webjson).InsertGetId()
			gerro = err
			resdata = append(resdata, GetInterfaceToInt(res))
		} else {
			res, err := DB().Table("client_form_rule").Data(webjson).Where("id", webjson["id"]).Update()
			gerro = err
			resdata = append(resdata, GetInterfaceToInt(res))
		}
	}
	if gerro != nil {
		results.Failed(context, "更新失败", gerro)
	} else {
		results.Success(context, "更新成功！", resdata, nil)
	}
}

// 删除
func DelRuleItem(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("client_form_rule").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
}
