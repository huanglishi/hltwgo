package webmain

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

// 1添加素材
func AddMaterial(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	parameter["accountID"] = user.Accountid
	parameter["createtime"] = time.Now().Unix()
	addId, err := DB().Table("merchant_micweb_material").Data(parameter).InsertGetId()
	if err != nil {
		results.Failed(context, "添加失败", err)
	} else {
		results.Success(context, "添加成功！", addId, nil)
	}
}

// 2获取模板数据
func GetMaterial(context *gin.Context) {
	// page := context.DefaultQuery("page", "1")
	// _pageSize := context.DefaultQuery("pageSize", "20")
	// pageNo, _ := strconv.Atoi(page)
	// pageSize, _ := strconv.Atoi(_pageSize)
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("merchant_micweb_material").Where("accountID", user.Accountid).OrWhere("type", 0).Fields("id,title,image,type,jsondata").
		// Limit(pageSize).Page(pageNo).
		Order("id asc").Get()
	results.Success(context, "获取模板数据", list, nil)
}

// 3使用模板数据
func UseMaterial(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	data, _ := DB().Table("merchant_micweb_material").Where("id", id).Fields("templateJson,component").First()
	results.Success(context, "使用模板数据", data, nil)
}

// 4 删除-素材
func DelMaterial(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("merchant_micweb_material").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
