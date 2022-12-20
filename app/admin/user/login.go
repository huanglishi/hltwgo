package user

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 登录
func Lonin(context *gin.Context) {
	// username := context.PostForm("username")
	// password := context.PostForm("password")
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
	userdata, err := DB().Table("admin_user").Fields("id,password,salt,name,avatar").Where("username", username).First()
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
	//token
	token := utils.GenerateToken(&utils.UserClaims{
		ID:             userdata["id"].(int64),
		Accountid:      0,
		Name:           userdata["name"].(string),
		Username:       username,
		StandardClaims: jwt.StandardClaims{},
	})
	// log.Printf("TOken解析: %v\n", token)
	DB().Table("admin_user").Where("id", userdata["id"]).Data(map[string]interface{}{"status": 1, "lastLoginTime": time.Now().Unix(), "lastLoginIp": utils.GetRequestIP(context)}).Update()
	//登录日志
	DB().Table("login_logs").
		Data(map[string]interface{}{"type": 1, "uid": userdata["id"], "out_in": "in",
			"createtime": time.Now().Unix(), "loginIP": utils.GetRequestIP(context)}).
		Insert()
	results.Success(context, "登录成功！", token, nil)
	// strtoken := string(token)
	// tokenstu := utils.ParseToken(strtoken)
	// log.Printf("TOken解析: %v\n", tokenstu.Name)
}

// 退出登录
func Logout(context *gin.Context) {
	getuser, ok := context.Get("user") //取值 实现了跨中间件取值
	if !ok {
		results.Failed(context, "用户id不存在！", ok)
		return
	}
	user := getuser.(*utils.UserClaims)
	res, err := DB().Table("admin_user").Where("id", user.ID).Data(map[string]interface{}{"status": 0}).Update()
	if err != nil {
		results.Failed(context, "退出登录失败！", err)
	} else {
		//登录日志
		DB().Table("login_logs").
			Data(map[string]interface{}{"type": 1, "uid": user.ID, "out_in": "out",
				"createtime": time.Now().Unix(), "loginIP": utils.GetRequestIP(context)}).
			Insert()
		results.Success(context, "退出登录成功！", res, nil)
	}
}

// 刷新token
func Refreshtoken(context *gin.Context) {
	// token := context.PostForm("token")
	token := context.Request.Header.Get("Authorization")
	newtoken := utils.Refresh(token)
	results.Success(context, "刷新token", newtoken, nil)
}
