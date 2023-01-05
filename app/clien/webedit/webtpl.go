package webedit

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

// 1 获取模板分类
func GetTplGroup(context *gin.Context) {
	dList, _ := DB().Table("client_micweb_tpl_group").Where("status", 0).Fields("id,name,pid,remark").Order("weigh asc").Get()
	dList = utils.GetTreeArray(dList, 0, "")
	results.Success(context, "获取模板分类", dList, nil)
}

// 2 获取网站模板列表
func GetWebTplList(context *gin.Context) {
	weblist, err := DB().Table("client_micweb_tpl_main").Fields("id,cid,home_id,title,details,type,image").Order("id desc").Get()
	if err != nil {
		results.Failed(context, "获取网站模板失败", err)
	} else {
		results.Success(context, "获取网站模板列表", weblist, nil)
	}
}

// 3 获取网站模板列表
func GetWebTpl(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	if id == "0" {
		results.Failed(context, "请传参数id", nil)
	} else {
		data, err := DB().Table("client_micweb_tpl_main").Where("id", id).Fields("id,cid,home_id,title,details,type,image,footer_tabbar").First()
		if err != nil {
			results.Failed(context, "获取网站模板失败", err)
		} else {
			pagelist, _ := DB().Table("client_micweb_tpl_main_page").Where("main_id", id).Order("id desc").Get()
			if pagelist != nil {
				data["pagelist"] = pagelist
			} else {
				data["pagelist"] = make([]gorose.Data, 0)
			}
			results.Success(context, "获取网站模板数据", data, nil)
		}
	}
}

// 4 添加-整站模板
func SaveWebTpl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	parameter["cuid"] = user.ClientID
	parameter["accountID"] = user.Accountid
	parameter["createtime"] = time.Now().Unix()
	//获取子页面数据
	if _, ok := parameter["micweb_id"]; !ok {
		results.Failed(context, "请传参数网站的micweb_id", nil)
	} else {
		micweb_id := parameter["micweb_id"]
		delete(parameter, "micweb_id")
		if _, ok := parameter["footer_tabbar"]; ok {
			parameter["footer_tabbar"] = JSONMarshalToString(parameter["footer_tabbar"])
		}
		if _, ok := parameter["top_tabbar"]; ok {
			parameter["top_tabbar"] = JSONMarshalToString(parameter["top_tabbar"])
		}
		if _, ok := parameter["side_tabbar"]; ok {
			parameter["side_tabbar"] = JSONMarshalToString(parameter["side_tabbar"])
		}
		addId, err := DB().Table("client_micweb_tpl_main").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			savealltplpage(micweb_id, addId) //添加页面
			results.Success(context, "添加成功！", addId, nil)
		}
	}
}

// 4.2批量处理模板页面数据
func savealltplpage(micweb_id interface{}, main_id interface{}) {
	//批量提交
	pagelist, _ := DB().Table("client_micweb_page").Where("micweb_id", micweb_id).Order("id asc").Get()
	save_arr := []map[string]interface{}{}
	for _, val := range pagelist {
		save_arr = append(save_arr, map[string]interface{}{
			"main_id":            main_id,
			"item_id":            val["id"],
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
		})
	}
	DB().Table("client_micweb_tpl_main_page").Data(save_arr).Insert()
}

// 删除-整站-假删
func DelWebTpl(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_micweb_tpl_main").Where("id", parameter["id"]).Data(map[string]interface{}{"isdel": 1, "delTime": time.Now().Unix()}).Update()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
}

// 删除-整站-真删
func DelWebTpl_real(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_micweb_tpl_main").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		DB().Table("client_micweb_tpl_main_page").Where("main_id", parameter["id"]).Delete()
		results.Success(context, "删除成功！", res2, nil)
	}
}
