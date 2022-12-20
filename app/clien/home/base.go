package home

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

// 1 获取文章
func GetArticleList(context *gin.Context) {
	list, err := DB().Table("admin_article").Where("type", "th").Fields("id,title,star").Get()
	if err != nil {
		results.Failed(context, "获取文章失败", err)
	} else {
		if list == nil {
			list = make([]gorose.Data, 0)
		}
		results.Success(context, "获取文章列表", list, nil)
	}
}

// 1 获取文章
func GetArticle(context *gin.Context) {
	list, err := DB().Table("admin_article").Where("type", "th").First()
	if err != nil {
		results.Failed(context, "获取文章失败", err)
	} else {
		results.Success(context, "获取文章", list, nil)
	}
}

// 1 添加文章的赞
func PushStar(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	star, _ := DB().Table("admin_article").Where("id", parameter["id"]).Value("star")
	res2, err := DB().Table("admin_article").Where("id", parameter["id"]).Data(map[string]interface{}{"star": GetInterfaceToInt(star) + 1}).Update()
	if err != nil {
		results.Failed(context, "点赞失败！", err)
	} else {
		results.Success(context, "点赞成功！", res2, nil)
	}
}

// 2 获取微站信息
func GetMicweb(context *gin.Context) {
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	list, err := DB().Table("client_micweb").Where("cuid", user.ClientID).Fields("id,title,status,approval_err").First()
	if err != nil {
		results.Failed(context, "获取微站信息失败", err)
	} else {
		results.Success(context, "获取微站信息", list, nil)
	}
}

// 3 编辑微站
func SaveMicweb(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_micweb").Where("id", parameter["id"]).Data(parameter).Update()
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

// 4 发布微站
func PublishMicweb(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_micweb").Where("id", parameter["id"]).Data(map[string]interface{}{"status": 1, "publishtime": time.Now().Unix()}).Update()
	if err != nil {
		results.Failed(context, "已提交发布失败！", err)
	} else {
		results.Success(context, "已提交发布，请等待审核！", res2, nil)
	}
}
