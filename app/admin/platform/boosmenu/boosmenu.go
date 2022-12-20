package boosmenu

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取菜单-子树结构
func Getlist(context *gin.Context) {
	menuList, _ := DB().Table("merchant_auth_rule").Order("orderNo asc").Get()
	menuList = GetMenuChildrenArray(menuList, 0)
	results.Success(context, "获取全部菜单列表", menuList, nil)
}

// 获取菜单父级数据
func GetParentList(context *gin.Context) {
	menuList, _ := DB().Table("merchant_auth_rule").WhereIn("type", []interface{}{0, 1}).Fields("id,parentMenu,menuName").Order("orderNo asc").Get()
	menuList = GetMenuChildrenArray(menuList, 0)
	// rulenum := GetTreeArray(rule, 0, "")
	// // rulenum := make([]interface{}, 0)
	// list_text := getTreeList_txt(rulenum, "title")
	// maxid, _ := DB().Table("merchant_auth_rule").Max("id")
	results.Success(context, "菜单父级数据！", menuList, nil)
}

// 角色授权-获取菜单
func GetMenuList(context *gin.Context) {
	menuList, _ := DB().Table("merchant_auth_rule").Fields("id,parentMenu,menuName").Order("orderNo asc").Get()
	menuList = GetMenuChildrenArray(menuList, 0)
	results.Success(context, "获取菜单数据", menuList, nil)
}

// 添加菜单
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
	parameter["createtime"] = time.Now().Unix()
	if f_id == 0 {
		addId, err := DB().Table("merchant_auth_rule").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加菜单失败", err)
		} else {
			if addId != 0 {
				DB().Table("merchant_auth_rule").
					Data(map[string]interface{}{"orderNo": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("merchant_auth_rule").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新菜单失败", err)
		} else {
			results.Success(context, "更新成功！", res, nil)
		}
	}
}

// 更新状态
func UpMenuLock(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("merchant_auth_rule").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
	res2, err := DB().Table("merchant_auth_rule").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除菜单失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
