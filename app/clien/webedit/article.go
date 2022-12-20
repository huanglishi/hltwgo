package webedit

import (
	"huling/utils/results"
	utils "huling/utils/tool"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取文章列表
func GetArticleList(context *gin.Context) {
	title := context.DefaultQuery("title", "")
	cid := context.DefaultQuery("cid", "")
	getall := context.DefaultQuery("getall", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	whereMap := DB().Table("client_article_manage").Where("cuid", user.ClientID).Where("status", 0)
	whereMap2 := DB().Table("client_article_manage").Where("cuid", user.ClientID).Where("status", 0)
	if title != "" {
		whereMap.Where("title", "like", "%"+title+"%")
		whereMap2.Where("title", "like", "%"+title+"%")
	}
	if cid != "" {
		article_ids, _ := DB().Table("client_article_cid").Where("cid", cid).Pluck("article_id")
		whereMap.WhereIn("id", article_ids.([]interface{}))
		whereMap2.WhereIn("id", article_ids.([]interface{}))
	}
	if getall != "0" {
		whereMap.Limit(pageSize).Page(pageNo)
	}
	list, err := whereMap.Fields("id,type,title,des,link,author,image,releasetime").Order("top desc , weigh desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			//分组
			cids, _ := DB().Table("client_article_cid").Where("article_id", val["id"]).Pluck("cid")
			catename, _ := DB().Table("client_article_cate").WhereIn("id", cids.([]interface{})).Pluck("name")
			val["catename"] = catename
		}
		var totalCount int64
		totalCount, _ = whereMap2.Count()
		results.Success(context, "获取文章列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"getall":   getall,
			"items":    list,
		}, nil)
	}
}

// 获取文章详情内容
func GetArticle(context *gin.Context) {
	id := context.DefaultQuery("id", "")
	data, err := DB().Table("client_article_manage").Where("id", id).Fields("id,type,title,des,author,image,releasetime,content").First()
	if err != nil {
		results.Failed(context, "获取文章详情内容失败", err)
	} else {
		results.Success(context, "获取文章详情内容", data, nil)
	}
}

// 获取文章分类
func GetArticleCate(context *gin.Context) {
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	data, err := DB().Table("client_article_cate").Where("cuid", user.ClientID).Fields("id,pid,name,weigh").Order("weigh asc , id asc").Get()
	if err != nil {
		results.Failed(context, "获取文章分类失败", err)
	} else {
		results.Success(context, "获取文章分类", data, nil)
	}
}
