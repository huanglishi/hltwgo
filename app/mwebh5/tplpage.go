package mwebh5

import (
	"huling/utils/results"

	"github.com/gin-gonic/gin"
)

// 2 获取单页模板
func GetTplPage(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	if id == "0" {
		results.Failed(context, "清传参数id", nil)
	}
	data, err := DB().Table("client_micweb_tplpage_page").Where("id", id).First()
	if err != nil {
		results.Failed(context, "获取单页模板失败", err)
	} else {
		results.Success(context, "获取单页模板", data, nil)
	}
}
