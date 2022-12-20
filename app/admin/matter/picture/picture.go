package picture

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取数据列表
func Getlist(context *gin.Context) {
	page := context.DefaultQuery("page", "1")
	typeid := context.DefaultQuery("type", "2")
	_pageSize := context.DefaultQuery("pageSize", "10")
	title := context.DefaultQuery("title", "")
	status := context.DefaultQuery("status", "0")
	orderid := context.DefaultQuery("orderid", "")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := DB().Table("client_picture")
	CMDB := DB().Table("client_picture")
	if status != "0" {
		MDB = MDB.Where("status", status)
		CMDB = CMDB.Where("status", status)
	}
	if orderid != "" {
		MDB = MDB.Where("orderid", orderid)
		CMDB = CMDB.Where("orderid", orderid)
	}
	if typeid != "2" {
		MDB = MDB.Where("type", typeid)
		CMDB = CMDB.Where("type", typeid)
	}
	if title != "" {
		MDB = MDB.Where("title", "like", "%"+title+"%")
		CMDB = CMDB.Where("title", "like", "%"+title+"%")
	}

	list, err := MDB.Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			catename, _ := DB().Table("client_picture_cate").Where("id", val["cid"]).Value("name")
			val["catename"] = catename
		}
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

// 添加
func Add(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	parameter["createtime"] = time.Now().Unix()
	res, err := DB().Table("client_picture").
		Data(parameter).
		Where("id", parameter["id"]).
		Update()
	if err != nil {
		results.Failed(context, "更新失败", err)
	} else {
		results.Success(context, "更新成功！", res, nil)
	}
}

// 更新状态
func UpLock(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("client_picture").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
func Del(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	file_list, _ := DB().Table("client_picture").WhereIn("id", ids.([]interface{})).Pluck("url")
	res2, err := DB().Table("client_picture").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		del_file(file_list.([]interface{}))
		results.Success(context, "删除成功1！", res2, file_list)
	}
}

// 删除本地附件
func del_file(file_list []interface{}) {
	for _, val := range file_list {
		dir := fmt.Sprintf("./%s", val)
		os.Remove(dir)
	}
}

// 上传图片
func UploadFile(context *gin.Context) {
	// 单个文件
	cid := context.DefaultPostForm("cid", "")
	typeId := context.DefaultPostForm("type", "")
	Id := context.DefaultPostForm("id", "0")
	file, err := context.FormFile("file")
	if err != nil {
		results.Failed(context, "获取数据失败，", err)
		return
	}
	nowTime := time.Now().Unix()      //当前时间
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	//判断文件是否已经传过
	fileContent, _ := file.Open()
	var byteContainer []byte
	byteContainer = make([]byte, 1000000)
	fileContent.Read(byteContainer)
	m_d5 := md5.New()
	m_d5.Write(byteContainer)
	sha1_str := hex.EncodeToString(m_d5.Sum(nil))
	//查找该用户是否传过
	attachment, _ := DB().Table("client_picture").Where("uid", user.ID).
		Where("sha1", sha1_str).Fields("id,name,title,url,filesize,mimetype,storage").First()
	if attachment != nil { //文件是否已经存在
		//更新到最前面
		var nid interface{}
		if Id != "0" {
			DB().Table("client_picture").Data(map[string]interface{}{"title": attachment["title"], "name": attachment["name"], "url": attachment["url"]}).Where("id", Id).Update()
			nid = Id
		} else {
			delete(attachment, "id")
			attachment["cid"] = cid
			attachment["type"] = typeId
			file_id, _ := DB().Table("client_picture").Data(attachment).InsertGetId()
			nid = file_id
			//更新排序
			DB().Table("client_picture").Data(map[string]interface{}{"weigh": file_id}).Where("id", file_id).Update()
		}
		results.Success(context, "文件已上传", map[string]interface{}{"id": nid, "title": attachment["title"], "url": attachment["url"]}, nil)
	} else {
		file_path := fmt.Sprintf("%s%s%s", "resource/uploads/", time.Now().Format("20060102"), "/")
		//如果没有filepath文件目录就创建一个
		if _, err := os.Stat(file_path); err != nil {
			if !os.IsExist(err) {
				os.MkdirAll(file_path, os.ModePerm)
			}
		}
		//上传到的路径
		filename_arr := strings.Split(file.Filename, ".")
		//重新名片-lunix系统不支持中文
		name_str := md5Str(fmt.Sprintf("%v%s", nowTime, filename_arr[0]))      //组装文件保存名字
		file_Filename := fmt.Sprintf("%s%s%s", name_str, ".", filename_arr[1]) //文件加.后缀
		path := file_path + file_Filename
		// 上传文件到指定的目录
		err = context.SaveUploadedFile(file, path)
		if err != nil { //上传失败
			context.JSON(200, gin.H{
				"uid":      sha1_str,
				"name":     file.Filename,
				"status":   "error",
				"response": "上传失败",
				"time":     nowTime,
			})
		} else { //上传成功
			//保存数据
			dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
			Insertdata := map[string]interface{}{
				"uid":        user.ID,
				"type":       typeId,
				"cid":        cid,
				"sha1":       sha1_str,
				"title":      filename_arr[0],
				"name":       file.Filename,
				"url":        path,
				"storage":    dir + strings.Replace(path, "/", "\\", -1),
				"createtime": nowTime,
				"filesize":   file.Size,
				"mimetype":   file.Header["Content-Type"][0],
			}
			//保存数据
			var nid interface{}
			if Id != "0" {
				DB().Table("client_picture").Data(Insertdata).Where("id", Id).Update()
				nid = Id
			} else {
				file_id, _ := DB().Table("client_picture").Data(Insertdata).InsertGetId()
				nid = file_id
				//更新排序
				DB().Table("client_picture").Data(map[string]interface{}{"weigh": file_id}).Where("id", file_id).Update()
			}
			//返回数据
			results.Success(context, "上传成功", map[string]interface{}{"id": nid, "title": Insertdata["title"], "url": Insertdata["url"]}, nil)
		}
	}
	context.Abort()
}

// md5加密
func md5Str(origin string) string {
	m := md5.New()
	m.Write([]byte(origin))
	return hex.EncodeToString(m.Sum(nil))
}
