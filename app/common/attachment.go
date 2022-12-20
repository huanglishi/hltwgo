package common

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取分类
func GetCatelist(context *gin.Context) {
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	filetype := context.DefaultQuery("filetype", "1")
	list, _ := DB().Table("attachment_cate").Where("accountID", user.Accountid).Where("filetype", filetype).Order("id asc").Get()
	results.Success(context, "获取列表", list, nil)
}

// 获取附件列表-分页
func GetImgList(context *gin.Context) {
	cid := context.DefaultQuery("cid", "0")
	rooting := context.DefaultQuery("rooting", "")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := DB().Table("attachment").Fields("id,title,url,cover_url,mimetype")
	MDBC := DB().Table("attachment")
	if cid != "0" {
		MDB.Where("cid", cid)
		MDBC.Where("cid", cid)
	} else {
		//当前用户
		getuser, _ := context.Get("user")
		user := getuser.(*utils.UserClaims)
		MDB.Where("accountID", user.Accountid)
		MDBC.Where("accountID", user.Accountid)
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			if val["cover_url"] != "" {
				val["imgs"] = rooting + val["cover_url"].(string)
			} else {
				val["imgs"] = rooting + val["url"].(string)
			}
		}
		var totalCount int64
		totalCount, _ = MDBC.Count()
		results.Success(context, "获取全部列表", map[string]interface{}{
			"total":    totalCount,
			"items":    list,
			"pageSize": pageSize,
			"page":     page,
		}, nil)
	}
}

// 添加分类
func CateAdd(context *gin.Context) {
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
		parameter["usetype"] = 1
		parameter["accountID"] = user.Accountid
		addId, err := DB().Table("attachment_cate").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("attachment_cate").
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

// 更新附件信息
func Upfile(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	upres, err := DB().Table("attachment").Where("id", parameter["id"]).Data(parameter).Update()
	if err != nil {
		results.Failed(context, "更新失败", err)
	} else {
		results.Success(context, "更新成功！", upres, nil)
	}
	context.Abort()
	return
}

// 删除分组
func DelCate(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("attachment_cate").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 删除附件及记录
func Delfile(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	attachment, _ := DB().Table("attachment").Where("id", parameter["id"]).First()
	if attachment != nil {
		//删除本地文件
		//判断文件是否存在
		_, err := os.Lstat(attachment["url"].(string))
		if os.IsNotExist(err) {
			os.Remove(attachment["url"].(string))
		}
		res2, err := DB().Table("attachment").Where("id", parameter["id"]).Delete()
		if err != nil {
			results.Failed(context, "删除附件数据失败", err)
		} else {
			results.Success(context, "删除成功！", res2, attachment)
		}
	} else {
		results.Failed(context, "删除附件失败", "文件找不到")
	}
	context.Abort()
	return
}
