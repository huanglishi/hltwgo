package mwebh5

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 注册用户账号
func Register(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if _, ok := parameter["micweb_id"]; !ok || parameter["micweb_id"] == "" {
		results.Failed(context, "请传参数：micweb_id（网站id）", nil)
	} else {
		micweb, _ := DB().Table("client_micweb").Where("id", parameter["micweb_id"]).Fields("cuid,accountID").First()
		if parameter["password"] == nil || parameter["password"] == "" || parameter["username"] == nil || parameter["username"] == "" {
			results.Failed(context, "请填写账号或密码", nil)
		} else {
			haseuse, _ := DB().Table("client_member").Where("username", parameter["username"]).Fields("id").First()
			if haseuse != nil {
				results.Failed(context, "您输入账号已注册", nil)
			} else {
				parameter["cuid"] = micweb["cuid"]
				parameter["accountID"] = micweb["accountID"]
				delete(parameter, "micweb_id") //删除多余字段
				rnd := rand.New(rand.NewSource(6))
				salt := strconv.Itoa(rnd.Int())
				mdpass := parameter["password"].(string) + salt
				parameter["password"] = utils.Md5(mdpass)
				parameter["salt"] = salt
				parameter["createtime"] = time.Now().Unix()
				addId, err := DB().Table("client_member").Data(parameter).InsertGetId()
				if err != nil {
					results.Failed(context, "注册失败", err)
				} else {
					results.Success(context, "注册成功", addId, nil)
				}
			}
		}
	}
}

// 登录
func Lonin(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter == nil {
		results.Failed(context, "请传参数：账号或密码", nil)
		return
	}
	username := parameter["username"].(string)
	password := parameter["password"].(string)
	if username == "" {
		results.Failed(context, "请提交用户账号！", nil)
		return
	}
	userdata, err := DB().Table("client_member").Fields("id,accountID,cuid,password,salt,name,avatar").Where("username", username).First()
	if userdata == nil || err != nil {
		results.Failed(context, "账号不存在！", nil)
		return
	}
	if userdata["status"] == 2 {
		results.Failed(context, "账号已经被锁住了", nil)
		return
	}
	pass := utils.Md5(password + userdata["salt"].(string))
	if pass != userdata["password"].(string) {
		results.Failed(context, "您输入的密码不正确！", pass+"="+userdata["password"].(string))
		return
	}
	//登录成功
	delete(userdata, "password")
	delete(userdata, "salt")
	results.Success(context, "登录成功！", userdata, nil)
}

// 修改用户资料
func UpUserInfo(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res, err := DB().Table("client_member").
		Data(parameter).
		Where("id", parameter["id"]).
		Update()
	if err != nil {
		results.Failed(context, "修改失败", err)
	} else {
		results.Success(context, "修改成功！", res, nil)
	}
}
