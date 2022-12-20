package file

import (
	"huling/utils/results"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取分组分类
func GetCateList(context *gin.Context) {
	MDB := DB().Table("client_picture_cate")
	menuList, _ := MDB.Fields("id,name,type").Order("weigh asc").Get()
	results.Success(context, "获取数据列表", menuList, nil)
}

// 获取图片
func GetPicture(context *gin.Context) {
	typeid := context.DefaultQuery("type", "0")
	cid := context.DefaultQuery("cid", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "15")
	searchword := context.DefaultQuery("searchword", "")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := DB().Table("client_picture").Where("status", 0).Where("type", typeid).Where("cid", cid)
	CMDB := DB().Table("client_picture").Where("status", 0).Where("type", typeid).Where("cid", cid)
	if searchword != "" {
		MDB = MDB.Where("title", "like", "%"+searchword+"%")
		CMDB = CMDB.Where("title", "like", "%"+searchword+"%")
	}

	list, err := MDB.Fields("id,cid,type,title,url").Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		var totalCount int64
		totalCount, _ = CMDB.Count()
		results.Success(context, "获取列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}
