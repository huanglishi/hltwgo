package role

import (
	"encoding/json"
	"fmt"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

// 获取数据列表-子树结构
func Getlist(context *gin.Context) {
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	role_id, _ := DB().Table("merchant_auth_role_access").Where("uid", user.ID).Pluck("role_id")
	role_ids := GetAllChilIds(role_id.([]interface{})) //批量获取子节点id
	all_role_id := mergeArr(role_id.([]interface{}), role_ids)
	fmt.Printf("当前计算机CPU核数: %v ", user.ID)
	roleList, _ := DB().Table("merchant_auth_role").WhereIn("id", all_role_id).OrWhere("uid", user.ID).Order("weigh asc").Get()
	roleList = GetTreeArray(roleList, 0, "")
	if roleList == nil {
		roleList = make([]gorose.Data, 0)
	}
	results.Success(context, "获取拥有角色列表", roleList, nil)
}

// 获取父级数据
func GetParentList(context *gin.Context) {
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	role_id, _ := DB().Table("merchant_auth_role_access").Where("uid", user.ID).Pluck("role_id")
	role_ids := GetAllChilIds(role_id.([]interface{})) //批量获取子节点id
	all_role_id := mergeArr(role_id.([]interface{}), role_ids)
	list, _ := DB().Table("merchant_auth_role").WhereIn("id", all_role_id).Fields("id,pid,name").Order("weigh asc").Get()
	list = GetMenuChildrenArray(list, 0)
	if list == nil {
		list = make([]gorose.Data, 0)
	}
	results.Success(context, "角色父级数据！", list, nil)
}

// 表单选择数据-单选用
func GetAllList(context *gin.Context) {
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	role_id, _ := DB().Table("merchant_auth_role_access").Where("uid", user.ID).Pluck("role_id")
	role_ids := GetAllChilIds(role_id.([]interface{})) //批量获取子节点id
	all_role_id := mergeArr(role_id.([]interface{}), role_ids)
	menuList, _ := DB().Table("merchant_auth_role").WhereIn("id", all_role_id).Where("status", 0).Fields("id as value,pid,name as label").Order("weigh asc").Get()
	menuList = GetMenuChildrenArraylist(menuList, 0)
	results.Success(context, "表单选择角色数据", menuList, nil)
}

// 在表单选择数据-多选用
func GetRoleList(context *gin.Context) {
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	role_id, _ := DB().Table("merchant_auth_role_access").Where("uid", user.ID).Pluck("role_id")
	role_ids := GetAllChilIds(role_id.([]interface{})) //批量获取子节点id
	all_role_id := mergeArr(role_id.([]interface{}), role_ids)
	menuList, _ := DB().Table("merchant_auth_role").WhereIn("id", all_role_id).Where("status", 0).Fields("id as value,pid,name as label").Order("weigh asc").Get()
	results.Success(context, "表单选择角色多选用数据", menuList, nil)
}

// 添加数据
func Add(context *gin.Context) {
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
	if parameter["menu"] != nil && parameter["menu"] != "*" {
		rules := GetRulesID(parameter["menu"]) //获取子菜单包含的父级ID
		rudata := rules.([]interface{})
		var rulesStr []string
		for _, v := range rudata {
			str := fmt.Sprintf("%v", v) //interface{}强转string
			rulesStr = append(rulesStr, str)
		}
		parameter["rules"] = strings.Join(rulesStr, ",")
		parameter["menu"] = JSONMarshalToString(parameter["menu"])
	}
	parameter["createtime"] = time.Now().Unix()
	if f_id == 0 {
		parameter["accountID"] = user.Accountid
		addId, err := DB().Table("merchant_auth_role").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 {
				DB().Table("merchant_auth_role").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("merchant_auth_role").
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
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("merchant_auth_role").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

// 删除菜单
func Del(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("merchant_auth_role").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除菜单失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
