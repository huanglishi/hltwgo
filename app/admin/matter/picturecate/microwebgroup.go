package picturecate

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取数据列表
func Getlist(context *gin.Context) {
	keyword_name := context.DefaultQuery("name", "")
	keyword_status := context.DefaultQuery("status", "")
	typeid := context.DefaultQuery("type", "3")
	MDB := DB().Table("client_picture_cate")
	if keyword_name != "" {
		MDB = MDB.Where("name", "like", "%"+keyword_name+"%")
	}
	if keyword_status != "" {
		MDB = MDB.Where("status", keyword_status)
	}
	if typeid != "3" {
		MDB = MDB.Where("type", typeid)
	}
	menuList, _ := MDB.Order("weigh asc").Get()
	results.Success(context, "获取数据列表", menuList, nil)
}

// 获取选择分组
func GetCateTree(context *gin.Context) {
	keyword := context.DefaultQuery("keyword", "")
	MDB := DB().Table("client_picture_cate")
	if keyword != "" {
		MDB = MDB.Where("name", "like", "%"+keyword+"%")
	}
	menuList, _ := MDB.Fields("id,name,type").Order("weigh asc").Get()
	results.Success(context, "获取数据列表", menuList, nil)
}

// 添加
func Add(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	parameter["uid"] = user.ID
	var f_id float64 = 0
	if parameter["id"] != nil {
		f_id = parameter["id"].(float64)
	}
	if f_id == 0 {
		parameter["createtime"] = time.Now().Unix()
		addId, err := DB().Table("client_picture_cate").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			if addId != 0 {
				DB().Table("client_picture_cate").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("client_picture_cate").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			results.Success(context, "更新成功！", res, nil)
		}
	}
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
	res2, err := DB().Table("client_picture_cate").WhereIn("id", ids_arr).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

// 更新父级
func UpGrouppid(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	b_ids, _ := json.Marshal(parameter["ids"])
	var ids_arr []interface{}
	json.Unmarshal([]byte(b_ids), &ids_arr)
	res2, err := DB().Table("client_picture_cate").WhereIn("id", ids_arr).Data(map[string]interface{}{"pid": parameter["pid"]}).Update()
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
func Del(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("client_picture_cate").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
