package webmain

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取数据列表-分页先用Limit在page分则无小-type=1是平台
func Getlist(context *gin.Context) {
	cid := context.DefaultQuery("cid", "0")
	title := context.DefaultQuery("title", "")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := DB().Table("merchant_micweb").Fields("id,groupid,title,des,status,createtime")
	MDBCON := DB().Table("merchant_micweb")
	if cid != "0" {
		sub_cids := GetAllChilId(cid)
		allsubids := append(sub_cids, cid)
		MDB = MDB.WhereIn("groupid", allsubids)
		MDBCON = MDBCON.WhereIn("groupid", allsubids)
	}
	if title != "" {
		MDB = MDB.Where("title", "like", "%"+title+"%")
		MDBCON = MDBCON.Where("title", "like", "%"+title+"%")
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			groupname, _ := DB().Table("merchant_micweb_group").Where("id", val["groupid"]).Value("name")
			val["groupname"] = groupname
			//获取入口首页id
			homeid, _ := DB().Table("merchant_micweb_item").Where("micweb_id", val["id"]).Where("ishome", 1).Value("id")
			if homeid != nil {
				val["homeid"] = homeid
			} else {
				nhomeid, _ := DB().Table("merchant_micweb_item").Where("micweb_id", val["id"]).Value("id")
				val["homeid"] = nhomeid
			}
		}
		var totalCount int64
		totalCount, _ = MDBCON.Count()
		results.Success(context, "获取全部列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 2获取所有子级ID
func GetAllChilId(id interface{}) []interface{} {
	var subids []interface{}
	sub_ids, _ := DB().Table("merchant_micweb_group").Where("pid", id).Pluck("id")
	if len(sub_ids.([]interface{})) > 0 {
		for _, sid := range sub_ids.([]interface{}) {
			subids = append(subids, sid)
			subids = append(subids, GetAllChilId(sid)...)
		}
	}
	return subids
}

// 获取父级数据
func GetParentList(context *gin.Context) {
	menuList, _ := DB().Table("merchant_micweb").Fields("id,pid,name").Order("id asc").Get()
	menuList = GetMenuChildrenArray(menuList, 0)
	results.Success(context, "部门父级数据！", menuList, nil)
}

// 保存微站
func SaveWebData(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	var f_id int = 0
	var zero int = 0
	if parameter["id"] != nil {
		f_id = GetInterfaceToInt(parameter["id"])
	}
	parameter["updatetime"] = time.Now().Unix()
	if f_id == zero {
		parameter["accountID"] = user.Accountid
		parameter["createtime"] = time.Now().Unix()
		addId, err := DB().Table("merchant_micweb").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res2, err := DB().Table("merchant_micweb").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			results.Success(context, "更新成功！", f_id, res2)
		}
	}
}

// 更新状态
func UpStatus(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("merchant_micweb").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res2, nil)
	}
}

// 1.1保存单个页面-保存微站信息
func SavePageData(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	//微站保存数据
	var webSave map[string]interface{}
	bdata, _ := json.Marshal(parameter["webdata"])
	_ = json.Unmarshal(bdata, &webSave)
	webSave["uid"] = user.ID
	var f_id float64 = 0
	if webSave["id"] != nil {
		f_id = InterfaceToFloat64(webSave["id"])
	}
	webSave["updatetime"] = time.Now().Unix()
	if f_id == 0 {
		webSave["accountID"] = user.Accountid
		webSave["createtime"] = time.Now().Unix()
		addId, err := DB().Table("merchant_micweb").Data(webSave).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			pageId := saveOnepage(context, parameter["pagedata"], addId, user.ID, user.Accountid)
			results.Success(context, "添加成功！", map[string]interface{}{"webid": addId, "pageid": pageId}, nil)
		}
	} else {
		res2, err := DB().Table("merchant_micweb").
			Data(webSave).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			pageId := saveOnepage(context, parameter["pagedata"], f_id, user.ID, user.Accountid)
			results.Success(context, "更新成功！", map[string]interface{}{"webid": f_id, "pageid": pageId}, res2)
		}
	}
}

// 1.2保存单个页面数据
func saveOnepage(context *gin.Context, pagedata interface{}, micweb_id interface{}, uid interface{}, accountID interface{}) interface{} {
	timestamp := time.Now().Unix()
	b, _ := json.Marshal(&pagedata)
	var mjson map[string]interface{}
	_ = json.Unmarshal(b, &mjson)
	var item_id float64 = 0
	if mjson["id"] != nil {
		item_id = InterfaceToFloat64(mjson["id"])
	}
	mjson["updatetime"] = timestamp
	var rbId interface{}
	if item_id == 0 { //新增
		mjson["micweb_id"] = micweb_id
		mjson["uid"] = uid
		mjson["accountID"] = accountID
		mjson["createtime"] = timestamp
		addId, err := DB().Table("merchant_micweb_item").Data(mjson).InsertGetId()
		if err != nil {
			results.Failed(context, "添加页面数据失败", err)
			context.Abort()
		} else {
			rbId = addId
		}
	} else {
		_, err := DB().Table("merchant_micweb_item").
			Data(mjson).
			Where("id", item_id).
			Update()
		if err != nil {
			results.Failed(context, "更新页面数据失败", err)
			context.Abort()
		} else {
			rbId = item_id
		}
	}
	return rbId
}

// 设置页是否为首页
func ChangHomePage(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//修改全部为未选择
	DB().Table("merchant_micweb_item").Where("micweb_id", parameter["micweb_id"]).Data(map[string]interface{}{"ishome": 0}).Update()
	//更新当前项
	res2, err := DB().Table("merchant_micweb_item").Where("id", parameter["id"]).Data(map[string]interface{}{"ishome": parameter["ishome"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		results.Success(context, "更新成功！", res2, nil)
	}
}

// 2.1批量全部保存微站及页面
func SaveAllPageData(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	//微站保存数据
	var webSave map[string]interface{}
	bdata, _ := json.Marshal(parameter["webdata"])
	_ = json.Unmarshal(bdata, &webSave)
	webSave["uid"] = user.ID
	var f_id float64 = 0
	if webSave["id"] != nil {
		f_id = InterfaceToFloat64(webSave["id"])
	}
	webSave["updatetime"] = time.Now().Unix()
	if f_id == 0 {
		webSave["accountID"] = user.Accountid
		webSave["createtime"] = time.Now().Unix()
		addId, err := DB().Table("merchant_micweb").Data(webSave).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			saveallpage(parameter["pagelist"].([]interface{}), addId, user.ID, user.Accountid)
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res2, err := DB().Table("merchant_micweb").
			Data(webSave).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			saveallpage(parameter["pagelist"].([]interface{}), f_id, user.ID, user.Accountid)
			pageids, _ := DB().Table("merchant_micweb_item").Where("micweb_id", f_id).Fields("id,uuid").Get()
			results.Success(context, "更新成功！", map[string]interface{}{"webid": f_id, "pageids": pageids}, res2)
		}
	}
}

// 2.2批量处理页面数据
func saveallpage(pagelist []interface{}, micweb_id interface{}, uid interface{}, accountID interface{}) {
	//批量提交
	timestamp := time.Now().Unix()
	save_arr := []map[string]interface{}{}
	for _, val := range pagelist {
		b, _ := json.Marshal(&val)
		var mjson map[string]interface{}
		_ = json.Unmarshal(b, &mjson)
		var item_id float64 = 0
		if mjson["id"] != nil {
			item_id = InterfaceToFloat64(mjson["id"])
		}
		mjson["updatetime"] = timestamp
		if item_id == 0 { //新增
			mjson["micweb_id"] = micweb_id
			mjson["uid"] = uid
			mjson["accountID"] = accountID
			mjson["createtime"] = timestamp
			save_arr = append(save_arr, mjson)
		} else {
			DB().Table("merchant_micweb_item").
				Data(mjson).
				Where("id", item_id).
				Update()
		}
	}
	DB().Table("merchant_micweb_item").Data(save_arr).Insert()
}

// 3.1添加-单个模板
func AddTpl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	parameter["accountID"] = user.Accountid
	parameter["createtime"] = time.Now().Unix()
	addId, err := DB().Table("merchant_micweb_tpl").Data(parameter).InsertGetId()
	if err != nil {
		results.Failed(context, "添加失败", err)
	} else {
		results.Success(context, "添加成功！", addId, nil)
	}
}

// 3.2获取模板数据
func GetTpl(context *gin.Context) {
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "20")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	list, _ := DB().Table("merchant_micweb_tpl").Where("accountID", user.Accountid).OrWhere("type", 0).Fields("id,title,image,type").Limit(pageSize).Page(pageNo).Order("id asc").Get()
	results.Success(context, "获取模板数据", list, nil)
}

// 3.3使用模板数据
func GetTplData(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	data, _ := DB().Table("merchant_micweb_tpl").Where("id", id).Fields("templateJson,component").First()
	results.Success(context, "使用模板数据", data, nil)
}

// 4.1添加-整站模板
func AddWebTpl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	parameter["accountID"] = user.Accountid
	parameter["createtime"] = time.Now().Unix()
	//获取子页面数据
	pagelist := parameter["pagelist"]
	delete(parameter, "pagelist")
	addId, err := DB().Table("merchant_micweb_tpl_main").Data(parameter).InsertGetId()
	if err != nil {
		results.Failed(context, "添加失败", err)
	} else {
		savealltplpage(pagelist.([]interface{}), addId)
		results.Success(context, "添加成功！", addId, nil)
	}
}

// 4.2批量处理模板页面数据
func savealltplpage(pagelist []interface{}, main_id interface{}) {
	//批量提交
	save_arr := []map[string]interface{}{}
	for _, val := range pagelist {
		b, _ := json.Marshal(&val)
		var mjson map[string]interface{}
		_ = json.Unmarshal(b, &mjson)
		pageSetup, _ := json.Marshal(mjson["pageSetup"])
		pageComponents, _ := json.Marshal(mjson["pageComponents"])
		//id
		webb, _ := json.Marshal(mjson["pageSetup"])
		var webjson map[string]interface{}
		_ = json.Unmarshal(webb, &webjson)
		save_arr = append(save_arr, map[string]interface{}{
			"main_id":      main_id,
			"item_id":      webjson["id"],
			"templateJson": pageSetup,
			"component":    pageComponents,
		})
	}
	DB().Table("merchant_micweb_tpl_main_page").Data(save_arr).Insert()
}

// 更新状态
func UpLock(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("merchant_micweb").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res2, nil)
	}

}

// 4 删除-全部-微站和页面
func Del(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("merchant_micweb").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		//并且删除页面
		DB().Table("merchant_micweb_item").WhereIn("micweb_id", ids.([]interface{})).Delete()
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 5 删除页面
func DelPage(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("merchant_micweb_item").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 6获取微站编辑页面数据
func GetEditData(context *gin.Context) {
	//当前用户-确保是同账号下的数据
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	//数据ID
	id := context.DefaultQuery("id", "0")
	//1获取微站
	webdata, _ := DB().Table("merchant_micweb").Where("accountID", user.Accountid).Fields("id,groupid,title,des").Where("id", id).First()
	//2获取页面
	pagelist, _ := DB().Table("merchant_micweb_item").Where("accountID", user.Accountid).Where("micweb_id", id).Fields("id,uuid,orderNum,micweb_id,name,details,templateJson,component").Order("orderNum asc").Get()
	results.Success(context, "微站编辑页面数据", map[string]interface{}{"webdata": webdata, "pagelist": pagelist}, nil)
}

// 7.1添加二维码模板
func AddQrtpl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	parameter["accountID"] = user.Accountid
	parameter["createtime"] = time.Now().Unix()
	addId, err := DB().Table("merchant_micweb_qrtpl").Data(parameter).InsertGetId()
	if err != nil {
		results.Failed(context, "添加失败", err)
	} else {
		results.Success(context, "添加成功！", addId, nil)
	}
}

// 7.2获取二维码模板数据
func GetQrTpl(context *gin.Context) {
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	MDB := DB().Table("merchant_micweb_qrtpl").Where("accountID", user.Accountid).OrWhere("type", 0)
	list, _ := MDB.Fields("id,title,image,jsondata").Limit(pageSize).Page(pageNo).Order("id asc").Get()
	var totalCount int64
	totalCount, _ = MDB.Count()
	results.Success(context, "获取二维码模板数据", map[string]interface{}{"list": list, "total": totalCount}, nil)
}

// 7.3 删除二维码模板
func DelQrTpl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("merchant_micweb_qrtpl").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 8.1 获取选择模板分类
func GetSelectTplGroup(context *gin.Context) {
	dList, _ := DB().Table("merchant_micweb_tpl_group").Where("status", 0).Fields("id,name,pid,remark").Order("weigh asc").Get()
	dList = utils.GetTreeArray(dList, 0, "")
	results.Success(context, "获取选择模板分类", dList, nil)
}

// 8.2 获取选择模板
func GetSelectTplList(context *gin.Context) {
	key_pcid := context.DefaultQuery("pcid", "0")
	key_cid := context.DefaultQuery("cid", "0")
	MDB := DB().Table("merchant_micweb_tpl_main").Where("isdel", 0)
	if key_cid != "0" && key_pcid != key_cid {
		MDB = MDB.Where("cid", key_cid)
	} else {
		cids, _ := DB().Table("merchant_micweb_tpl_group").Where("status", 0).Where("pid", key_pcid).Pluck("id")
		MDB = MDB.WhereIn("cid", cids.([]interface{}))
	}
	datalist, _ := MDB.Order("id desc").Get()
	locall_imgurl, _ := DB().Table("merchant_config").Where("keyname", "locall_imgurl").Value("keyvalue")
	for _, val := range datalist {
		if val["image"] != "" && val["image"] != nil {
			val["image"] = strings.Replace(val["image"].(string), "http://localhost:8098", locall_imgurl.(string), -1)
		}
	}
	results.Success(context, "获取网站模板数据", datalist, nil)
}

// 8.3 删除整站模板
func DelWebTpl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("merchant_micweb_tpl_main").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		//删除子页面
		DB().Table("merchant_micweb_tpl_main_page").Where("main_id", parameter["id"]).Delete()
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 删除单个模板
func DelTpl(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("merchant_micweb_tpl").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
