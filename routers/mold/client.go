package rmold

import (
	"huling/app/clien/article"
	"huling/app/clien/file"
	"huling/app/clien/form"
	"huling/app/clien/home"
	"huling/app/clien/members"
	"huling/app/clien/microweb"
	"huling/app/clien/product"
	"huling/app/clien/resource"
	"huling/app/clien/system"
	"huling/app/clien/user"
	"huling/app/clien/webedit"
	"huling/app/common"
	microwebgroup "huling/app/merchant/microweb/group"
	"huling/app/merchant/microweb/webmain"

	"github.com/gin-gonic/gin"
)

/**C端API*/
func ApiClient(R *gin.Engine) {
	//项目根模块
	clienPath := R.Group("/client")
	{
		//模块公共的接口
		commonrPath := clienPath.Group("/common")
		{
			//文件管理
			attachmentRouter := commonrPath.Group("/attachment")
			{
				attachmentRouter.POST("/addCate", common.CateAdd)
				attachmentRouter.POST("/upfile", common.Upfile)
				attachmentRouter.DELETE("/delCate", common.DelCate)
				attachmentRouter.GET("/getList", common.GetCatelist)
				attachmentRouter.GET("/getImgList", common.GetImgList)
				attachmentRouter.DELETE("/delfile", common.Delfile)
			}
			//消息管理
			messageRouter := commonrPath.Group("/message")
			{
				messageRouter.GET("/getList", common.GetMessagelist)
				messageRouter.POST("/setRead", common.SetRead)
			}
		}
		//用户模块
		userPath := clienPath.Group("/user")
		{
			userPath.POST("/add", user.AddParam)
			userPath.GET("/userinfo", user.GetInfo)
			userPath.POST("/login", user.Lonin)
			userPath.GET("/refreshtoken", user.Refreshtoken)
			userPath.GET("/logout", user.Logout)
			userPath.POST("/updata", user.Updata)
			userPath.GET("/list", user.QueryParam)
			userPath.DELETE("/del", user.DelParam)
			userPath.POST("/changePwd", user.ChangePwd)
			//用户数据
			userPath.GET("/getUserData", user.GetUserData)
			userPath.POST("/upUserInfo", user.UpUserInfo)
			userPath.POST("/upAvatar", user.UpAvatar)
		}
		//权限模块
		authPath := clienPath.Group("/auth")
		{
			//路由菜单
			auth := authPath.Group("/rule")
			{
				auth.GET("/getMenuList", user.GetMenuList)
				auth.GET("/getPermCode", user.GetPermCode)
			}
		}
		//首页
		homePath := clienPath.Group("/home")
		{
			//统计
			modelsPath := homePath.Group("/models")
			{
				modelsPath.GET("/getModels", home.GetModels)
			}
			//基础-文章-微站
			basePath := homePath.Group("/base")
			{
				basePath.GET("/getArticleList", home.GetArticleList)
				basePath.GET("/getArticle", home.GetArticle)
				basePath.POST("/pushStar", home.PushStar)
				basePath.GET("/getMicweb", home.GetMicweb)
				basePath.POST("/saveMicweb", home.SaveMicweb)
				basePath.POST("/publishMicweb", home.PublishMicweb)
			}
			//首页模板
			micwebtplPath := homePath.Group("/micwebtpl")
			{
				micwebtplPath.GET("/getTplGroup", home.GetTplGroup)
				micwebtplPath.GET("/getTpl", home.GetTpl)
				micwebtplPath.GET("/getCustomtpl", home.GetCustomtpl)
				micwebtplPath.GET("/getCustomtpldeltail", home.GetCustomtpldeltail)
				micwebtplPath.POST("/saveCustom", home.SaveCustom)
			}
		}
		//企业微站
		microwebPath := clienPath.Group("/microweb")
		{
			//微站-编辑器-旧版
			accountPath := microwebPath.Group("/webmain")
			{
				accountPath.GET("/getSelectTplGroup", webmain.GetSelectTplGroup)
				accountPath.GET("/getSelectTplList", webmain.GetSelectTplList)
				accountPath.GET("/getGroupTree", microwebgroup.GetGroupTree)
				accountPath.GET("/getEditData", webmain.GetEditData)

				accountPath.GET("/getList", webmain.Getlist)
				accountPath.GET("/getGroupList", microwebgroup.GetCateList)
				accountPath.GET("/getTpl", webmain.GetTpl)
				accountPath.GET("/getTplData", webmain.GetTplData)
				accountPath.GET("/getQrTpl", webmain.GetQrTpl)
				accountPath.GET("/getMaterial", webmain.GetMaterial)
				accountPath.POST("/saveWebData", webmain.SaveWebData)
				accountPath.POST("/upStatus", webmain.UpStatus)
				accountPath.POST("/savePageData", webmain.SavePageData)
				accountPath.POST("/saveAllPageData", webmain.SaveAllPageData)
				accountPath.POST("/addTpl", webmain.AddTpl)
				accountPath.POST("/addWebTpl", webmain.AddWebTpl)
				accountPath.POST("/upLock", webmain.UpLock)
				accountPath.POST("/addQrtpl", webmain.AddQrtpl)
				accountPath.POST("/changHomePage", webmain.ChangHomePage)
				accountPath.POST("/addMaterial", webmain.AddMaterial)
				accountPath.DELETE("/del", webmain.Del)
				accountPath.DELETE("/delPage", webmain.DelPage)
				accountPath.DELETE("/delQrTpl", webmain.DelQrTpl)
				accountPath.DELETE("/delWebTpl", webmain.DelWebTpl)
				accountPath.DELETE("/delTpl", webmain.DelTpl)

			}
			//模板
			webtplPath := microwebPath.Group("/webtpl")
			{
				webtplPath.GET("/getSelectTplGroup", microweb.GetSelectTplGroup)
				webtplPath.GET("/getSelectTplList", microweb.GetSelectTplList)
				webtplPath.DELETE("/delWebTpl", microweb.DelWebTpl)
			}
			//编辑器接口
			webeditPath := microwebPath.Group("/webedit")
			{
				//获取文章
				webeditPath.GET("/getArticleList", webedit.GetArticleList)
				webeditPath.GET("/getArtic", webedit.GetArticle)
				webeditPath.GET("/getArticleCate", webedit.GetArticleCate)
				//获取产品
				webeditPath.GET("/getProductList", webedit.GetProductList)
				webeditPath.GET("/getProduct", webedit.GetProduct)
				webeditPath.GET("/getProductCate", webedit.GetProductCate)
				//获取表单
				webeditPath.GET("/getFormList", webedit.GetFormList)
				webeditPath.GET("/getFormField", webedit.GetFormField)
				webeditPath.GET("/getFormRule", webedit.GetFormRule)
				//轻站信息
				webeditPath.POST("/saveFooterTabBar", webedit.SaveFooterTabBar)
				webeditPath.POST("/saveMicweb", webedit.SaveMicweb)
				webeditPath.GET("/getMicweb", webedit.GetMicweb)
				webeditPath.GET("/getMicwebPage", webedit.GetMicwebPage)
				webeditPath.DELETE("/delMicwebPage", webedit.DelMicwebPage)
				webeditPath.POST("/saveMicwebPabe", webedit.SaveMicwebPabe)
				webeditPath.POST("/upPageIshome", webedit.UpPageIshome)
				webeditPath.GET("/getMicwebPageList", webedit.GetMicwebPageList)
				//地址
				webeditPath.POST("/saveAddress", webedit.SaveAddress)
				webeditPath.GET("/getAddressList", webedit.GetAddressList)
				webeditPath.GET("/getAddress", webedit.GetAddress)
				webeditPath.DELETE("/delAddress", webedit.DelAddress)
				//客服
				webeditPath.POST("/saveService", webedit.SaveService)
				webeditPath.GET("/getServiceList", webedit.GetServiceList)
				webeditPath.GET("/getService", webedit.GetService)
				webeditPath.DELETE("/delService", webedit.DelService)
				//模板-整站
				webeditPath.GET("/getTplGroup", webedit.GetTplGroup)
				webeditPath.GET("/getWebTplList", webedit.GetWebTplList)
				webeditPath.GET("/getWebTpl", webedit.GetWebTpl)
				webeditPath.POST("/saveWebTpl", webedit.SaveWebTpl)
				webeditPath.DELETE("/delWebTpl", webedit.DelWebTpl)
				//模板-单页
				webeditPath.GET("/getTplPageGroup", webedit.GetTplPageGroup)
				webeditPath.GET("/getTplpage", webedit.GetTplpage)
				webeditPath.POST("/saveTplpage", webedit.SaveTplpage)
				webeditPath.DELETE("/delTplpage", webedit.DelTplpage)

			}
		}
		//1.系统设置
		systemPath := clienPath.Group("/system")
		{
			//支付配置
			paymentConfigPath := systemPath.Group("/paymentConfig")
			{
				paymentConfigPath.GET("/getPayIinfo", system.GetPayIinfo)
				paymentConfigPath.POST("/savePay", system.SavePay)
				paymentConfigPath.POST("/uploadFile", system.UploadFile)
			}
			//接口请求
			apitestPath := systemPath.Group("/apitest")
			{
				//分组
				apitestPath.GET("/getGroupList", system.GetGroupList)
				apitestPath.POST("/saveGroup", system.SaveGroup)
				apitestPath.POST("/upStatus", system.UpStatus)
				apitestPath.GET("/getFormGroupList", system.GetFormGroupList)
				apitestPath.DELETE("/delGroup", system.DelGroup)
				//接口数据
				apitestPath.GET("/getApiList", system.GetApiList)
				apitestPath.GET("/getDBField", system.GetDBField)
				apitestPath.POST("/saveData", system.SaveData)
				apitestPath.DELETE("/delData", system.DelData)
				apitestPath.POST("/upLock", system.UpLock)
			}
		}
		//2.附件
		filePath := clienPath.Group("/file")
		{
			//管理数据
			managePath := filePath.Group("/manage")
			{
				managePath.GET("/getFiles", file.GetFiles)
				managePath.POST("/saveFile", file.SaveFile)
				managePath.POST("/uploadFile", file.UploadFile)
				managePath.POST("/upFile", file.UpFile)
				managePath.POST("/upImgPid", file.UpImgPid)
				managePath.DELETE("/delFile", file.DelFile)
				managePath.GET("/getCateList", file.GetCateList)
				managePath.GET("/getPicture", file.GetPicture)
			}
		}
		//附件资源
		resourcePath := clienPath.Group("/resource")
		{
			//管理数据
			managePath := resourcePath.Group("/manage")
			{
				managePath.POST("/upFile", resource.UpFile)
				managePath.DELETE("/delFile", resource.DelFile)
				managePath.DELETE("/delTest", resource.DelTest)
				managePath.GET("/getPicture", resource.GetPicture)
			}
			//管理数据
			catePath := resourcePath.Group("/cate")
			{
				catePath.GET("/getList", resource.Getlist)
				catePath.GET("/getParentList", resource.GetParentList)
				catePath.POST("/upFile", resource.UpFile)
				catePath.POST("/saveFile", file.SaveFile)
				catePath.DELETE("/delFile", resource.DelFileAndImg)
			}
		}
		//3.文章
		articlePath := clienPath.Group("/article")
		{
			//文章分类
			catePath := articlePath.Group("/cate")
			{
				catePath.GET("/getList", article.GetCateList)
				catePath.POST("/saveCate", article.SaveCate)
				catePath.POST("/upStatus", article.UpCateStatus)
				catePath.GET("/getFormCateList", article.GetFormCateList)
				catePath.DELETE("/delCate", article.DelCate)
			}
			//文章管理
			paymentConfigPath := articlePath.Group("/manage")
			{
				paymentConfigPath.GET("/getList", article.GetList)
				paymentConfigPath.POST("/saveArticle", article.SaveArticle)
				paymentConfigPath.DELETE("/delArticle", article.DelArticle)
				paymentConfigPath.POST("/upLock", article.UpLock)
				paymentConfigPath.GET("/getArticle", article.GetArticle)
				paymentConfigPath.POST("/upTop", article.UpTop)
			}
		}
		//4.产品
		productPath := clienPath.Group("/product")
		{
			//分类
			catePath := productPath.Group("/cate")
			{
				catePath.GET("/getList", product.GetCateList)
				catePath.POST("/saveCate", product.SaveCate)
				catePath.POST("/upStatus", product.UpCateStatus)
				catePath.GET("/getFormCateList", product.GetFormCateList)
				catePath.DELETE("/delCate", product.DelCate)
			}
			//管理
			managePath := productPath.Group("/manage")
			{
				managePath.GET("/getList", product.GetList)
				managePath.POST("/saveProduct", product.SaveProduct)
				managePath.DELETE("/delProduct", product.DelProduct)
				managePath.POST("/upLock", product.UpLock)
				managePath.GET("/getProduct", product.GetProduct)
				managePath.POST("/upTop", product.UpTop)
			}
			//参数
			proPath := productPath.Group("/pro")
			{
				proPath.GET("/getList", product.GetproList)
				proPath.POST("/savePro", product.SavePro)
				proPath.DELETE("/delPro", product.DelPro)
				proPath.POST("/upPro", product.UpPro)
				proPath.POST("/upWeigh", product.UpWeigh)
			}
			//参数值
			prolistPath := productPath.Group("/prolist")
			{
				prolistPath.GET("/getProlist", product.GetProlist)
				prolistPath.POST("/saveProlist", product.SaveProlist)
				prolistPath.DELETE("/delProlist", product.DelProlist)
				prolistPath.POST("/upProlist", product.UpProlist)
				prolistPath.POST("/upWeighlist", product.UpWeighlist)
			}
			//标签
			labelPath := productPath.Group("/label")
			{
				labelPath.GET("/getList", product.GetLabelList)
				labelPath.GET("/getFormLabelList", product.GetFormLabelList)
				labelPath.POST("/saveLabel", product.SaveLabel)
				labelPath.POST("/upStatus", product.UpLabelStatus)
				labelPath.DELETE("/delLabel", product.DelLabel)
			}
			//产品订单
			orderPath := productPath.Group("/order")
			{
				orderPath.GET("/getList", product.GetOrderList)
				orderPath.GET("/getOrder", product.GetOrder)
				orderPath.POST("/upOrder", product.UpOrder)
				orderPath.POST("/upOrderField", product.UpOrderField)
			}
		}
		//5.会员
		memberPath := clienPath.Group("/member")
		{
			//分类
			groupPath := memberPath.Group("/group")
			{
				groupPath.GET("/getList", members.GetGroupList)
				groupPath.GET("/getGroupList", members.GetGroupFormList)
				groupPath.POST("/saveGroup", members.SaveGroup)
				groupPath.POST("/upStatus", members.UpGroupStatus)
				groupPath.DELETE("/delGroup", members.DelGroup)
			}
			//管理
			managePath := memberPath.Group("/manage")
			{
				managePath.GET("/getList", members.GetList)
				managePath.POST("/upLock", members.UpLock)
				managePath.POST("/saveMember", members.SaveMember)
			}
		}
		//6.表单
		formPath := clienPath.Group("/form")
		{
			//表单管理
			managePath := formPath.Group("/manage")
			{
				managePath.GET("/getList", form.GetList)
				managePath.POST("/saveForm", form.SaveForm)
				managePath.POST("/upLock", form.UpLock)
				managePath.DELETE("/delForm", form.Del)
			}
			//表单项
			itemPath := formPath.Group("/item")
			{
				itemPath.GET("/getItemList", form.GetItemList)
				itemPath.POST("/saveItem", form.SaveItem)
				itemPath.POST("/upItem", form.UpItem)
				itemPath.POST("/upRequired", form.UpRequired)
				itemPath.POST("/upWeigh", form.UpWeigh)
				itemPath.DELETE("/delItem", form.DelItem)
			}
			//表单规则
			rulePath := formPath.Group("/rule")
			{
				rulePath.GET("/getRuleAndselectData", form.GetRuleAndselectData)
				rulePath.POST("/saveRule", form.SaveRule)
				rulePath.DELETE("/delRuleItem", form.DelRuleItem)
			}
			//表单数据
			dataPath := formPath.Group("/data")
			{
				dataPath.GET("/getFormField", form.GetFormField)
				dataPath.GET("/getFormDataList", form.GetFormDataList)
				itemPath.DELETE("/delData", form.DelData)
			}
		}

	}

}
