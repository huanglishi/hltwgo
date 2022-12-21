package rmold

import (
	"huling/app/mwebh5"
	"huling/app/mwebh5/wxpay"

	"github.com/gin-gonic/gin"
)

/**微站前端*/
func ApiMwebh5(R *gin.Engine) {
	//项目根模块
	microwebPath := R.Group("/mwebh5")
	{
		//微站
		webdataPath := microwebPath.Group("/webdata")
		{
			//获取微站数据
			webdataPath.GET("/getData", mwebh5.GetData)
			webdataPath.GET("/getPreviewData", mwebh5.GetPreviewData)
			//新版
			webdataPath.GET("/getHome", mwebh5.GetHome)
			webdataPath.GET("/getPage", mwebh5.GetPage)
			webdataPath.POST("/addVisitRecord", mwebh5.AddVisitRecord)
			//文章
			webdataPath.GET("/getArticle", mwebh5.GetArticle)
			//产品
			webdataPath.GET("/getProduct", mwebh5.GetProduct)
			//表单字段
			webdataPath.POST("/saveForm", mwebh5.SaveForm)
			//登录
			webdataPath.POST("/register", mwebh5.Register)
			webdataPath.POST("/lonin", mwebh5.Lonin)
			webdataPath.POST("/upUserInfo", mwebh5.UpUserInfo)
			//订单
			webdataPath.GET("/getOrderList", mwebh5.GetOrderList)
			webdataPath.GET("/getOrderDetail", mwebh5.GetOrderDetail)
			webdataPath.POST("/addOrder", mwebh5.AddOrder)
			//收货地址
			webdataPath.POST("/saveAddress", mwebh5.SaveAddress)
			webdataPath.GET("/getAddressList", mwebh5.GetAddressList)
			webdataPath.GET("/getAddress", mwebh5.GetAddress)
			webdataPath.DELETE("/delAddress", mwebh5.DelAddress)
			//上传图片
			webdataPath.POST("/uploadImage", mwebh5.UploadImage)
			//模板
			webdataPath.GET("/getTplPage", mwebh5.GetTplPage)
			//整站
			webdataPath.GET("/getWebtpl", mwebh5.GetWebtpl)
			webdataPath.GET("/getWebtplPage", mwebh5.GetWebtplPage)
			//功能木块接口
			webdataPath.GET("/searchAll", mwebh5.SearchAll)

		}
		//微信支付
		wxpayPath := microwebPath.Group("/wxpay")
		{
			//获取h5支付-h5_url
			wxpayPath.GET("/h5url", wxpay.Geth5url)
			wxpayPath.POST("/submitOrder", wxpay.SubmitOrder)
		}
	}

}
