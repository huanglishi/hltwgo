package form

import (
	"encoding/json"
	"fmt"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取表单字段
func GetFormField(context *gin.Context) {
	form_id := context.DefaultQuery("form_id", "0")
	list, _ := DB().Table("client_form_item").Where("form_id", form_id).Fields("id as dataIndex,name as title").Order("weigh asc").Get()
	fromtitle, _ := DB().Table("client_form").Where("id", form_id).Value("name")
	results.Success(context, "获取表单项数据", map[string]interface{}{"list": list, "title": fromtitle}, nil)
}

// 获取表单数据
func GetFormDataList(context *gin.Context) {
	keyword := context.DefaultQuery("keyword", "")
	username := context.DefaultQuery("username", "")
	form_id := context.DefaultQuery("form_id", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	MDB := DB().Table("client_form_record a").LeftJoin("client_form_value b on a.id = b.record_id")
	//填写字段查询
	if keyword != "" {
		MDB.Where("value", "like", "%"+keyword+"%")
	}
	//用户查询
	if username != "" {
		userids, _ := DB().Table("client_member").Where("name", "like", "%"+username+"%").Pluck("id")
		MDB.WhereIn("member_id", userids.([]interface{}))
	}
	list, err := MDB.Where("form_id", form_id).Where("cuid", user.ClientID).Limit(pageSize).Page(pageNo).Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		//字段赋值
		fielddata, _ := DB().Table("client_member").Where("form_id", form_id).Pluck("id")
		for _, val := range list {
			for _, item := range fielddata.([]interface{}) {
				//字段的key
				keys := fmt.Sprintf("%v", item)
				//字段值
				valuedata, _ := DB().Table("client_form_value").Where("form_item_id", item).Where("record_id", val["record_id"]).Value("value")
				val[keys] = valuedata
			}
		}
		var totalCount int64
		totalCount, _ = MDB.Where("form_id", form_id).Where("cuid", user.ClientID).Count()
		results.Success(context, "获取文章列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 删除-填写记录
func DelData(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("client_form_record").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
}
