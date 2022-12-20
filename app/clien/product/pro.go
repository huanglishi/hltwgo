package product

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// 获取参数列表
func GetproList(context *gin.Context) {
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	num, _ := DB().Table("client_product_manage_pro").Where("cuid", user.ID).Count()
	var zero int64 = 0
	if num == zero { //如果是第一次加载则初始数据
		weighId, _ := DB().Table("client_product_manage_pro").Data(map[string]interface{}{
			"cuid": user.ClientID, "accountID": user.Accountid, "keyname": "no", "name": "编号", "des": "产品的编号", "type": 2}).InsertGetId()
		if weighId != 0 {
			DB().Table("client_product_manage_pro").
				Data(map[string]interface{}{"weigh": weighId}).
				Where("id", weighId).
				Update()
			save_arr := []map[string]interface{}{}

			save_arr = append(save_arr, map[string]interface{}{
				"cuid": user.ClientID, "accountID": user.Accountid, "keyname": "type", "name": "类型", "des": "产品的类型", "type": 2, "weigh": weighId + 1})
			save_arr = append(save_arr, map[string]interface{}{
				"cuid": user.ClientID, "accountID": user.Accountid, "keyname": "spec", "name": "规格", "des": "产品的规格", "type": 2, "weigh": weighId + 2})
			save_arr = append(save_arr, map[string]interface{}{
				"cuid": user.ClientID, "accountID": user.Accountid, "keyname": "material", "name": "材质", "des": "产品的材质", "type": 2, "weigh": weighId + 3})
			save_arr = append(save_arr, map[string]interface{}{
				"cuid": user.ClientID, "accountID": user.Accountid, "keyname": "color", "name": "颜色", "des": "产品的颜色", "type": 2, "weigh": weighId + 4})
			save_arr = append(save_arr, map[string]interface{}{
				"cuid": user.ClientID, "accountID": user.Accountid, "keyname": "brand", "name": "品牌", "des": "产品的品牌", "type": 2, "weigh": weighId + 5})
			save_arr = append(save_arr, map[string]interface{}{
				"cuid": user.ClientID, "accountID": user.Accountid, "keyname": "model", "name": "型号", "des": "产品的型号", "type": 2, "weigh": weighId + 6})
			save_arr = append(save_arr, map[string]interface{}{
				"cuid": user.ClientID, "accountID": user.Accountid, "keyname": "mprice", "name": "市场价", "des": "市场价格", "type": 1, "weigh": weighId + 7})
			save_arr = append(save_arr, map[string]interface{}{
				"cuid": user.ClientID, "accountID": user.Accountid, "keyname": "price", "name": "价格", "des": "实际交易价", "type": 1, "weigh": weighId + 8})
			save_arr = append(save_arr, map[string]interface{}{
				"cuid": user.ClientID, "accountID": user.Accountid, "keyname": "sales", "name": "销量", "des": "产品的销量", "type": 1, "weigh": weighId + 9})
			DB().Table("client_product_manage_pro").Data(save_arr).Insert()
		}
	}
	from := context.DefaultQuery("from", "")
	MDB := DB().Table("client_product_manage_pro").Where("cuid", user.ClientID)
	if from == "product" {
		MDB.Where("status", 0)
	}
	list, _ := MDB.Order("weigh asc").Get()
	for _, val := range list {
		if from == "pro" {
			val["edit"] = false
		}
		if from == "product" {
			product_id := context.DefaultQuery("product_id", "")
			if product_id != "" {
				pro_val, _ := DB().Table("client_product_manage_pro_val").Where("product_id", product_id).Where("pro_id", val["id"]).Value("val")
				if pro_val != nil {
					val["val"] = pro_val
				} else {
					val["val"] = ""
				}
			} else {
				val["val"] = ""
			}
			//获取配置值
			if val["type"] != 1 {
				vallist, _ := DB().Table("client_product_manage_pro_list").Where("pro_id", val["id"]).Where("status", 0).Fields("id,value").Order("weigh asc").Get()
				val["vallist"] = vallist
			}
		}
	}
	results.Success(context, "产品参数数据！", list, nil)
}

// 保存/修改数据
func SavePro(context *gin.Context) {
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
	delete(parameter, "id")
	delete(parameter, "edit")
	if f_id == 0 {
		parameter["cuid"] = user.ClientID
		parameter["accountID"] = user.Accountid
		addId, err := DB().Table("client_product_manage_pro").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败！", err)
		} else {
			if addId != 0 {
				DB().Table("client_product_manage_pro").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		_, err := DB().Table("client_product_manage_pro").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			results.Success(context, "更新成功！", f_id, user)
		}
	}
}

// 删除
func DelPro(context *gin.Context) {
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("client_product_manage_pro").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
}

// 更新状态
func UpPro(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_product_manage_pro").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

// 更新排序
func UpWeigh(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res1, err := DB().Table("client_product_manage_pro").Where("id", parameter["id"]).Data(map[string]interface{}{"weigh": parameter["rpweigh"]}).Update()
	if err != nil {
		results.Failed(context, "更新失败！", err)
	} else {
		DB().Table("client_product_manage_pro").Where("id", parameter["rpId"]).Data(map[string]interface{}{"weigh": parameter["weigh"]}).Update()
		msg := "更新成功！"
		if res1 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(context, msg, res1, nil)
	}
}
