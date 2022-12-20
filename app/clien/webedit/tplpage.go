package webedit

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

// 1 获取单页模板分类
func GetTplPageGroup(context *gin.Context) {
	typestr := context.DefaultQuery("type", "0")
	MBD := DB().Table("client_micweb_tplpage_group").Where("status", 0)
	if typestr != "0" {
		MBD.Where("type", typestr)
	}
	dList, err := MBD.Fields("id,name,type,remark").Order("weigh asc").Get()
	if err != nil {
		results.Failed(context, "获取模板分类失败", err)
	} else {
		results.Success(context, "获取模板分类", dList, nil)
	}
}

// 2 获取单页模板
func GetTplpage(context *gin.Context) {
	group_id := context.DefaultQuery("group_id", "0")
	MBD := DB().Table("client_micweb_tplpage_page")
	if group_id != "0" {
		MBD.Where("group_id", group_id)
	}
	list, err := MBD.Order("id desc").Get()
	if err != nil {
		results.Failed(context, "获取单页模板失败", err)
	} else {
		locall_imgurl, _ := DB().Table("merchant_config").Where("keyname", "locall_imgurl").Value("keyvalue")
		for _, val := range list {
			if val["image"] != "" && val["image"] != nil {
				val["image"] = strings.Replace(val["image"].(string), "http://192.168.1.118:8098", locall_imgurl.(string), -1)
			}
		}
		results.Success(context, "获取单页模板", list, nil)
	}
}

// 3 添加单页模板
func SaveTplpage(context *gin.Context) {
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
	//JSON转字符串
	if parameter["component"] != nil {
		parameter["component"] = JSONMarshalToString(parameter["component"])
	}
	//JSON转字符串
	if parameter["templateJson"] != nil {
		parameter["templateJson"] = JSONMarshalToString(parameter["templateJson"])
	}
	if f_id == 0 {
		parameter["cuid"] = user.ClientID
		parameter["accountID"] = user.Accountid
		addId, err := DB().Table("client_micweb_tplpage_page").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("client_micweb_tplpage_page").
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
func DelTplpage(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_micweb_tplpage_page").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
}
