package file

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"huling/utils/results"
	utils "huling/utils/tool"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// 上传图片
func UploadFile(context *gin.Context) {
	// 单个文件
	pid := context.DefaultPostForm("pid", "")
	filetype := context.DefaultPostForm("filetype", "image") //文件类型
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
	// file_uniname := fmt.Sprintf("%s%s%v", file.Filename, day_time, user.ID)
	sha1_str := hex.EncodeToString(m_d5.Sum(nil))
	//查找该用户是否传过
	attachment, _ := DB().Table("client_attachment").Where("cuid", user.ClientID).
		Where("sha1", sha1_str).Fields("id,pid,name,title,type,url,filesize,mimetype,storage,cover_url").First()
	if attachment != nil { //文件是否已经存在
		//更新到最前面
		maxId, _ := DB().Table("client_attachment").Where("cuid", user.ClientID).Order("weigh desc").Value("id")
		if maxId != nil {
			DB().Table("client_attachment").Data(map[string]interface{}{"weigh": maxId.(int64) + 1, "pid": pid}).Where("id", attachment["id"]).Update()
		}
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
			var cover_url string = ""
			if filetype == "video" {
				ftype = 2
				//封面路径
				fmt.Println("dir:", dir)
				var ffmpegPath string = strings.Replace(dir, "\\", "/", -1) + "/" + "ffmpeg/bin/ffmpeg"
				vurlpath := strings.Replace(dir, "\\", "/", -1) + "/" + path
				cover_url = invokeFfmpeg(vurlpath, file_path+name_str, ffmpegPath)
			}
			Insertdata := map[string]interface{}{
				"accountID":  user.Accountid,
				"cuid":       user.ClientID,
				"type":       ftype,
				"pid":        pid,
				"sha1":       sha1_str,
				"title":      filename_arr[0],
				"name":       file.Filename,
				"url":        path,
				"cover_url":  cover_url, //视频封面
				"storage":    dir + "/" + path,
				"createtime": nowTime,
				"filesize":   file.Size,
				"mimetype":   file.Header["Content-Type"][0],
			}
			//保存数据
			file_id, _ := DB().Table("client_attachment").Data(Insertdata).InsertGetId()
			//更新排序
			DB().Table("client_attachment").Data(map[string]interface{}{"weigh": file_id}).Where("id", file_id).Update()
			//返回数据
			getdata, _ := DB().Table("client_attachment").Where("id", file_id).Fields("id,pid,name,title,type,url,filesize,mimetype,storage,cover_url").First()
			results.Success(context, "上传成功", getdata, nil)
		}
	}
	context.Abort()
}

// 视频截取第一帧作为封面
func invokeFfmpeg(urlpath string, path string, ffmpegPath string) string {
	fmt.Println("urlpath:", urlpath)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(14400000)*time.Millisecond)
	cmd := exec.CommandContext(ctx, ffmpegPath,
		"-loglevel", "error",
		"-i", urlpath,
		"-ss", "1",
		"-f", "image2",
		"./"+path+".jpg")
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	defer cancel()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	var outputerror string
	err := cmd.Run()
	if err != nil {
		outputerror += fmt.Sprintf("cmderr:%v;", err)
	}
	if stderr.Len() != 0 {
		outputerror += fmt.Sprintf("stderr:%v;", stderr.String())
	}
	if ctx.Err() != nil {
		outputerror += fmt.Sprintf("ctxerr:%v;", ctx.Err())
	}
	// fmt.Println("invokeFfmpeg err:", outputerror)
	return path + ".jpg"
}

// md5加密
func md5Str(origin string) string {
	m := md5.New()
	m.Write([]byte(origin))
	return hex.EncodeToString(m.Sum(nil))
}
