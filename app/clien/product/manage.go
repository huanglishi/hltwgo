package product

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取列表
func GetList(context *gin.Context) {
	title := context.DefaultQuery("title", "")
	status := context.DefaultQuery("status", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	whereMap := DB().Table("client_product_manage").Where("cuid", user.ClientID)
	whereMap2 := DB().Table("client_product_manage").Where("cuid", user.ClientID)
	if status != "0" {
		whereMap.Where("status", status)
		whereMap2.Where("status", status)
	}
	if title != "" {
		whereMap.Where("title", "like", "%"+title+"%")
		whereMap2.Where("title", "like", "%"+title+"%")
	}
	list, err := whereMap.Limit(pageSize).Page(pageNo).Fields("id,type,title,images,views,top,releasetime,status").Order("top desc , weigh desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			//分组
			cids, _ := DB().Table("client_product_cid").Where("product_id", val["id"]).Pluck("cid")
			catename, _ := DB().Table("client_product_cate").WhereIn("id", cids.([]interface{})).Pluck("name")
			val["catename"] = catename
			//标签
			lids, _ := DB().Table("client_product_lid").Where("product_id", val["id"]).Pluck("lid")
			labelname, _ := DB().Table("client_product_label").WhereIn("id", lids.([]interface{})).Pluck("name")
			val["labelname"] = labelname
			//多图
			var parameter []interface{}
			_ = json.Unmarshal([]byte(val["images"].(string)), &parameter)
			val["images"] = parameter
		}
		var totalCount int64
		totalCount, _ = whereMap2.Count()
		results.Success(context, "获取列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 添加
func SaveProduct(context *gin.Context) {
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
	if parameter["releasetime"] != nil {
		var LOC, _ = time.LoadLocation("Asia/Shanghai")
		tim, _ := time.ParseInLocation("2006-01-02 15:04:05", parameter["releasetime"].(string), LOC)
		parameter["releasetime"] = tim.Unix()
	}
	list := parameter["cid"]
	labelist := parameter["lid"]
	pro_list := parameter["pro_list"]
	delete(parameter, "cid")
	delete(parameter, "lid")
	delete(parameter, "catename")
	delete(parameter, "pendingStatus")
	delete(parameter, "pro_list")
	if parameter["images"] != nil {
		parameter["images"] = JSONMarshalToString(parameter["images"])
	}
	if f_id == 0 {
		parameter["createtime"] = time.Now().Unix()
		parameter["accountID"] = user.Accountid
		parameter["cuid"] = user.ClientID
		addId, err := DB().Table("client_product_manage").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败！", err)
		} else {
			if addId != 0 {
				DB().Table("client_product_manage").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			addAllCid(list.([]interface{}), addId)
			addAllLid(labelist.([]interface{}), addId)
			addAllPro(pro_list.([]interface{}), addId)
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("client_product_manage").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			addAllCid(list.([]interface{}), f_id)
			addAllLid(labelist.([]interface{}), f_id)
			addAllPro(pro_list.([]interface{}), f_id)
			results.Success(context, "更新成功！", res, user)
		}
	}
}

// 1批量添加分类
func addAllCid(list []interface{}, produc_id interface{}) {
	//批量提交
	DB().Table("client_product_cid").Where("product_id", produc_id).Delete()
	save_arr := []map[string]interface{}{}
	for _, val := range list {
		save_arr = append(save_arr, map[string]interface{}{"product_id": produc_id, "cid": val})
	}
	DB().Table("client_product_cid").Data(save_arr).Insert()
}

// 1批量添加标签
func addAllLid(list []interface{}, produc_id interface{}) {
	//批量提交
	DB().Table("client_product_lid").Where("product_id", produc_id).Delete()
	save_arr := []map[string]interface{}{}
	for _, val := range list {
		save_arr = append(save_arr, map[string]interface{}{"product_id": produc_id, "lid": val})
	}
	DB().Table("client_product_lid").Data(save_arr).Insert()
}

// 2批量添加参数
func addAllPro(list []interface{}, produc_id interface{}) {
	//批量提交
	DB().Table("client_product_manage_pro_val").Where("product_id", produc_id).Delete()
	save_arr := []map[string]interface{}{}
	for _, val := range list {
		webb, _ := json.Marshal(val)
		var webjson map[string]interface{}
		_ = json.Unmarshal(webb, &webjson)
		save_arr = append(save_arr, map[string]interface{}{"product_id": produc_id, "pro_id": webjson["id"], "val": webjson["val"]})
	}
	DB().Table("client_product_manage_pro_val").Data(save_arr).Insert()
}

// 更新状态
func UpLock(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_product_manage").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

// 删除
func DelProduct(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("client_product_manage").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}

// 获取编辑内容
func GetProduct(context *gin.Context) {
	id := context.DefaultQuery("id", "")
	data, err := DB().Table("client_product_manage").Where("id", id).First()
	if err != nil {
		results.Failed(context, "获取大内容字段失败", err)
	} else {
		//分组
		cids, _ := DB().Table("client_product_cid").Where("product_id", id).Pluck("cid")
		data["cid"] = cids
		//标签
		lids, _ := DB().Table("client_product_lid").Where("product_id", id).Pluck("lid")
		data["lid"] = lids
		//多图
		var parameter []interface{}
		_ = json.Unmarshal([]byte(data["images"].(string)), &parameter)
		data["images"] = parameter
		results.Success(context, "获取编辑内容", data, nil)
	}
}

// 置顶
func UpTop(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	var top int64 = 0
	if parameter["checked"] == true {
		maxtop, _ := DB().Table("client_product_manage").Where("top", ">", 0).Order("top desc").Value("top")
		if maxtop != nil {
			top = maxtop.(int64) + 1
		} else {
			top = 1
		}
	}
	res2, err := DB().Table("client_product_manage").Where("id", parameter["id"]).Data(map[string]interface{}{"top": top}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, top, nil)
	}
}
