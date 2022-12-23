package mwebh5

//产品
import (
	"encoding/json"
	"huling/utils/results"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

// 提交表单数据
func SaveForm(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//添加一条提交记录
	userinfo, _ := DB().Table("client_member").Where("id", parameter["uid"]).Fields("cuid,accountID").First()
	recordId, reerr := DB().Table("client_form_record").Data(map[string]interface{}{
		"cuid":       userinfo["cuid"],
		"accountID":  userinfo["accountID"],
		"member_id":  parameter["uid"],
		"form_id":    parameter["form_id"],
		"createtime": time.Now().Unix(),
	}).InsertGetId()
	if reerr != nil {
		results.Failed(context, "添加提交记录失败", reerr)
	} else {
		save_arr := []map[string]interface{}{}
		for _, val := range parameter["item_list"].([]interface{}) {
			webb, _ := json.Marshal(val)
			var webjson map[string]interface{}
			_ = json.Unmarshal(webb, &webjson)
			save_arr = append(save_arr, map[string]interface{}{
				"form_item_id": webjson["form_item_id"],
				"record_id":    recordId,
				"value":        webjson["value"],
			})
		}
		res, valerr := DB().Table("client_form_value").Data(save_arr).Insert()
		if valerr != nil {
			results.Failed(context, "添加表单值失败", valerr)
		} else {
			results.Success(context, "提交表单", res, nil)
		}
	}
}

// 获取表单字段列表
func GetFormField(context *gin.Context) {
	form_id := context.DefaultQuery("form_id", "")
	if form_id == "" {
		results.Failed(context, "请传参数form_id", nil)
	} else {
		list, err := DB().Table("client_form_item").Where("form_id", form_id).Get()
		if err != nil {
			results.Failed(context, "获取表单字段失败", err)
		} else {
			rule, Rerr := DB().Table("client_form_rule").Where("form_id", form_id).Get()
			rulelist := rule
			if Rerr != nil || rule == nil {
				rulelist = make([]gorose.Data, 0)
			}
			results.Success(context, "获取表单字段和规则", map[string]interface{}{"list": list, "rule": rulelist}, nil)
		}
	}
}
