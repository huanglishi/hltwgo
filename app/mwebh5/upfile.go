package mwebh5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"huling/utils/results"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 上传图片
func UploadImage(context *gin.Context) {
	// 单个文件
	pid := context.DefaultPostForm("pid", "0")
	uid := context.DefaultPostForm("uid", "0")
	file, err := context.FormFile("file")
	if err != nil {
		results.Failed(context, "获取图片内容失败，", err)
		return
	}
	nowTime := time.Now().Unix() //当前时间
	userinfo, _ := DB().Table("client_member").Where("id", uid).Fields("cuid,accountID").First()
	//判断文件是否已经传过
	fileContent, _ := file.Open()
	var byteContainer []byte
	byteContainer = make([]byte, 1000000)
	fileContent.Read(byteContainer)
	m_d5 := md5.New()
	m_d5.Write(byteContainer)
	sha1_str := hex.EncodeToString(m_d5.Sum(nil))
	//查找该用户是否传过
	attachment, _ := DB().Table("client_attachment").Where("cuid", userinfo["cuid"]).
		Where("sha1", sha1_str).Fields("id,pid,name,title,type,url,filesize,mimetype").First()
	if attachment != nil { //文件是否已经存在
		rooturl, _ := DB().Table("merchant_config").Where("keyname", "rooturl").Value("keyvalue")
		attachment["url"] = fmt.Sprintf("%s%s", rooturl, attachment["url"])
		results.Success(context, "文件已上传", attachment, nil)
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
			var ftype int64 = 0
			Insertdata := map[string]interface{}{
				"accountID":  userinfo["accountID"],
				"cuid":       userinfo["cuid"],
				"type":       ftype,
				"pid":        pid,
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
			file_id, _ := DB().Table("client_attachment").Data(Insertdata).InsertGetId()
			//更新排序
			DB().Table("client_attachment").Data(map[string]interface{}{"weigh": file_id}).Where("id", file_id).Update()
			//返回数据
			getdata, _ := DB().Table("client_attachment").Where("id", file_id).Fields("id,pid,name,title,type,url,filesize,mimetype").First()
			rooturl, _ := DB().Table("merchant_config").Where("keyname", "rooturl").Value("keyvalue")
			getdata["url"] = fmt.Sprintf("%s%s", rooturl, getdata["url"])
			results.Success(context, "上传成功", getdata, nil)
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
