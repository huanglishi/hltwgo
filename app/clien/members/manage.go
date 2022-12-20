package members

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取列表
func GetList(context *gin.Context) {
	uname := context.DefaultQuery("name", "")
	status := context.DefaultQuery("status", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	whereMap := DB().Table("client_member").Where("cuid", user.ClientID)
	whereMap2 := DB().Table("client_member").Where("cuid", user.ClientID)
	if status != "0" {
		whereMap.Where("status", status)
		whereMap2.Where("status", status)
	}
	if uname != "" {
		whereMap.Where("name", "like", "%"+uname+"%")
		whereMap2.Where("name", "like", "%"+uname+"%")
	}
	list, err := whereMap.Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			//分组
			groupname, _ := DB().Table("client_member_group").Where("id", val["cid"]).Pluck("name")
			if groupname != nil {
				val["groupname"] = groupname
			} else {
				val["groupname"] = "未分组"
			}
		}
		var totalCount int64
		totalCount, _ = whereMap2.Count()
		results.Success(context, "获取文章列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 添加
func SaveArticle(context *gin.Context) {
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
	if parameter["releasetime"] != nil {
		var LOC, _ = time.LoadLocation("Asia/Shanghai")
		tim, _ := time.ParseInLocation("2006-01-02 15:04:05", parameter["releasetime"].(string), LOC)
		parameter["releasetime"] = tim.Unix()
	}
	delete(parameter, "catename")
	delete(parameter, "pendingStatus")
	if f_id == 0 {
		parameter["cuid"] = user.ClientID
		parameter["createtime"] = time.Now().Unix()
		parameter["accountID"] = user.Accountid
		addId, err := DB().Table("client_member").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("client_member").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			results.Success(context, "更新成功！", res, user)
		}
	}
}

// 更新状态
func UpLock(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_member").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

// 删除
func Del(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("client_member").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
