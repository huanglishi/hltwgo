package packagedesign

import (
	"encoding/json"
	"fmt"
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
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := DB().Table("merchant_packagedesign")
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("id asc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		for _, val := range list {
			paylist, _ := DB().Table("merchant_packagedesign_paylist").Where("packg_id", val["id"]).Get()
			val["paylist"] = paylist
		}
		var totalCount int64
		totalCount, _ = DB().Table("merchant_packagedesign").Count()
		results.Success(context, "获取全部列表", map[string]interface{}{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}, nil)
	}
}

// 获取套餐数据
func UseList(context *gin.Context) {
	userid := context.DefaultQuery("userid", "0")
	menuList, _ := DB().Table("merchant_packagedesign").Fields("id,name,image,type,masterckeck,price,count").Order("id asc").Get()
	for _, val := range menuList {
		packages, _ := DB().Table("merchant_user_package").Where("user_id", userid).Where("pid", val["id"]).Value("number")
		if packages != nil {
			val["checks"] = 1
			val["number"] = packages
		} else {
			val["checks"] = 0
			val["number"] = 0
		}
	}
	results.Success(context, "获取套餐数据", menuList, nil)
}

// 添加套餐
func SaveData(context *gin.Context) {
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
	if parameter["menu"] != nil {
		rudata := parameter["menu"].([]interface{})
		var rulesStr []string
		for _, v := range rudata {
			str := fmt.Sprintf("%v", v) //interface{}强转string
			rulesStr = append(rulesStr, str)
		}
		parameter["menu_str"] = strings.Join(rulesStr, ",")
		parameter["menu"] = JSONMarshalToString(parameter["menu"])

	}
	parameter["createtime"] = time.Now().Unix()
	paylist := parameter["paylist"]
	//删除数组
	delete(parameter, "paylist")
	if f_id == 0 {
		addId, err := DB().Table("merchant_packagedesign").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			addPayPackge(paylist.([]interface{}), addId)
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("merchant_packagedesign").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			addPayPackge(paylist.([]interface{}), f_id)
			results.Success(context, "更新成功！", res, nil)
		}
	}
}

// 批量支付预置套餐
func addPayPackge(plist []interface{}, packg_id interface{}) {
	DB().Table("merchant_packagedesign_paylist").Where("packg_id", packg_id).Delete()
	//批量提交
	save_arr := []map[string]interface{}{}
	for _, val := range plist {
		marr := val.(map[string]interface{})
		marr["packg_id"] = packg_id
		save_arr = append(save_arr, marr)
	}
	DB().Table("merchant_packagedesign_paylist").Data(save_arr).Insert()
}

// 删除
func Del(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := DB().Table("merchant_packagedesign").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(context, "删除失败", err)
	} else {
		results.Success(context, "删除成功！", res2, nil)
	}
	context.Abort()
	return
}
