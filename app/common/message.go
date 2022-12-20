package common

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取消息
func GetMessagelist(context *gin.Context) {
	usertype := context.DefaultQuery("usertype", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	list, err := DB().Table("common_message").Fields("id,type,title,path,content,isread,createtime").
		WhereIn("usertype", []interface{}{0, usertype}).Where("touid", user.ID).
		// Limit(pageSize).Page(pageNo).
		Order("id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		var totalCount int64
		totalCount, _ = DB().Table("common_message").WhereIn("usertype", []interface{}{0, usertype}).Where("touid", user.ID).Count()
		results.Success(context, "获取全部列表", map[string]interface{}{
			"total":    totalCount,
			"items":    list,
			"pageSize": pageSize,
			"page":     pageNo,
		}, nil)
	}
}

// 更新为已读
func SetRead(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("common_message").WhereIn("id", ids_arr).Data(map[string]interface{}{"isread": 1}).Update()
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
