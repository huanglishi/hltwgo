package resource

import (
	"encoding/json"
	"fmt"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

// 获取图片
func GetPicture(context *gin.Context) {
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "15")
	searchword := context.DefaultQuery("searchword", "")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	MDB := DB().Table("client_attachment").Where("cuid", user.ClientID).WhereIn("type", []interface{}{0, 2})
	CMDB := DB().Table("client_attachment").Where("cuid", user.ClientID).WhereIn("type", []interface{}{0, 2})
	if searchword != "" {
		MDB = MDB.Where("title", "like", "%"+searchword+"%")
		CMDB = CMDB.Where("title", "like", "%"+searchword+"%")
	}
	list, err := MDB.Fields("id,type,title,url,cover_url").Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		var totalCount int64
		totalCount, _ = CMDB.Count()
		//全部图片
		var allnumber int64
		var allsize interface{}
		var fileSize interface{}
		allnumber, _ = DB().Table("client_attachment").Where("cuid", user.ClientID).Where("type", 0).Count()
		allsize, _ = DB().Table("client_attachment").Where("cuid", user.ClientID).Where("type", 0).Sum("filesize")
		fileSize, _ = DB().Table("client_user_config").Where("cuid", user.ClientID).Value("fileSize")
		if list == nil {
			list = make([]gorose.Data, 0)
		}
		results.Success(context, "获取列表", map[string]interface{}{
			"page":      pageNo,
			"pageSize":  pageSize,
			"total":     totalCount,
			"allnumber": allnumber,
			"allsize":   allsize,
			"items":     list,
			"fileSize":  fileSize,
		}, nil)
	}
}

// 更新文件名称
func UpFile(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_attachment").Where("id", parameter["id"]).Data(map[string]interface{}{"title": parameter["title"]}).Update()
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

// 删除
func DelFile(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	file_list, _ := DB().Table("client_attachment").WhereIn("id", ids.([]interface{})).Pluck("url")
	res2, err := DB().Table("client_attachment").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		del_file(file_list.([]interface{}))
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 删除本地附件
func del_file(file_list []interface{}) {
	for _, val := range file_list {
		dir := fmt.Sprintf("./%s", val)
		os.Remove(dir)
	}
}

// 测试上传图片
func DelTest(context *gin.Context) {
	dir := "./resource/uploads/1222.png"
	err := os.Remove(dir)
	if err != nil {
		// 删除失败
		results.Failed(context, "删除失败", err)
	} else {
		// 删除成功
		results.Success(context, "附件删除成功！", 1, nil)
	}
}
