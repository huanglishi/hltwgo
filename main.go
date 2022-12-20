package main

import (
	"fmt"
	"huling/app/model"
	. "huling/routers"
	"huling/utils/Toolconf"
	"runtime"
	"strconv"

	_ "github.com/go-sql-driver/mysql" //只执行github.com/go-sql-driver/mysql的init函数
	"github.com/gohouse/gorose/v2"     //数据库操作
	// //swag接口文档
	// _ "huling/docs" // 这里需要引入本地已生成文档
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

var db *gorose.Engin
var runmode = Toolconf.AppConfig.String("runmode")

func InitDB() {
	var err error
	db, err = gorose.Open(&gorose.Config{Driver: "mysql", Dsn: Toolconf.AppConfig.String("db.user") + ":" + Toolconf.AppConfig.String("db.password") + "@tcp(" + Toolconf.AppConfig.String("db.host") + ":" + Toolconf.AppConfig.String("db.port") + ")/" + Toolconf.AppConfig.String("db.name") + "?charset=utf8mb4&parseTime=true&loc=Local", SetMaxOpenConns: 100, SetMaxIdleConns: 10})
	if err != nil {
		fmt.Println("链接数据库错误，请检查数据库链接！")
	}
	if runmode == "div" { //Gin 框架在运行的时候默认是debug模式
		fmt.Printf("数据库已连接:%v\n", Toolconf.AppConfig.String("db.host")+"下的："+Toolconf.AppConfig.String("db.name"))
	}
	model.DB = db
}

// @title           Swagger接口文档
// @version         1.1
// @description     自动生成api文档
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8088
// @BasePath
// @securityDefinitions.basic  BasicAuth
func main() {
	//多核并行任务
	cpu_num, _ := strconv.Atoi(Toolconf.AppConfig.String("cpunum"))
	mycpu := runtime.NumCPU()
	if cpu_num > mycpu { //如果配置cpu核数大于当前计算机核数，则等当前计算机核数
		cpu_num = mycpu
	}
	if cpu_num > 0 {
		if runmode == "dev" { //Gin 框架在运行的时候默认是debug模式
			fmt.Printf("当前计算机CPU核数: %v个,调用：%v个\n", mycpu, cpu_num)
		}
		runtime.GOMAXPROCS(cpu_num)
	} else {
		if runmode == "dev" { //Gin 框架在运行的时候默认是debug模式
			fmt.Printf("当前计算机CPU核数: %v个,调用：%v个\n", mycpu, mycpu)
		}
		runtime.GOMAXPROCS(mycpu)
	}
	InitDB() //连接数据库
	//加入swagger的路由，可以支持在页面访问
	// R.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run("里面不指定端口号默认为8088")
	R.Run(Toolconf.AppConfig.String("httpport"))
}
