package rmold

import (
	"huling/app/admin/group"
	"huling/app/common"
	"huling/app/merchant/account"
	"huling/app/merchant/dashboard"
	"huling/app/merchant/dept"
	"huling/app/merchant/member"
	microwebgroup "huling/app/merchant/microweb/group"
	"huling/app/merchant/microweb/webmain"
	"huling/app/merchant/ordermanag"
	"huling/app/merchant/role"
	"huling/app/merchant/user"

	"github.com/gin-gonic/gin"
)

/**客户管理后台*/
func ApiMerchant(R *gin.Engine) {
	//项目根模块
	adminPath := R.Group("/merchant")
	{
		//模块公共的接口
		commonrPath := adminPath.Group("/common")
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
			//mapapi
			mapapiRouter := commonrPath.Group("/mapapi")
			{
				mapapiRouter.GET("/weather", common.GetWeather)
			}
		}
		//面板统计
		dashboardPath := adminPath.Group("/dashboard")
		{
			//文件管理
			analysisRouter := dashboardPath.Group("/analysis")
			{
				analysisRouter.GET("/getNumList", dashboard.GetNumList)
				analysisRouter.GET("/getWebPayList", dashboard.GetWebPayList)
				analysisRouter.POST("/savaPayOrder", dashboard.SavaPayOrder)
				analysisRouter.GET("/isHaseMicroweb", dashboard.IsHaseMicroweb)
				analysisRouter.GET("/isHaseSubAccount", dashboard.IsHaseSubAccount)
			}
		}
		//用户模块
		userPath := adminPath.Group("/user")
		{
			userPath.POST("/add", user.AddParam)
			userPath.GET("/userinfo", user.GetInfo)
			userPath.POST("/login", user.Lonin)
			userPath.GET("/refreshtoken", user.Refreshtoken)
			userPath.GET("/logout", user.Logout)
			userPath.POST("/updata", user.Updata)
			userPath.GET("/list", user.QueryParam)
			userPath.DELETE("/del", user.DelParam)
			userPath.POST("/changepassword", user.Changepassword)
		}

		//权限模块
		authPath := adminPath.Group("/auth")
		{
			//路由菜单
			auth := authPath.Group("/rule")
			{
				auth.GET("/getMenuList", user.GetMenuList)
				auth.GET("/getPermCode", user.GetPermCode)
			}
		}
		//系统管理
		systemPath := adminPath.Group("/system")
		{
			//分组
			groupPath := systemPath.Group("/group")
			{
				groupPath.GET("/getlist", group.Getlist_group)
				groupPath.GET("/getparent", group.Getparent_group)
				groupPath.GET("/grouptree", group.Getgroup_tree)
				groupPath.POST("/add", group.Add_group)
				groupPath.DELETE("/del", group.Del_group)
			}
			//部门
			deptPath := systemPath.Group("/dept")
			{
				deptPath.GET("/getDeptList", dept.Getlist)
				deptPath.GET("/getParentList", dept.GetParentList)
				deptPath.POST("/save", dept.Add)
				deptPath.POST("/upLock", dept.UpLock)
				deptPath.POST("/upGrouppid", dept.UpGrouppid)
				deptPath.DELETE("/del", dept.Del)
			}
			//角色
			rolePath := systemPath.Group("/role")
			{
				rolePath.GET("/getList", role.Getlist)
				rolePath.GET("/getParentList", role.GetParentList)
				rolePath.POST("/save", role.Add)
				rolePath.POST("/upLock", role.UpLock)
				rolePath.DELETE("/del", role.Del)
				rolePath.GET("/getMenuList", user.GetRoleMenuList)
			}
			//账号
			accountPath := systemPath.Group("/account")
			{
				accountPath.GET("/getList", account.Getlist)
				accountPath.GET("/getAllRoleList", role.GetAllList)
				accountPath.GET("/getRoleList", role.GetRoleList)
				accountPath.GET("/getLoginLogList", account.GetLoginLogList)
				accountPath.GET("/getAccount", account.GetAccount)
				accountPath.POST("/isAccountExist", account.IsAccountExist)
				accountPath.POST("/save", account.Add)
				accountPath.POST("/upLock", account.UpLock)
				accountPath.POST("/upAvatar", account.UpAvatar)
				accountPath.POST("/changePwd", account.ChangePwd)
				accountPath.DELETE("/del", account.Del)
			}
		}
		//企业微站
		microwebPath := adminPath.Group("/microweb")
		{
			//分组
			groupPath := microwebPath.Group("/group")
			{
				groupPath.GET("/getList", microwebgroup.Getlist)
				groupPath.GET("/getParentList", microwebgroup.GetParentList)
				groupPath.POST("/save", microwebgroup.Add)
				groupPath.POST("/upLock", microwebgroup.UpLock)
				groupPath.POST("/upGrouppid", microwebgroup.UpGrouppid)
				groupPath.DELETE("/del", microwebgroup.Del)
			}
			//微站
			accountPath := microwebPath.Group("/webmain")
			{
				accountPath.GET("/getList", webmain.Getlist)
				accountPath.GET("/getGroupList", microwebgroup.GetCateList)
				accountPath.GET("/getGroupTree", microwebgroup.GetGroupTree)
				accountPath.GET("/getEditData", webmain.GetEditData)
				accountPath.GET("/getTpl", webmain.GetTpl)
				accountPath.GET("/getTplData", webmain.GetTplData)
				accountPath.GET("/getQrTpl", webmain.GetQrTpl)
				accountPath.GET("/getSelectTplGroup", webmain.GetSelectTplGroup)
				accountPath.GET("/getSelectTplList", webmain.GetSelectTplList)
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
			//成员
			memberPath := microwebPath.Group("/member")
			{
				memberPath.GET("/getList", member.Getlist)
				memberPath.POST("/saveMemberData", member.SaveMemberData)
				memberPath.POST("/isAccountExist", member.IsAccountExist)
				memberPath.POST("/upStatus", member.UpStatus)
			}
		}
		//定单管理
		ordermanagPath := adminPath.Group("/ordermanag")
		{
			//套餐
			servicepackagePath := ordermanagPath.Group("/servicepackage")
			{
				servicepackagePath.GET("/getOrderList", ordermanag.GetOrderList)
			}

		}
	}

}
