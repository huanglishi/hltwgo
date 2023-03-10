package utils

import (
	"huling/utils/Toolconf"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 用户信息类，作为生成token的参数
type UserClaims struct {
	ID        int64  `json:"id"`
	Accountid int64  `json:"accountid"`
	ClientID  int64  `json:"clientid"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	//jwt-go提供的标准claim
	jwt.StandardClaims
}

var (
	//自定义的token秘钥
	secret = []byte("16849841325189456f489")
	//该路由下不校验token
	noVerify = []string{
		"/admin/user/login",
		"/admin/user/refreshtoken",
		"/merchant/user/login",
		"/merchant/user/refreshtoken",
		"/common/uploadfile/getimage",
		"/common/uploadfile/getimagebase",
		"/common/api/registry",
		"/client/user/login",
		"/client/cancelapp/user/login",
		"/client/user/refreshtoken",
		"/MP_verify_0HkL8VPApFK29Tb8.txt",
		"/mwebh5/wxpay/paynotify",
		"/mwebh5/wxpay/paynotify1",
	}
	// effectTime = 2 * time.Minute //两分钟
)

// token有效时间（纳秒）
var Expirre = Toolconf.AppConfig.String("tokenouttime")
var effectTime = time.Duration(getiInt(Expirre)) * time.Minute //分钟单位
// var effectTime = time.Duration(getiInt(Expirre)) * time.Hour //2小时的时间
// 写个返回int64-默认2个小时
func getiInt(num string) int64 {
	intnum, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return 2 * 60 //默认2个小时
	} else {
		return intnum
	}
}

// 返回超时时间
func TokenOutTime(claims *UserClaims) int64 {
	return time.Now().Add(effectTime).Unix()
}

// 生成token
func GenerateToken(claims *UserClaims) interface{} {
	//设置token有效期，也可不设置有效期，采用redis的方式
	//   1)将token存储在redis中，设置过期时间，token如没过期，则自动刷新redis过期时间，
	//   2)通过这种方式，可以很方便的为token续期，而且也可以实现长时间不登录的话，强制登录
	//本例只是简单采用 设置token有效期的方式，只是提供了刷新token的方法，并没有做续期处理的逻辑
	claims.ExpiresAt = time.Now().Add(effectTime).Unix()
	//生成token
	sign, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		//这里因为项目接入了统一异常处理，所以使用panic并不会使程序终止，如不接入，可使用原始方式处理错误
		//接入统一异常可参考 https://blog.csdn.net/u014155085/article/details/106733391
		panic(err)
	}
	return map[string]interface{}{"sign": sign, "expiresat": claims.ExpiresAt}
}

// 验证token
func JwtVerify(c *gin.Context) {
	//过滤是否验证token
	//c.Request.RequestURI全部地址
	if strings.Contains(c.Request.URL.Path, "swagger") { //过滤文档接口
		return
	} else if strings.Contains(c.Request.URL.Path, "mwebh5") { //过滤前端
		return
	} else if IsContain(noVerify, c.Request.URL.Path) {
		return
	}
	// token := c.Request.Header.Get("Access-Token")
	token := c.GetHeader("Authorization")
	if token == "" {
		panic("token 不存在")
	}
	//验证token，并存储在请求中
	c.Set("user", ParseToken(token))
}

// 解析Token
func ParseToken(tokenString string) *UserClaims {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		panic(err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		panic("The token is invalid")
	}
	return claims
}

// 更新token
func Refresh(tokenString string) interface{} {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		panic(err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		panic("The token is invalid")
	}
	jwt.TimeFunc = time.Now
	claims.StandardClaims.ExpiresAt = time.Now().Add(effectTime).Unix()
	return GenerateToken(claims)
}
