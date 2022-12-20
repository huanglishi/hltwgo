package handler

import (
	"fmt"
	"huling/app/model"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

func DB() gorose.IOrm {
	return model.DB.NewOrm()
}
func Recover(c *gin.Context) {
	defer func() {
		// token := c.Request.Header.Get("Access-Token")
		// log.Printf("获取消息头: %v\n", token)
		//记录访问日志
		if c.Request != nil {
			db := DB()
			db.Table("logs").Data(map[string]interface{}{"path": c.Request.URL.Path,
				"method": c.Request.Method, "clientIP": GetRequestIP(c), "createtime": time.Now().Unix()}).Insert()
		}
		if r := recover(); r != nil {
			//打印错误堆栈信息
			// typestr := reflect.TypeOf(r)
			// log.Printf("panic返回类型: %v\n", typestr)
			// log.Printf("数组: %v\n", m)
			log.Printf("panic: %s\n", r)
			fmt.Println(recover())
			debug.PrintStack()
			//封装通用json返回
			//c.JSON(http.StatusOK, Result.Fail(errorToString(r)))
			//Result.Fail不是本例的重点，因此用下面代码代替
			// c.JSON(http.StatusOK, gin.H{
			// 	"code":    1,
			// 	"message": errorToString(r),
			// 	"result":    nil,
			// })
			//使用Contains()函数
			erromsg := errorToString(r)
			res1 := strings.Contains(erromsg, "token")
			if res1 {
				c.JSON(http.StatusOK, gin.H{
					"code":    401, //重新登录
					"message": erromsg,
					"result":  "",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code":    1, //错误
					"message": erromsg,
					"result":  nil,
				})
			}
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}

// 获取ip
func GetRequestIP(c *gin.Context) string {
	reqIP := c.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

// recover错误，转string
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}
