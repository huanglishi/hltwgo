package common

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 上传单文件
func OneFile(context *gin.Context) {
	// cid := context.DefaultQuery("cid", "1")
	cid := context.DefaultPostForm("cid", "1")
	// 单个文件
	file, err := context.FormFile("file")
	if err != nil {
		results.Failed(context, "获取数据失败，", err)
		return
	}
	nowTime := time.Now().Unix()      //当前时间
	getuser, _ := context.Get("user") //取值 实现了跨中间件取值
	user := getuser.(*utils.UserClaims)
	//时间查询-获取当天时间
	day_time := time.Now().Format("2006-01-02")
	//文件唯一性
	file_uniname := fmt.Sprintf("%s%s%v", file.Filename, day_time, user.ID)
	sha1_str := md5Str(file_uniname)
	//开始
	day_star, _ := time.Parse("2006-01-02 15:04:05", day_time+" 00:00:00")
	day_star_times := day_star.Unix() //时间戳
	//结束
	day_end, _ := time.Parse("2006-01-02 15:04:05", day_time+" 23:59:59")
	day_end_times := day_end.Unix() //时间戳
	attachment, _ := DB().Table("attachment").Where("uid", user.ID).
		WhereBetween("uploadtime", []interface{}{day_star_times, day_end_times}).
		Where("sha1", sha1_str).Fields("id,title,url").First()
	if attachment != nil { //文件是否已经存在
		context.JSON(200, gin.H{
			"id":       attachment["id"],
			"uid":      sha1_str,
			"name":     attachment["name"],
			"status":   "done",
			"url":      attachment["url"],
			"response": "文件已上传",
			"time":     nowTime,
		})
		context.Abort()
		return
	}
	file_path := fmt.Sprintf("%s%s%s", "resource/uploads/", time.Now().Format("20060102"), "/")
	//如果没有filepath文件目录就创建一个
	if _, err := os.Stat(file_path); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(file_path, os.ModePerm)
		}
	}
	//上传到的路径
	filename_arr := strings.Split(file.Filename, ".")
	name_str := md5Str(fmt.Sprintf("%v%s", nowTime, filename_arr[0])) //组装文件保存名字
	//path := 'resource/uploads/20060102150405test.xlsx'
	file_Filename := fmt.Sprintf("%s%s%s", name_str, ".", filename_arr[1]) //文件加.后缀
	path := file_path + file_Filename
	// fmt.Println("path1:", path) //路径+文件名上传
	// 上传文件到指定的目录
	err = context.SaveUploadedFile(file, path)
	if err != nil {
		context.JSON(200, gin.H{
			"uid":      sha1_str,
			"name":     file.Filename,
			"status":   "error",
			"response": "上传失败",
			"time":     nowTime,
		})
	} else {
		//保存数据
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

		//判断是否是视频-获取封面
		var cover_url string = ""
		if strings.Contains(file.Header["Content-Type"][0], "video/mp4") {
			//封面路径
			// fmt.Println("dir:", path)
			var ffmpegPath string = "ffmpeg"
			vurlpath := dir + strings.Replace("/"+path, "/", "\\", -1)
			cover_url = getLastFrame(vurlpath, file_path+name_str, ffmpegPath)
		}

		Insertdata := map[string]interface{}{
			"accountID":  user.Accountid,
			"cid":        cid,
			"uid":        user.ID,
			"sha1":       sha1_str,
			"title":      filename_arr[0],
			"name":       file.Filename,
			"url":        path,
			"cover_url":  cover_url,
			"storage":    dir + strings.Replace(path, "/", "\\", -1),
			"uploadtime": nowTime,
			"updatetime": nowTime,
			"filesize":   file.Size,
			"mimetype":   file.Header["Content-Type"][0],
		}
		file_id, _ := DB().Table("attachment").Data(Insertdata).InsertGetId()
		context.JSON(200, gin.H{
			"id":        file_id,
			"uid":       sha1_str,
			"name":      file.Filename,
			"status":    "done",
			"url":       path,
			"thumb":     path,
			"cover_url": cover_url,
			"response":  "上传成功",
			// "file":     file.Header,
			"time": nowTime,
		})
	}
}

// 获取视频中最后一帧的图片 url=视频地址,path=图片地址
func getLastFrame(url string, path string, ffmpegPath string) string {
	// fmt.Println("视频url:", url)
	// fmt.Println("图片path:", path)
	// fmt.Println("ffmpeg:", ffmpegPath)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(50000)*time.Millisecond)
	cmd := exec.CommandContext(ctx, ffmpegPath,
		"-loglevel", "error",
		"-y",
		"-ss", "13",
		"-t", "1",
		"-i", url,
		"-vframes", "1",
		path+".jpg")
	defer cancel()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	var outputerror string
	err := cmd.Run()
	// fmt.Println("cmdRun-err0:", err)
	if err != nil {
		outputerror += fmt.Sprintf("lastframecmd—err1:%v;", err)
	}
	if stderr.Len() != 0 {
		outputerror += fmt.Sprintf("lastframestd—err2:%v;", stderr.String())
	}
	if ctx.Err() != nil {
		outputerror += fmt.Sprintf("lastframectx—err3:%v;", ctx.Err())
	}
	// fmt.Println("outputerror4:", outputerror)
	return path + ".jpg"
}

// 显示图片
func GetImage(context *gin.Context) {
	imageName := context.Query("url")
	context.File(imageName)
}

// 显示图片base64
func Getimagebase(context *gin.Context) {
	imageName := context.Query("url")
	file, _ := ioutil.ReadFile(imageName)
	context.Writer.WriteString(string(file))
}
func md5Str(origin string) string {
	m := md5.New()
	m.Write([]byte(origin))
	return hex.EncodeToString(m.Sum(nil))
}

func Testpath(context *gin.Context) {
	log.Printf("测试调度: %v\n", time.Now().Unix())
	results.Success(context, "测试调度2", 100, nil)
}
