package webedit

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取信息内容
func GetMicweb(context *gin.Context) {
	tplid := context.DefaultQuery("tplid", "0")
	micweb_id_df := context.DefaultQuery("micweb_id", "0") //选填
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	var micweb_id int = 0
	if micweb_id_df == "0" {
		micweb_select_id, _ := DB().Table("client_micweb").Where("cuid", user.ClientID).Where("is_select", 1).Value("id")
		if micweb_select_id == nil { //没有选择中-随机查找一个
			micweb_def_id, _ := DB().Table("client_micweb").Where("cuid", user.ClientID).Value("id")
			micweb_id = GetInterfaceToInt(micweb_def_id)
		} else {
			micweb_id = GetInterfaceToInt(micweb_select_id)
		}
	} else {
		micweb_id = GetInterfaceToInt(micweb_id_df)
	}
	if tplid != "0" { //如果是传模板则用模板覆盖-复制
		tpldata, tplerr := DB().Table("client_micweb_tpl_main").Where("id", tplid).Fields("id,cid,home_id,title,details,type,image,footer_tabbar,top_tabbar,side_tabbar,copyright_text").First()
		if tplerr != nil {
			results.Failed(context, "获取模板数据失败，导致无法复制模板", tplerr)
			return
		}
		//1更新底部导航
		_, micweberr := DB().Table("client_micweb").Data(map[string]interface{}{"footer_tabbar": tpldata["footer_tabbar"], "top_tabbar": tpldata["top_tabbar"], "side_tabbar": tpldata["side_tabbar"], "copyright_text": tpldata["copyright_text"]}).
			Where("cuid", user.ClientID).Where("is_select", 1).Update()
		if micweberr != nil {
			results.Failed(context, "获取站点底部导航失败！", micweberr)
			return
		}
		//2复制模板全部页面
		msg, tplpageerr := copyTplpage(tplid, micweb_id, user.ID, user.Accountid)
		if tplpageerr != nil {
			results.Failed(context, msg, tplpageerr)
			return
		}
	}
	data, err := DB().Table("client_micweb").Where("id", micweb_id).Fields("id,title,des,status,approval_err,footer_tabbar,top_tabbar,side_tabbar,copyright_text").First()
	if err != nil {
		results.Failed(context, "获取微站信息失败", err)
	} else {
		if data["selectval"] != nil {
			data["selectval"] = StingToJSON(data["selectval"])
		}
		//字符串转JSON
		// if data["footer_tabbar"] != nil {
		// 	var parameter interface{}
		// 	_ = json.Unmarshal([]byte(data["footer_tabbar"].(string)), &parameter)
		// 	data["footer_tabbar"] = parameter
		// }
		//获取站点下的页面
		list, _ := DB().Table("client_micweb_page").Where("micweb_id", data["id"]).Fields("id,ishome,name,uuid,orderNum").Order("orderNum asc").Get()
		if list != nil {
			data["pagelist"] = list
		} else {
			data["pagelist"] = make([]interface{}, 0)
		}
		//判断是否开通支付
		ispay, _ := DB().Table("client_system_paymentconfig").Where("cuid", user.ClientID).Value("mchAPIv3Key")
		if ispay != nil && ispay != "" {
			data["ispay"] = true
		} else {
			data["ispay"] = false
		}
		results.Success(context, "获取微站信息!", data, nil)
	}
}

// 4.2批量复制模板数据
func copyTplpage(tplid interface{}, micweb_id interface{}, uid interface{}, accountID interface{}) (string, error) {
	//删除之前页面
	_, delerr := DB().Table("client_micweb_page").Where("micweb_id", micweb_id).Delete()
	if delerr != nil {
		return "删除页面错误!", delerr
	} else {
		//批量提交
		pagelist, tplerr := DB().Table("client_micweb_tpl_main_page").Where("main_id", tplid).Order("id asc").Get()
		if tplerr != nil {
			return "获取模板页面失败!", tplerr
		}
		save_arr := []map[string]interface{}{}
		for _, val := range pagelist {
			save_arr = append(save_arr, map[string]interface{}{
				"micweb_id":          micweb_id,
				"uid":                uid,
				"accountID":          accountID,
				"ishome":             val["ishome"],
				"name":               val["name"],
				"orderNum":           val["orderNum"],
				"uuid":               val["uuid"],
				"templateJson":       val["templateJson"],
				"component":          val["component"],
				"banners":            val["banners"],
				"show_banner":        val["show_banner"],
				"show_top_tabbar":    val["show_top_tabbar"],
				"show_side_tabbar":   val["show_side_tabbar"],
				"show_footer_tabbar": val["show_footer_tabbar"],
				"float_btn":          val["float_btn"],
				"show_float_btn":     val["show_float_btn"],
				"returntop":          val["returntop"],
				"show_returntop":     val["show_returntop"],
				"createtime":         time.Now().Unix(),
			})
		}
		DB().Table("client_micweb_page").Data(save_arr).Insert()
		return "添加页面成功", nil
	}
}

// 保存微站页面
func SaveMicwebPabe(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	var f_id float64 = 0
	if parameter["id"] != nil {
		f_id = parameter["id"].(float64)
	}
	//JSON转字符串
	if _, ok := parameter["component"]; ok && parameter["component"] != nil {
		parameter["component"] = JSONMarshalToString(parameter["component"])
	}
	//JSON转字符串
	if _, ok := parameter["templateJson"]; ok && parameter["templateJson"] != nil {
		parameter["templateJson"] = JSONMarshalToString(parameter["templateJson"])
	}
	if f_id == 0 {
		parameter["uid"] = user.ID
		parameter["createtime"] = time.Now().Unix()
		parameter["accountID"] = user.Accountid
		parameter["uuid"] = utils.Getuuid()
		if _, ok := parameter["micweb_id"]; !ok { //不传站点信息在获取
			micweb_id, _ := DB().Table("client_micweb").Where("cuid", user.ClientID).Where("is_select", 1).Value("id")
			parameter["micweb_id"] = micweb_id
		}
		hasehome, _ := DB().Table("client_micweb_page").Where("micweb_id", parameter["micweb_id"]).Where("ishome", 1).Fields("id").First()
		if hasehome == nil {
			parameter["ishome"] = 1
		} else {
			parameter["ishome"] = 0
		}
		addId, err := DB().Table("client_micweb_page").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "保存失败", err)
		} else {
			results.Success(context, "保存成功！", map[string]interface{}{"id": addId, "uuid": parameter["uuid"]}, nil)
		}
	} else {
		res, err := DB().Table("client_micweb_page").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新页面数据失败", err)
		} else {
			results.Success(context, "更新页面数据成功！", f_id, res)
		}
	}
}

// 获取微站页面内容
func GetMicwebPage(context *gin.Context) {
	micweb_id := context.DefaultQuery("micweb_id", "0")
	uuid := context.DefaultQuery("uuid", "0")
	if micweb_id == "0" {
		results.Failed(context, "请传网站的：micweb_id(网站的ID)", nil)
		return
	}
	if uuid == "0" {
		results.Failed(context, "请传页面的：uuid", nil)
		return
	}
	data, err := DB().Table("client_micweb_page").Where("micweb_id", micweb_id).Where("uuid", uuid).First()
	if err != nil {
		results.Failed(context, "获取微站页面内容失败", err)
	} else {
		//字符串转JSON
		if data["component"] != nil {
			data["component"] = StingToJSON(data["component"])
		}
		//字符串转JSON
		if data["templateJson"] != nil {
			var parameter interface{}
			_ = json.Unmarshal([]byte(data["templateJson"].(string)), &parameter)
			data["templateJson"] = parameter
		}
		results.Success(context, "获取微站页面内容", data, nil)
	}
}

// 获取微站页面列表-用于选择跳转
func GetMicwebPageList(context *gin.Context) {
	micweb_id := context.DefaultQuery("micweb_id", "0")
	keyword := context.DefaultQuery("keyword", "")
	if micweb_id == "0" {
		results.Failed(context, "请传网站的：micweb_id(网站的ID)", nil)
		return
	}
	MDB := DB().Table("client_micweb_page")
	if keyword != "" {
		MDB.Where("name", "like", "%"+keyword+"%")
	}
	list, err := MDB.Where("micweb_id", micweb_id).Fields("id,uuid,name").Get()
	if err != nil {
		results.Failed(context, "获取微站页面列表失败", err)
	} else {
		results.Success(context, "获取微站页面列表", list, nil)
	}
}

// 保存微站底部导航
func SaveFooterTabBar(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//JSON转字符串
	if parameter["footer_tabbar"] != nil {
		parameter["footer_tabbar"] = JSONMarshalToString(parameter["footer_tabbar"])
	}
	res, err := DB().Table("client_micweb").
		Data(map[string]interface{}{"footer_tabbar": parameter["footer_tabbar"]}).
		Where("id", parameter["id"]).
		Update()
	if err != nil {
		results.Failed(context, "保存失败", err)
	} else {
		results.Success(context, "保存成功！", res, nil)
	}
}

// 保存微站内容字段不固定
func SaveMicweb(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)

	if _, ok := parameter["id"]; !ok || parameter["id"] == "" {
		results.Failed(context, "参数id不能空", nil)
	} else {
		id := parameter["id"]
		// delete(parameter, "id")
		res, err := DB().Table("client_micweb").
			Where("id", id).
			Data(parameter).
			Update()
		if err != nil {
			results.Failed(context, "保存失败", err)
		} else {
			results.Success(context, "保存成功！", res, nil)
		}
	}
}

// 删除页面
func DelMicwebPage(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	delpage, _ := DB().Table("client_micweb_page").Where("id", parameter["id"]).Fields("micweb_id,ishome").First()
	res2, err := DB().Table("client_micweb_page").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除页面失败", err)
	} else {
		//设置id最小为首页
		var one int = 1
		if GetInterfaceToInt(delpage["ishome"]) == one {
			homeid, _ := DB().Table("client_micweb_page").Where("micweb_id", delpage["micweb_id"]).Order("id asc").Value("id")
			DB().Table("client_micweb_page").Data(map[string]interface{}{"ishome": 1}).Where("id", homeid).Update()
		}
		results.Success(context, "删除页面成功！", res2, nil)
	}
}

// 更新页面是否为首页
func UpPageIshome(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res, err := DB().Table("client_micweb_page").
		Data(map[string]interface{}{"ishome": parameter["ishome"]}).
		Where("id", parameter["id"]).
		Update()
	if err != nil {
		results.Failed(context, "更新失败", err)
	} else {
		results.Success(context, "更新成功！", res, nil)
	}
}
