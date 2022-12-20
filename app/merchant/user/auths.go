package user

import (
	"huling/utils/results"
	utils "huling/utils/tool"
	"strings"

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
	MDB := DB().Table("merchant_auth_rule").Fields("id,title,menuName,orderNo,type,parentMenu,icon,routePath,routeName,component,redirect,permission,isExt,keepalive,hideMenu,hideBreadcrumb,hideChildrenInMenu,currentActiveMenu,hideTab")
	//获取用户权限
	role_id, _ := DB().Table("merchant_auth_role_access").Where("uid", user.ID).Pluck("role_id")
	if role_id == nil {
		results.Failed(context, "您没有使用权限", ok)
		return
	}
	menu_ids, _ := DB().Table("merchant_auth_role").WhereIn("id", role_id.([]interface{})).Pluck("rules")
	getmenus := ArrayMerge(menu_ids.([]interface{}))
	nemu_list, err := MDB.WhereIn("id", getmenus).Where("status", 0).WhereIn("type", []interface{}{0, 1}).Order("orderNo asc").Get()
	if err != nil {
		results.Failed(context, "查找菜单列表失败！", err)
		return
	}
	rulenum := GetMenuArray(nemu_list, 0)
	results.Success(context, "获取用户菜单列表!", rulenum, nil)
}

/*
获取角色选择菜单
*/
func GetRoleMenuList(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	pid := context.DefaultQuery("pid", "0")
	MDB := DB().Table("merchant_auth_rule")
	var zero string = "0"
	if id == zero || pid == zero { //获取本账号所拥有的权限
		//账号信息
		getuser, _ := context.Get("user") //当前用户
		user := getuser.(*utils.UserClaims)
		role_id, _ := DB().Table("merchant_auth_role_access").Where("uid", user.ID).Pluck("role_id")
		menu_id, _ := DB().Table("merchant_auth_role").WhereIn("id", role_id.([]interface{})).Pluck("rules")
		if !IsContain(menu_id.([]interface{}), "*") { //不是超级权限-过滤菜单权限
			getmenus := ArraymoreMerge(menu_id.([]interface{}))
			MDB = MDB.WhereIn("id", getmenus)
		}
	} else {
		//获取用户权限
		pid, _ := DB().Table("merchant_auth_role").Where("id", id).Value("pid") //获取父级权限
		menu_id_str, _ := DB().Table("merchant_auth_role").Where("id", pid).Value("rules")
		if !strings.Contains(menu_id_str.(string), "*") { //不是超级权限-过滤菜单权限
			getmenus := utils.Axplode(menu_id_str)
			MDB = MDB.WhereIn("id", getmenus)
		}
	}
	menuList, _ := MDB.Fields("id,parentMenu,menuName").Order("orderNo asc").Get()
	menuList = GetMenuChildrenArray(menuList, 0)
	results.Success(context, "获取菜单数据", menuList, nil)
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
	MDB := DB().Table("merchant_auth_rule")
	//获取用户权限
	role_id, _ := DB().Table("merchant_auth_role_access").Where("uid", user.ID).Pluck("role_id")
	menu_id, _ := DB().Table("merchant_auth_role").WhereIn("id", role_id.([]interface{})).Pluck("rules")
	if !IsContain(menu_id.([]interface{}), "*") { //不是超级权限-过滤菜单权限
		getmenus := ArrayMerge(menu_id.([]interface{}))
		MDB = MDB.WhereIn("id", getmenus)
	}
	list, _ := MDB.Where("type", 2).Pluck("permission")
	results.Success(context, "获取用户权限", list, nil)
}
