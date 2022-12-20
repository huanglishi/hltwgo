package resource

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// 获取文件分类
func Getlist(context *gin.Context) {
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	menuList, _ := DB().Table("client_attachment").Where("cuid", user.ClientID).Where("type", 1).Fields("id,pid,title,weigh,createtime").Order("weigh asc").Get()
	for _, val := range menuList {
		numcount, _ := DB().Table("client_attachment").Where("pid", val["id"]).Where("type", 0).Count()
		val["num"] = numcount
	}
	// menuList = GetMenuChildrenArray(menuList, 0)
	menuList = GetTreeArray_x(menuList, 0, "")
	results.Success(context, "获取文件分类", menuList, nil)
}

// 获取文件分类-父级
func GetParentList(context *gin.Context) {
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	keyword := context.DefaultQuery("keyword", "")
	MDB := DB().Table("client_attachment")
	if keyword != "" {
		MDB = MDB.Where("title", "like", "%"+keyword+"%")
	}
	list, _ := MDB.Fields("id,pid,title").Where("type", 1).Where("cuid", user.ClientID).Order("weigh asc").Get()
	rulenum := GetTreeArray_x(list, 0, "")
	list_text := GetTreeList_txt(rulenum, "title")
	results.Success(context, "父级分组数据！", list_text, nil)
}

// 删除-文件及文件夹-子文件
func DelFileAndImg(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	//获取要删除的数据
	file_ids := GetAllChilIds(ids.([]interface{})) //批量获取子节点id
	file_ids_new := mergeArr(ids.([]interface{}), file_ids)
	res2, err := DB().Table("client_attachment").WhereIn("id", file_ids_new).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
}

// 1批量获取子节点id
func GetAllChilIds(ids []interface{}) []interface{} {
	var allsubids []interface{}
	for _, id := range ids {
		sub_ids := GetAllChilId(id)
		allsubids = append(allsubids, sub_ids...)
	}
	return allsubids
}

// 2获取所有子级ID
func GetAllChilId(id interface{}) []interface{} {
	var subids []interface{}
	sub_ids, _ := DB().Table("client_attachment").Where("pid", id).Pluck("id")
	if len(sub_ids.([]interface{})) > 0 {
		for _, sid := range sub_ids.([]interface{}) {
			subids = append(subids, sid)
			subids = append(subids, GetAllChilId(sid)...)
		}
	}
	return subids
}

// 　合并数组
func mergeArr(a, b []interface{}) []interface{} {
	var arr []interface{}
	for _, i := range a {
		arr = append(arr, i)
	}
	for _, j := range b {
		arr = append(arr, j)
	}
	return arr
}
