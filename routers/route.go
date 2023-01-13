package routers

import (
	rmold "huling/routers/mold"
	"huling/utils/Toolconf"
	"huling/utils/handler"
	utils "huling/utils/tool"
	"net/http"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	R *gin.Engine
)

func LimitHandler(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "您的操作太频繁，请稍后再试！",
				"result":  nil,
			})
			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}

// 验证接口合法性
func validityAPi() gin.HandlerFunc {
	return func(c *gin.Context) {
		var apisecret = Toolconf.AppConfig.String("apisecret")
		encrypt := c.Request.Header.Get("verify-encrypt")
		verifytime := c.Request.Header.Get("verify-time")
		mdsecret := utils.Md5(apisecret + verifytime)
		// fmt.Printf("验证接: %v,新：%v\n", encrypt, mdsecret)
		if encrypt == "" || mdsecret != encrypt {
			if strings.Contains(c.Request.URL.Path, "resource") || strings.Contains(c.Request.URL.Path, "views") ||
				strings.Contains(c.Request.URL.Path, "webmin") || strings.Contains(c.Request.URL.Path, "webadmin") ||
				strings.Contains(c.Request.URL.Path, "webclient") ||
				strings.Contains(c.Request.URL.Path, "webbusiness") ||
				strings.Contains(c.Request.URL.Path, "MP_verify_0HkL8VPApFK29Tb8.txt") ||
				strings.Contains(c.Request.URL.Path, "common/uploadfile/getimage") { //过滤附件访问接口
				c.Next()
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code":    1,
					"message": "您的请求不合法，请按规范请求数据!",
					"result":  nil,
				})
				c.Abort()
			}
		} else {
			c.Next()
		}
	}
}

func init() {
	//Gin 框架在运行的时候默认是debug模式
	runmode := Toolconf.AppConfig.String("runmode")
	if runmode == "dev" {
		gin.SetMode(gin.DebugMode) //ReleaseMode 为方便调试，Gin 框架在运行的时候默认是debug模式，在控制台默认会打印出很多调试日志，上线的时候我们需要关闭debug模式，改为release模式。
	} else if runmode == "pro" {
		gin.SetMode(gin.ReleaseMode)
	} else if runmode == "test" {
		gin.SetMode(gin.TestMode)
	}
	R = gin.Default()
	//0.跨域访问-注意跨域要放在gin.Default下
	str_arr := strings.Split(Toolconf.AppConfig.String("allowurl"), `,`)
	R.Use(cors.New(cors.Config{
		AllowOrigins: str_arr,
		// AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"X-Requested-With", "Content-Type", "verify-encrypt", "Authorization", "verify-time"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	R.MaxMultipartMemory = 8 << 20 // 8 MiB
	//判断接口是否合法
	R.Use(validityAPi())

	//1.限流rate-limit 中间件
	lmt := tollbooth.NewLimiter(100, nil)
	lmt.SetMessage("您访问过于频繁，系统安全检查认为恶意攻击。")
	R.Use(LimitHandler(lmt))
	//2.部署vue项目
	// R.LoadHTMLGlob("viewst/*.html")              // 添加入口index.html
	R.Static("/resource", "./resource")                                                          // 附件
	R.StaticFile("/MP_verify_0HkL8VPApFK29Tb8.txt", "./resource/MP_verify_0HkL8VPApFK29Tb8.txt") // 附件
	R.StaticFile("/favicon.ico", "./resource/favicon.ico")
	// R.Static("/views", "./views")       // 轻站用户后台
	R.Static("/webmin", "./webmin")           // 轻站-手机查看页面
	R.Static("/webadmin", "./webadmin")       // 轻站总后台-admin
	R.Static("/webbusiness", "./webbusiness") // B端后台
	R.Static("/webclient", "./webclient")     // C端后台

	//4.注意 Recover 要尽量放在第一个被加载
	R.Use(handler.Recover)

	//5.验证token
	R.Use(utils.JwtVerify)
	//6.找不到路由
	R.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		c.JSON(400, gin.H{"code": 400, "message": "您" + method + "请求地址：" + path + "不存在！"})
	})
	//7.路由配置文件
	rmold.ApiCommon(R)   //公共接口
	rmold.ApiAdmin(R)    //后台管理模块
	rmold.ApiMerchant(R) //B端后台管理模块
	rmold.ApiMwebh5(R)   //微站手机端页面
	rmold.ApiClient(R)   //C端后台
}
