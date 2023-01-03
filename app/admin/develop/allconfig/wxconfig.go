package allconfig

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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取微信公众号
func Getwxinfo(context *gin.Context) {
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	data, _ := DB().Table("admin_system_wxconfig").Where("admin_id", user.ID).First()
	results.Success(context, "获取微信公众号", data, nil)
}

// 添加配置数据
func SaveWx(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	//查找是否有数据
	find, _ := DB().Table("admin_system_wxconfig").Where("admin_id", user.ID).Value("id")
	if find == nil {
		parameter["createtime"] = time.Now().Unix()
		parameter["admin_id"] = user.ID
		addId, err := DB().Table("admin_system_wxconfig").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("admin_system_wxconfig").
			Data(parameter).
			Where("admin_id", user.ID).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			results.Success(context, "更新成功！", res, nil)
		}
	}
}

// 获取支付
func GetPay(context *gin.Context) {
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	data, _ := DB().Table("admin_system_paymentconfig").Where("admin_id", user.ID).First()
	results.Success(context, "获取支付", data, nil)
}

// 添加支付配置
func SavePay(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	//查找是否有数据
	find, _ := DB().Table("admin_system_paymentconfig").Where("admin_id", user.ID).Value("id")
	if find == nil {
		parameter["createtime"] = time.Now().Unix()
		parameter["admin_id"] = user.ID
		addId, err := DB().Table("admin_system_paymentconfig").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("admin_system_paymentconfig").
			Data(parameter).
			Where("admin_id", user.ID).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			results.Success(context, "更新成功！", res, nil)
		}
	}
}

// 上传图片
func UploadFile(context *gin.Context) {
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
	attachment, _ := DB().Table("attachment").Where("cuid", user.ClientID).
		Where("sha1", sha1_str).Fields("id,name,title,url").First()
	if attachment != nil { //文件是否已经存在
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
			Insertdata := map[string]interface{}{
				"cid":        user.ID,
				"sha1":       sha1_str,
				"title":      filename_arr[0],
				"name":       file.Filename,
				"url":        path,
				"storage":    dir + strings.Replace(path, "/", "\\", -1),
				"uploadtime": nowTime,
				"filesize":   file.Size,
				"mimetype":   file.Header["Content-Type"][0],
			}
			//保存数据
			file_id, _ := DB().Table("attachment").Data(Insertdata).InsertGetId()
			//返回数据
			getdata, _ := DB().Table("attachment").Where("id", file_id).Fields("id,name,title,url").First()
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
