package webedit

//客服
import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// 获取列表
func GetServiceList(context *gin.Context) {
	//当前用户
	product_id := context.DefaultQuery("product_id", "0")
	MDB := DB().Table("client_product_service")
	list, err := MDB.Where("product_id", product_id).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "获取客服列表失败", err)
	} else {
		results.Success(context, "获取客服列表", list, nil)
	}
}

// 获取单条
func GetService(context *gin.Context) {
	id := context.DefaultQuery("id", "0")
	list, err := DB().Table("client_product_service").Where("id", id).First()
	if err != nil {
		results.Failed(context, "获取客服失败", err)
	} else {
		results.Success(context, "获取客服数据", list, nil)
	}
}

// 添加
func SaveService(context *gin.Context) {
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
	if f_id == 0 {
		parameter["cuid"] = user.ClientID
		parameter["accountID"] = user.Accountid
		addId, err := DB().Table("client_product_service").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("client_product_service").
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

// 删除
func DelService(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_product_service").Where("id", parameter["id"]).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
}
