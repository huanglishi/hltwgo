package webedit

// 表单
import (
	"huling/utils/results"
	utils "huling/utils/tool"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获列表
func GetFormList(context *gin.Context) {
	title := context.DefaultQuery("title", "")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	whereMap := DB().Table("client_form").Where("cuid", user.ClientID)
	whereMap2 := DB().Table("client_form").Where("cuid", user.ClientID)
	if title != "" {
		whereMap.Where("name", "like", "%"+title+"%")
		whereMap2.Where("name", "like", "%"+title+"%")
	}
	list, err := whereMap.Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		var totalCount int64
		totalCount, _ = whereMap2.Count()
		results.Success(context, "获取表单列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 获取表单字段
func GetFormField(context *gin.Context) {
	id := context.DefaultQuery("id", "")
	data, err := DB().Table("client_form_item").Where("form_id", id).Get()
	if err != nil {
		results.Failed(context, "获取表单字段失败", err)
	} else {
		results.Success(context, "获取表单字段", data, nil)
	}
}

// 获取表单规则
func GetFormRule(context *gin.Context) {
	form_id := context.DefaultQuery("form_id", "")
	data, err := DB().Table("client_form_rule").Where("form_id", form_id).Get()
	for _, val := range data {
		if val["show_item_ids"] != nil {
			val["show_item_ids"] = StingToJSON(val["show_item_ids"])
		}
		if val["show_item_text"] != nil {
			val["show_item_text"] = StingToJSON(val["show_item_text"])
		}
		if val["selectval"] != nil {
			val["selectval"] = StingToJSON(val["selectval"])
		}
	}
	if err != nil {
		results.Failed(context, "获取表单规则失败", err)
	} else {
		results.Success(context, "获取表单规则", data, nil)
	}
}
