package user

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

/*
1.获取用户信息
2.获取用户的
*/
func GetInfo(context *gin.Context) {
	getuser, ok := context.Get("user") //取值 实现了跨中间件取值
	if !ok {
		results.Failed(context, "用户id不存在！", ok)
		return
	}
	user := getuser.(*utils.UserClaims)
	userdata, err := DB().Table("client_user").Fields("id,accountID,username,name,avatar,status,password,salt,remark").Where("id", user.ID).First()
	if err != nil {
		results.Failed(context, "查找用户数据！", err)
		return
	}
	token := context.GetHeader("Authorization")
	roles_ids, err := DB().Table("client_auth_role_access").Where("uid", user.ID).Pluck("role_id")
	roles, err := DB().Table("client_auth_rule").WhereIn("id", roles_ids.([]interface{})).Fields("id,name").Get()

	mwurl, _ := DB().Table("merchant_config").Where("keyname", "mwurl").Value("keyvalue")
	tplpreviewurl, _ := DB().Table("merchant_config").Where("keyname", "tplpreviewurl").Value("keyvalue")
	rooturl, _ := DB().Table("merchant_config").Where("keyname", "rooturl").Value("keyvalue")
	if userdata["avatar"] == nil {
		userdata["avatar"] = "/resource/img/avatar.png"
	} else if !strings.Contains(userdata["avatar"].(string), "http") && rooturl != nil {
		userdata["avatar"] = rooturl.(string) + userdata["avatar"].(string)
	}
	//获取是否是管理账号
	is_admin, _ := DB().Table("merchant_user").Where("id", user.Accountid).Value("is_admin")
	//获取用户权限
	results.Success(context, "获取用户信息", map[string]interface{}{
		"userId":        userdata["id"],
		"accountID":     userdata["accountID"],
		"username":      userdata["username"],
		"password":      userdata["password"],
		"salt":          userdata["salt"],
		"token":         token,
		"realName":      userdata["name"],
		"avatar":        userdata["avatar"],
		"desc":          userdata["remark"],
		"is_admin":      is_admin,
		"homePath":      "/home",
		"mwurl":         mwurl,         //微站二维码
		"tplpreviewurl": tplpreviewurl, //轻站模板预览地址
		"rooturl":       rooturl,       //图片
		"roles":         roles,
	}, nil)
}

// 更新数据
func Updata(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("client_user").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
	if err != nil {
		context.JSON(200, gin.H{
			"code": 1,
			"msg":  "更新失败！",
			"data": err,
			"time": time.Now().Unix(),
		})
	} else {
		msg := "更新成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		context.JSON(200, gin.H{
			"code": 0,
			"msg":  msg,
			"data": res2,
			"time": time.Now().Unix(),
		})
	}

}

// 添加用户
func AddParam(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	parameter["createtime"] = time.Now().Unix()
	if parameter["valid_time"] != nil {
		valid_time, _ := time.Parse("01/02/2006", parameter["valid_time"].(string))
		parameter["valid_time"] = valid_time
	} else {
		parameter["valid_time"] = 0
	}
	f_id := parameter["id"].(float64)
	groupids := parameter["groupids"].([]interface{})
	delete(parameter, "groupids")
	if f_id == 0 {
		rnd := rand.New(rand.NewSource(6))
		salt := strconv.Itoa(rnd.Int())
		mdpass := parameter["password"].(string) + salt
		parameter["password"] = utils.Md5(mdpass)
		parameter["salt"] = salt
		addId, err := DB().Table("client_user").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			group_arr := []map[string]interface{}{}
			for _, val := range groupids {
				marr := map[string]interface{}{"uid": addId, "group_id": val}
				group_arr = append(group_arr, marr)
			}
			DB().Table("admin_auth_group_access").Where("uid", f_id).Delete()
			DB().Table("admin_auth_group_access").Data(group_arr).Insert()
			results.Success(context, "添加成功！", nil, nil)
		}
	} else {
		if parameter["password"] != nil {
			rnd := rand.New(rand.NewSource(6))
			salt := strconv.Itoa(rnd.Int())
			mdpass := parameter["password"].(string) + salt
			parameter["password"] = utils.Md5(mdpass)
			parameter["salt"] = salt
		}
		res, err := DB().Table("client_user").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			group_arr := []map[string]interface{}{}
			for _, val := range groupids {
				marr := map[string]interface{}{"uid": f_id, "group_id": val}
				group_arr = append(group_arr, marr)
			}
			DB().Table("admin_auth_group_access").Where("uid", f_id).Delete()
			DB().Table("admin_auth_group_access").Data(group_arr).Insert()
			results.Success(context, "更新成功！", res, group_arr)
		}
	}
	context.Abort()
	return
}

// 获取用列表
func QueryParam(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	group_ids, _ := DB().Table("admin_auth_group_access").Where("uid", user.ID).Pluck("group_id")                      //获取用户分组
	data_access_ids, _ := DB().Table("admin_auth_group").WhereIn("id", group_ids.([]interface{})).Pluck("data_access") //获取用户数据权限
	_pageNo := context.DefaultQuery("pageNo", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(_pageNo)
	pageSize, _ := strconv.Atoi(_pageSize)
	//权限判断
	var res []gorose.Data
	var err error
	var access_all int64 = 2
	var access_myandchrid int64 = 1
	if utils.In_array(access_all, data_access_ids) { //全部
		_res, _err := DB().Table("client_user").Fields("id, name,username,telephone,lastLoginIp,lastLoginTime,status,valid_time").Page(pageNo).Limit(pageSize).Order("id desc").Get()
		res = _res
		err = _err
	} else if utils.In_array(access_myandchrid, data_access_ids) { //自己和子集
		_res, _err := DB().Table("client_user").Fields("id, name,username,telephone,lastLoginIp,lastLoginTime,status,valid_time").Page(pageNo).Limit(pageSize).Order("id desc").Get()
		res = _res
		err = _err
	} else { //仅自己
		_res, _err := DB().Table("client_user").Where("uid", user.ID).Fields("id, name,username,telephone,lastLoginIp,lastLoginTime,status,valid_time").Page(pageNo).Limit(pageSize).Order("id desc").Get()
		res = _res
		err = _err
	}
	if err != nil {
		results.Failed(context, "查找失败", err)
	} else {
		// 统计数据
		var totalCount int64
		totalCount, _ = DB().Table("client_user").Count()
		_pageSize := int64(pageSize)
		totalPage := totalCount / _pageSize
		for _, val := range res {
			groupids, _ := DB().Table("admin_auth_group_access").Where("uid", val["id"]).Pluck("group_id")
			val["groupids"] = groupids
			groupname, _ := DB().Table("admin_auth_group").WhereIn("id", groupids.([]interface{})).Fields("id,name").Get()
			val["groupname"] = groupname
		}
		results.Success(context, "查找成功！", map[string]interface{}{
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalCount": totalCount,
			"totalPage":  totalPage,
			"data":       res}, nil)
	}
}

// 删除操作
func DelParam(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("client_user").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		//并删除角色权限
		DB().Table("client_auth_role_access").WhereIn("uid", ids.([]interface{})).Delete()
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return

}

// 修改密码
func ChangePwd(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//账号信息
	getuser, _ := context.Get("user") //当前用户
	user := getuser.(*utils.UserClaims)
	userdata, err := DB().Table("client_user").Where("id", user.ID).Fields("password,salt").First()
	if err != nil {
		results.Failed(context, "修改密码失败", err)
	} else {
		pass := utils.Md5(parameter["passwordOld"].(string) + userdata["salt"].(string))
		if userdata["password"] != pass {
			results.Failed(context, "原来密码输入错误！", err)
		} else {
			newpass := utils.Md5(parameter["passwordNew"].(string) + userdata["salt"].(string))
			res, err := DB().Table("client_user").
				Data(map[string]interface{}{"password": newpass}).
				Where("id", user.ID).
				Update()
			if err != nil {
				results.Failed(context, "修改密码失败", err)
			} else {
				results.Success(context, "修改密码成功！", res, nil)
			}
		}
	}
	context.Abort()
	return
}

// 获取编辑用户数据
func GetUserData(context *gin.Context) {
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	userdata, err := DB().Table("client_user").Fields("id,accountID,username,nickname,avatar,status,password,salt,tel,email,city,remark").Where("id", user.ID).First()
	if err != nil {
		results.Failed(context, " 获取用户数据失败", err)
	} else {
		results.Success(context, " 获取编辑用户数据", userdata, nil)
	}
}

// 更新用户数据
func UpUserInfo(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	res, err := DB().Table("client_user").
		Data(parameter).
		Where("id", user.ID).
		Update()
	if err != nil {
		results.Failed(context, "更新失败", err)
	} else {
		results.Success(context, " 更新用户数据成功", res, nil)
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
	res, err := DB().Table("client_user").Where("id", user.ID).Data(map[string]interface{}{"avatar": parameter["url"]}).Update()
	if err != nil {
		results.Failed(context, "更新头像失败！", err)
	} else {
		results.Success(context, " 更新头像成功", res, nil)
	}
}
