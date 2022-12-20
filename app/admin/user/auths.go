package user

import (
	"huling/utils/results"
	utils "huling/utils/tool"

	"github.com/gin-gonic/gin"
)

/*
获取用户菜单
*/
func GetMenuList(context *gin.Context) {
	getuser, ok := context.Get("user") //取值 实现了跨中间件取值
	if !ok {
		results.Failed(context, "用户id不存在！", ok)
		return
	}
	user := getuser.(*utils.UserClaims)
	MDB := DB().Table("admin_auth_rule").Fields("id,title,menuName,orderNo,type,parentMenu,icon,routePath,routeName,component,redirect,permission,isExt,keepalive,hideMenu,hideBreadcrumb,hideChildrenInMenu,currentActiveMenu,hideTab")
	//获取用户权限
	role_id, _ := DB().Table("admin_auth_role_access").Where("uid", user.ID).Pluck("role_id")
	menu_id, _ := DB().Table("admin_auth_role").WhereIn("id", role_id.([]interface{})).Pluck("rules")
	if !IsContain(menu_id.([]interface{}), "*") { //不是超级权限-过滤菜单权限
		getmenus := ArrayMerge(menu_id.([]interface{}))
		MDB = MDB.WhereIn("id", getmenus)
	}
	nemu_list, err := MDB.Where("status", 0).WhereIn("type", []interface{}{0, 1}).Order("orderNo asc").Get()
	if err != nil {
		results.Failed(context, "查找菜单列表失败！", err)
		return
	}
	rulenum := GetMenuArray(nemu_list, 0)
	results.Success(context, "获取用户菜单列表!", rulenum, nil)
}

/*
获取用户权限
*/
func GetPermCode(context *gin.Context) {
	getuser, ok := context.Get("user") //取值 实现了跨中间件取值
	if !ok {
		results.Failed(context, "用户id不存在！", ok)
		return
	}
	user := getuser.(*utils.UserClaims)
	MDB := DB().Table("admin_auth_rule")
	//获取用户权限
	role_id, _ := DB().Table("admin_auth_role_access").Where("uid", user.ID).Pluck("role_id")
	menu_id, _ := DB().Table("admin_auth_role").WhereIn("id", role_id.([]interface{})).Pluck("rules")
	if !IsContain(menu_id.([]interface{}), "*") { //不是超级权限-过滤菜单权限
		getmenus := ArrayMerge(menu_id.([]interface{}))
		MDB = MDB.WhereIn("id", getmenus)
	}
	list, _ := MDB.Where("type", 2).Pluck("permission")
	results.Success(context, "获取用户权限", list, nil)
}
