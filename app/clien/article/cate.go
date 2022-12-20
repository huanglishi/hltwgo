package article

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取分类列表
func GetCateList(context *gin.Context) {
	title := context.DefaultQuery("title", "")
	status := context.DefaultQuery("status", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	wheremap := DB().Table("client_article_cate").Where("cuid", user.ClientID)
	wheremap2 := DB().Table("client_article_cate").Where("cuid", user.ClientID)
	if status != "0" {
		wheremap.Where("status", status)
		wheremap2.Where("status", status)
	}
	if title != "" {
		wheremap.Where("title", "like", "%"+title+"%")
		wheremap2.Where("title", "like", "%"+title+"%")
	}
	list, err := wheremap.
		// Limit(pageSize).Page(pageNo).
		Order("weigh desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			//上一级
			pidname, _ := DB().Table("client_article_cate").Where("id", val["pid"]).Value("name")
			if pidname != nil {
				val["pidname"] = pidname
			} else {
				val["pidname"] = "无"
			}
		}
		var totalCount int64
		totalCount, _ = wheremap2.Count()
		rulenum := GetTreeArray(list, 0, "")
		list_text := GetTreeList_txt(rulenum, "name")
		results.Success(context, "获取分类列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list_text,
		}, nil)
	}
}

// 添加分类
func SaveCate(context *gin.Context) {
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
	if f_id == 0 {
		parameter["createtime"] = time.Now().Unix()
		parameter["cuid"] = user.ClientID
		parameter["accountID"] = user.Accountid
		addId, err := DB().Table("client_article_cate").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 {
				DB().Table("client_article_cate").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("client_article_cate").
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

// 获取分组tree-打成二维数组
func GetFormCateList(context *gin.Context) {
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	keyword := context.DefaultQuery("keyword", "")
	MDB := DB().Table("client_article_cate")
	if keyword != "" {
		MDB = MDB.Where("name", "like", "%"+keyword+"%")
	}
	list, _ := MDB.Fields("id,pid,name").Where("cuid", user.ClientID).Order("weigh desc").Get()
	rulenum := GetTreeArray(list, 0, "")
	list_text := GetTreeList_txt(rulenum, "name")
	results.Success(context, "分组数据！", list_text, nil)
}

// 删除
func DelCate(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("client_article_cate").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 更新状态
func UpCateStatus(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_article_cate").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
