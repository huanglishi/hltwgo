package rmold

import (
	"huling/app/common"

	"github.com/gin-gonic/gin"
)

// 公共路由
func ApiCommon(R *gin.Engine) {
	//公共模块
	commonRouter := R.Group("/common")
	{
		//表格排序
		tableRouter := commonRouter.Group("/table")
		{
			tableRouter.POST("/weigh", common.Weigh)
		}
		//文件上传
		uploadfileRouter := commonRouter.Group("/uploadfile")
		{
			uploadfileRouter.POST("/onefile", common.OneFile)
			uploadfileRouter.GET("/getimage", common.GetImage)
			uploadfileRouter.GET("/getimagebase", common.Getimagebase)
		}
		//测试
		testRouter := commonRouter.Group("/api")
		{
			testRouter.POST("/registry", common.Testpath)
		}

	}
}
