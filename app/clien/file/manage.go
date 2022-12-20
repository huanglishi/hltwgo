package file

import (
	"encoding/json"
	"fmt"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取附件列表
func GetFiles(context *gin.Context) {
	searchword := context.DefaultQuery("searchword", "")
	filetype := context.DefaultQuery("filetype", "image")
	pid := context.DefaultQuery("pid", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	whereMap := DB().Table("client_attachment").Where("cuid", user.ClientID).Where("pid", pid)
	if searchword != "" {
		whereMap.Where("title", "like", "%"+searchword+"%")
	}
	if filetype == "video" {
		whereMap.WhereIn("type", []interface{}{1, 2})
	} else { //默认图片
		whereMap.WhereIn("type", []interface{}{0, 1})
	}
	list, err := whereMap.
		// Limit(pageSize).Page(pageNo).
		Fields("id,pid,name,title,type,url,filesize,mimetype,storage,cover_url").Order("type desc,weigh asc , id asc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		var totalCount int64
		//获取目录菜单
		allids := getAllParentIds(pid)
		allids = append(allids, pid)
		dirmenu, _ := DB().Table("client_attachment").WhereIn("id", allids).Fields("id,pid,title").Get()
		results.Success(context, "获取附件列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"dirmenu":  dirmenu,
			"allids":   allids,
			"items":    list,
		}, nil)
	}
}

// 工具
func getAllParentIds(id interface{}) []interface{} {
	var parent_ids []interface{}
	parent_id, _ := DB().Table("client_attachment").Where("id", id).Value("pid")
	if parent_id != nil {
		parent_ids = append(parent_ids, parent_id)
		parent_ids = append(parent_ids, getAllParentIds(parent_id)...)
	}
	return parent_ids
}

// 添加文件夹
func SaveFile(context *gin.Context) {
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
		parameter["accountID"] = user.Accountid
		parameter["cuid"] = user.ClientID
		getcount, _ := DB().Table("client_attachment").Where("cuid", user.ClientID).Where("pid", parameter["pid"]).Where("title", "like", fmt.Sprintf("%s%v%s", "%", parameter["title"], "%")).Count()
		parameter["title"] = fmt.Sprintf("%s%v", parameter["title"], GetInterfaceToInt(getcount)+1)
		addId, err := DB().Table("client_attachment").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			//更新排序
			DB().Table("client_attachment").Data(map[string]interface{}{"weigh": addId}).Where("id", addId).Update()
			getdata, _ := DB().Table("client_attachment").Where("id", addId).Fields("id,pid,name,title,type,url,filesize,mimetype,storage").First()
			results.Success(context, "添加成功！", getdata, nil)
		}
	} else {
		res, err := DB().Table("client_attachment").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			results.Success(context, "更新成功！", res, user)
		}
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

// 更新图片目录
func UpImgPid(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_attachment").Where("id", parameter["imgid"]).Data(map[string]interface{}{"pid": parameter["pid"]}).Update()
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

// 获取文章大编辑内容
func GetArticle(context *gin.Context) {
	id := context.DefaultQuery("id", "")
	data, err := DB().Table("client_article_manage").Where("id", id).First()
	if err != nil {
		results.Failed(context, "获取文章大内容字段失败", err)
	} else {
		cids, _ := DB().Table("client_article_cid").Where("article_id", id).Pluck("cid")
		data["cid"] = cids
		results.Success(context, "文章大内容字段", data, nil)
	}
}
