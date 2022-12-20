package platformaccount

import (
	"encoding/json"
	"fmt"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取数据列表-分页先用Limit在page分则无小-type=1是平台
func Getlist(context *gin.Context) {
	cid := context.DefaultQuery("cid", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := DB().Table("merchant_user").Fields("id,is_admin,status,name,username,validtime,avatar,mobile,email,remark,groupid,remark,createtime,sub_account")
	if cid != "0" {
		MDB.Where("groupid", cid)
	}
	list, err := MDB.Where("type", 1).Limit(pageSize).Page(pageNo).Order("id asc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			depname, _ := DB().Table("admin_platform_group").Where("id", val["groupid"]).Value("name")
			val["depname"] = depname
			var zro int64 = 0
			if val["validtime"] != nil && val["validtime"].(int64) != zro {
				val["validtime"] = utils.DateToStr(val["validtime"].(int64), "ymd")
			} else {
				val["validtime"] = ""
			}
		}
		var totalCount int64
		totalCount, _ = DB().Table("merchant_user").Count()
		results.Success(context, "获取全部列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 获取父级数据
func GetParentList(context *gin.Context) {
	menuList, _ := DB().Table("merchant_user").Fields("id,pid,name").Order("id asc").Get()
	menuList = GetMenuChildrenArray(menuList, 0)
	results.Success(context, "部门父级数据！", menuList, nil)
}

// 添加用户
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
	if parameter["validtime"] != nil {
		parameter["validtime"] = utils.StrToTime(parameter["validtime"].(string))
	} else {
		parameter["validtime"] = 0
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
		parameter["avatar"] = "resource/staticfile/avatar.png"
		addId, err := DB().Table("merchant_user").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 { //添加账号id
				DB().Table("merchant_user").
					Data(map[string]interface{}{"accountID": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("merchant_user").
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

// 保存套餐设置
func SaveSetting(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("merchant_user").Where("id", parameter["id"]).Data(map[string]interface{}{"sub_account": parameter["sub_account"], "status": parameter["status"]}).Update()
	if err != nil {
		results.Failed(context, "设置提交失败！", err)
	} else {
		//1添加套餐数据
		getmenus := appAllpackge(parameter["pdatas"].([]interface{}), parameter["id"])
		if getmenus != nil {
			rules := GetRulesID(getmenus) //获取子菜单包含的父级ID
			rudata := rules.([]interface{})
			var rulesStr []string
			for _, v := range rudata {
				str := fmt.Sprintf("%v", v) //interface{}强转string
				rulesStr = append(rulesStr, str)
			}
			parameter["rules"] = strings.Join(rulesStr, ",")
			parameter["menu"] = JSONMarshalToString(getmenus)
		}
		//2 添加一个账号对应授权
		roleid, _ := DB().Table("merchant_auth_role_access").Where("uid", parameter["id"]).Value("role_id")
		if roleid != nil { //更新
			DB().Table("merchant_auth_role").Where("id", roleid).Data(map[string]interface{}{"rules": parameter["rules"], "menu": parameter["menu"]}).Update()
		} else { //添加
			//当前用户
			getuser, _ := context.Get("user")
			user := getuser.(*utils.UserClaims)
			addId, _ := DB().Table("merchant_auth_role").
				Data(map[string]interface{}{"uid": user.ID, "name": "总管理组", "rules": parameter["rules"], "menu": parameter["menu"],
					"remark": "账号的总管理员", "createtime": time.Now().Unix(), "weigh": 1}).InsertGetId()
			//添加权限
			DB().Table("merchant_auth_role_access").Data(map[string]interface{}{"uid": parameter["id"], "role_id": addId}).Insert()
		}
		results.Success(context, "设置提交成功", res2, nil)
	}
}

// 批量添加套餐
func appAllpackge(plist []interface{}, userid interface{}) []interface{} {
	//批量提交
	DB().Table("merchant_user_package").Where("user_id", userid).Delete()
	save_arr := []map[string]interface{}{}
	var pids = []interface{}{}
	for _, val := range plist {
		items := val.(map[string]interface{})
		marr := map[string]interface{}{"user_id": userid, "accountID": userid, "pid": items["pid"], "checks": items["checks"],
			"number": items["number"], "type": items["type"]}
		save_arr = append(save_arr, marr)
		pids = append(pids, items["pid"])
	}
	DB().Table("merchant_user_package").Data(save_arr).Insert()
	//获取菜单id
	menu_str, _ := DB().Table("merchant_packagedesign").WhereIn("id", pids).Pluck("menu_str")
	getmenus := ArrayMerge(menu_str.([]interface{}))
	return getmenus
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
	res2, err := DB().Table("merchant_user").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

// 检查账号是否存在
func IsAccountExist(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["id"] != nil {
		res1, err := DB().Table("merchant_user").Where("id", "!=", parameter["id"]).Where("username", parameter["account"]).Value("id")
		if err != nil {
			results.Failed(context, "验证失败", err)
		} else if res1 != nil {
			results.Failed(context, "账号已存在", err)
		} else {
			results.Success(context, "验证通过", res1, nil)
		}
	} else {
		res2, err := DB().Table("merchant_user").Where("username", parameter["account"]).Value("id")
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
	res2, err := DB().Table("merchant_user").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
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
	list, _ := DB().Table("admin_login_logs").Where("uid", userID).Where("type", 1).Limit(pageSize).Page(pageNo).Order("createtime desc").Get()
	var totalCount int64
	totalCount, _ = DB().Table("admin_login_logs").Where("uid", userID).Where("type", 1).Count()
	results.Success(context, "获取登录日志", map[string]interface{}{
		"page":     pageNo,
		"pageSize": pageSize,
		"total":    totalCount,
		"items":    list}, nil)
}

// 获取账号信息
func GetAccount(context *gin.Context) {
	//当前用户
	userID := context.DefaultQuery("id", "0")
	data, _ := DB().Table("merchant_user").Where("id", userID).First()
	results.Success(context, "获取账号信息", data, nil)
}
