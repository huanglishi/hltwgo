package rmold

import (
	"huling/app/admin/account"
	"huling/app/admin/dept"
	"huling/app/admin/develop/allconfig"
	dwebtplgroup "huling/app/admin/dictionary/group"
	"huling/app/admin/dictionary/teacharticle"
	tplpagegroup "huling/app/admin/dictionary/tplpage"
	"huling/app/admin/group"
	"huling/app/admin/matter/picture"
	"huling/app/admin/matter/picturecate"
	"huling/app/admin/menu"
	platformaccount "huling/app/admin/platform/account"
	boosmenu "huling/app/admin/platform/boosmenu"
	clientmenu "huling/app/admin/platform/clientmenu"
	platformgroup "huling/app/admin/platform/group"
	mwebapproval "huling/app/admin/platform/mwebapproval"
	packagedesign "huling/app/admin/platform/packagedesign"
	"huling/app/admin/platform/repayment"
	"huling/app/admin/platform/tplaplication"
	"huling/app/admin/platform/webtpldel"
	"huling/app/admin/role"
	"huling/app/admin/user"

	"github.com/gin-gonic/gin"
)

/**管理后台*/
func ApiAdmin(R *gin.Engine) {
	//项目根模块
	adminPath := R.Group("/admin")
	{
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
			//系统菜单
			menuPath := systemPath.Group("/menu")
			{
				menuPath.GET("/getlist", menu.Getlist)
				menuPath.GET("/getParentList", menu.GetParentList)
				menuPath.POST("/saveMenu", menu.Add)
				menuPath.POST("/upMenuLock", menu.UpMenuLock)
				menuPath.DELETE("/delMenu", menu.Del)
			}
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
				rolePath.GET("/getMenuList", menu.GetMenuList)
			}
			//账号
			accountPath := systemPath.Group("/account")
			{
				accountPath.GET("/getList", account.Getlist)
				accountPath.GET("/getParentList", account.GetParentList)
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
		//系统管理
		platformPath := adminPath.Group("/platform")
		{
			//分组
			groupPath := platformPath.Group("/group")
			{
				groupPath.GET("/getList", platformgroup.Getlist)
				groupPath.GET("/getParentList", platformgroup.GetParentList)
				groupPath.POST("/save", platformgroup.Add)
				groupPath.POST("/upLock", platformgroup.UpLock)
				groupPath.DELETE("/del", platformgroup.Del)
			}
			//账号
			accountPath := platformPath.Group("/account")
			{
				accountPath.GET("/getList", platformaccount.Getlist)
				accountPath.GET("/getGroupList", platformgroup.GetParentList)
				accountPath.GET("/getLoginLogList", platformaccount.GetLoginLogList)
				accountPath.GET("/getAccount", platformaccount.GetAccount)
				accountPath.GET("/getMenuList", boosmenu.GetMenuList)
				accountPath.POST("/isAccountExist", platformaccount.IsAccountExist)
				accountPath.POST("/save", platformaccount.Add)
				accountPath.POST("/saveSetting", platformaccount.SaveSetting)
				accountPath.POST("/upLock", platformaccount.UpLock)
				accountPath.DELETE("/del", platformaccount.Del)
			}

			//套餐设计
			packagedesignPath := platformPath.Group("/packagedesign")
			{
				packagedesignPath.GET("/getList", packagedesign.Getlist)
				packagedesignPath.GET("/useList", packagedesign.UseList)
				packagedesignPath.POST("/saveData", packagedesign.SaveData)
				packagedesignPath.DELETE("/delData", packagedesign.Del)
			}
			//轻站审批
			mwebapprovalPath := platformPath.Group("/mwebapproval")
			{
				mwebapprovalPath.GET("/getList", mwebapproval.Getlist)
				mwebapprovalPath.POST("/saveData", mwebapproval.SaveData)
				mwebapprovalPath.POST("/approvalMweb", mwebapproval.ApprovalMweb)
			}
			//缴费续费
			repaymentPath := platformPath.Group("/repayment")
			{
				repaymentPath.GET("/getList", repayment.Getlist)
				repaymentPath.POST("/doResult", repayment.DoResult)
			}
			//缴费续费
			tplaplicationPath := platformPath.Group("/tplaplication")
			{
				tplaplicationPath.GET("/getList", tplaplication.Getlist)
				tplaplicationPath.GET("/getDetail", tplaplication.GetDetail)
				tplaplicationPath.GET("/getTplList", tplaplication.GetTplList)
				tplaplicationPath.POST("/doResult", tplaplication.DoResult)
			}
			//轻站模板回收站
			webtpldelPath := platformPath.Group("/webtpldel")
			{
				webtpldelPath.GET("/getList", webtpldel.Getlist)
				webtpldelPath.POST("/restoreTpl", webtpldel.RestoreTpl)
			}
		}

		//开发者
		developPath := adminPath.Group("/develop")
		{
			//Boss菜单
			bossmenuPath := developPath.Group("/bossmenu")
			{
				bossmenuPath.GET("/getlist", boosmenu.Getlist)
				bossmenuPath.GET("/getParentList", boosmenu.GetParentList)
				bossmenuPath.POST("/saveMenu", boosmenu.Add)
				bossmenuPath.POST("/upMenuLock", boosmenu.UpMenuLock)
				bossmenuPath.DELETE("/delMenu", boosmenu.Del)
			}
			//Boss菜单
			clientmenuPath := developPath.Group("/clientmenu")
			{
				clientmenuPath.GET("/getlist", clientmenu.Getlist)
				clientmenuPath.GET("/getParentList", clientmenu.GetParentList)
				clientmenuPath.POST("/saveMenu", clientmenu.Add)
				clientmenuPath.POST("/upMenuLock", clientmenu.UpMenuLock)
				clientmenuPath.DELETE("/delMenu", clientmenu.Del)
			}
			allconfigPath := developPath.Group("/allconfig")
			{
				allconfigPath.GET("/getwxinfo", allconfig.Getwxinfo)
				allconfigPath.POST("/saveWx", allconfig.SaveWx)

				allconfigPath.GET("/getPay", allconfig.GetPay)
				allconfigPath.POST("/savePay", allconfig.SavePay)
				allconfigPath.POST("/uploadFile", allconfig.UploadFile)
			}

		}
		//素材管理
		matterPath := adminPath.Group("/matter")
		{
			//图片分组
			picturecatePath := matterPath.Group("/picturecate")
			{
				picturecatePath.GET("/getList", picturecate.Getlist)
				picturecatePath.POST("/save", picturecate.Add)
				picturecatePath.POST("/upLock", picturecate.UpLock)
				picturecatePath.DELETE("/del", picturecate.Del)
			}
			//图片
			picturePath := matterPath.Group("/picture")
			{
				picturePath.GET("/getList", picture.Getlist)
				picturePath.GET("/getCateTree", picturecate.GetCateTree)
				picturePath.POST("/uploadFile", picture.UploadFile)
				picturePath.POST("/save", picture.Add)
				picturePath.POST("/upLock", picture.UpLock)
				picturePath.DELETE("/del", picture.Del)
			}
		}
		//系统字典数据
		microwebPath := adminPath.Group("/dictionary")
		{
			//分组-整站
			groupPath := microwebPath.Group("/webtplgroup")
			{
				groupPath.GET("/getList", dwebtplgroup.Getlist)
				groupPath.GET("/getParentList", dwebtplgroup.GetParentList)
				groupPath.POST("/save", dwebtplgroup.Add)
				groupPath.POST("/upLock", dwebtplgroup.UpLock)
				groupPath.POST("/upGrouppid", dwebtplgroup.UpGrouppid)
				groupPath.DELETE("/del", dwebtplgroup.Del)
			}
			//分组=单页面
			tplpagePath := microwebPath.Group("/tplpage")
			{
				tplpagePath.GET("/getList", tplpagegroup.Getlist)
				tplpagePath.POST("/save", tplpagegroup.Add)
				tplpagePath.POST("/upLock", tplpagegroup.UpLock)
				tplpagePath.DELETE("/del", tplpagegroup.Del)
			}
			//教程引导文章
			teacharticlePath := microwebPath.Group("/teacharticle")
			{
				teacharticlePath.GET("/getList", teacharticle.GetList)
				teacharticlePath.GET("/getArticle", teacharticle.GetArticle)
				teacharticlePath.POST("/saveArticle", teacharticle.SaveArticle)
				teacharticlePath.POST("/upLock", teacharticle.UpLock)
				teacharticlePath.DELETE("/del", teacharticle.DelArticle)
			}

		}
	}

}
