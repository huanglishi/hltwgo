package account

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取数据列表-分页先用Limit
func Getlist(context *gin.Context) {
	deptId := context.DefaultQuery("deptId", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := DB().Table("admin_user").Fields("id,status,name,username,nickname,avatar,telephone,email,remark,deptid,remark,createtime")
	if deptId != "0" {
		MDB.Where("deptid", deptId)
	}
	list, _ := MDB.Limit(pageSize).Page(pageNo).Order("id asc").Get()
	for _, val := range list {
		roleid, _ := DB().Table("admin_auth_role_access").Where("uid", val["id"]).Pluck("role_id")
		rolename, _ := DB().Table("admin_auth_role").WhereIn("id", roleid.([]interface{})).Pluck("name")
		val["rolename"] = rolename
		val["roleid"] = roleid
		depname, _ := DB().Table("admin_auth_dept").Where("id", val["deptid"]).Value("name")
		val["depname"] = depname
	}
	var totalCount int64
	totalCount, _ = DB().Table("admin_user").Count()
	results.Success(context, "获取全部列表", map[string]interface{}{
		"page":     pageNo,
		"pageSize": pageSize,
		"total":    totalCount,
		"items":    list}, nil)
}

// 获取父级数据
func GetParentList(context *gin.Context) {
	menuList, _ := DB().Table("admin_user").Fields("id,pid,name").Order("id asc").Get()
	menuList = GetMenuChildrenArray(menuList, 0)
	results.Success(context, "部门父级数据！", menuList, nil)
}

// 添加添加
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
	var roleid []interface{}
	if parameter["roleid"] != nil {
		roleid = parameter["roleid"].([]interface{})
		delete(parameter, "roleid")
	}
	if parameter["password"] != nil && parameter["password"] != "" {
		rnd := rand.New(rand.NewSource(6))
		salt := strconv.Itoa(rnd.Int())
		mdpass := parameter["password"].(string) + salt
		parameter["password"] = utils.Md5(mdpass)
		parameter["salt"] = salt
	}
	parameter["createtime"] = time.Now().Unix()
	if f_id == 0 {
		addId, err := DB().Table("admin_user").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			//添加角色-多个
			appRoleAccess(roleid, addId)
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("admin_user").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			//添加角色-多个
			if roleid != nil {
				appRoleAccess(roleid, f_id)
			}
			results.Success(context, "更新成功！", res, nil)
		}
	}
}

// 添加授权
func appRoleAccess(roleids []interface{}, uid interface{}) {
	//批量提交
	DB().Table("admin_auth_role_access").Where("uid", uid).Delete()
	save_arr := []map[string]interface{}{}
	for _, val := range roleids {
		marr := map[string]interface{}{"uid": uid, "role_id": val}
		save_arr = append(save_arr, marr)
	}
	DB().Table("admin_auth_role_access").Data(save_arr).Insert()
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
	res2, err := DB().Table("admin_user").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

// 更新头像-返回用户新信息
func UpAvatar(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	res2, err := DB().Table("admin_user").Where("id", user.ID).Data(map[string]interface{}{"avatar": parameter["url"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		useinfo, err := DB().Table("admin_user").Fields("id,username,name,avatar,status").Where("id", user.ID).First()
		if err != nil {
			results.Failed(context, "查找用户数据！", err)
			return
		}
		token := context.GetHeader("Authorization")
		roles_ids, err := DB().Table("admin_auth_role_access").Where("uid", user.ID).Pluck("role_id")
		roles, err := DB().Table("admin_auth_role").WhereIn("id", roles_ids.([]interface{})).Fields("id,name").Get()
		//获取用户权限
		results.Success(context, "更新成功！", map[string]interface{}{
			"userId":   useinfo["id"],
			"username": useinfo["username"],
			"token":    token,
			"realName": useinfo["name"],
			"avatar":   useinfo["avatar"],
			"desc":     "",
			"roles":    roles,
		}, res2)
	}
}

// 修改密码
func ChangePwd(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	userdata, _ := DB().Table("admin_user").Where("id", user.ID).Fields("id,password,salt").First()
	pass := utils.Md5(parameter["passwordOld"].(string) + userdata["salt"].(string))
	n_pass := utils.Md5(parameter["passwordNew"].(string) + userdata["salt"].(string))
	if pass != userdata["password"].(string) {
		results.Failed(context, "原来密码不正确！", nil)
	} else if n_pass == userdata["password"].(string) {
		results.Failed(context, "新密码和原来密码一样", nil)
	} else {
		newpass := utils.Md5(parameter["passwordNew"].(string) + userdata["salt"].(string))
		res2, err := DB().Table("admin_user").Where("id", user.ID).Data(map[string]interface{}{"password": newpass}).Update()
		if err != nil {
			results.Failed(context, "修改失败！", err)
		} else {
			results.Success(context, "修改成功", res2, nil)
		}
	}

}

// 检查账号是否存在
func IsAccountExist(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["id"] != nil {
		res1, err := DB().Table("admin_user").Where("id", "!=", parameter["id"]).Where("username", parameter["account"]).Value("id")
		if err != nil {
			results.Failed(context, "验证失败", err)
		} else if res1 != nil {
			results.Failed(context, "账号已存在", err)
		} else {
			results.Success(context, "验证通过", res1, nil)
		}
	} else {
		res2, err := DB().Table("admin_user").Where("username", parameter["account"]).Value("id")
		if err != nil {
			results.Failed(context, "验证失败", err)
		} else if res2 != nil {
			results.Failed(context, "账号已存在", err)
		} else {
			results.Success(context, "验证通过", res2, nil)
		}
	}
}

// 删除
func Del(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("admin_user").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		//并删除角色权限
		DB().Table("admin_auth_role_access").WhereIn("uid", ids.([]interface{})).Delete()
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 获取登录日志
func GetLoginLogList(context *gin.Context) {
	//当前用户
	userID := context.DefaultQuery("uid", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	list, _ := DB().Table("login_logs").Where("uid", userID).Where("type", 1).Limit(pageSize).Page(pageNo).Order("createtime desc").Get()
	var totalCount int64
	totalCount, _ = DB().Table("login_logs").Where("uid", userID).Where("type", 1).Count()
	results.Success(context, "获取登录日志", map[string]interface{}{
		"page":     pageNo,
		"pageSize": pageSize,
		"total":    totalCount,
		"items":    list}, nil)
}

// 获取账号信息
func GetAccount(context *gin.Context) {
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	data, _ := DB().Table("admin_user").Where("id", user.ID).First()
	results.Success(context, "获取账号信息", data, nil)
}
