package webedit

//产品
import (
	"encoding/json"
	"fmt"
	"huling/utils/results"
	utils "huling/utils/tool"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gohouse/gorose/v2"
)

// 获列表
func GetProductList(context *gin.Context) {
	title := context.DefaultQuery("title", "")
	cid := context.DefaultQuery("cid", "")
	getall := context.DefaultQuery("getall", "0")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	whereMap := DB().Table("client_product_manage").Where("cuid", user.ClientID).Where("status", 0)
	whereMap2 := DB().Table("client_product_manage").Where("cuid", user.ClientID).Where("status", 0)
	if title != "" {
		whereMap.Where("title", "like", "%"+title+"%")
		whereMap2.Where("title", "like", "%"+title+"%")
	}
	if cid != "" {
		product_ids, _ := DB().Table("client_product_cid").Where("cid", cid).Pluck("product_id")
		whereMap.WhereIn("id", product_ids.([]interface{}))
		whereMap2.WhereIn("id", product_ids.([]interface{}))
	}
	if getall != "0" {
		whereMap.Limit(pageSize).Page(pageNo)
	}
	list, err := whereMap.Fields("id,type,title,des,images,createtime,releasetime").Order("top desc , weigh desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		prolist, _ := DB().Table("client_product_manage_pro").Where("cuid", user.ClientID).Where("status", 0).Fields("id,keyname,name,des,weigh,type").Order("weigh asc").Get()
		for _, val := range list {
			//分组
			cids, _ := DB().Table("client_product_cid").Where("product_id", val["id"]).Pluck("cid")
			catename, _ := DB().Table("client_product_cate").WhereIn("id", cids.([]interface{})).Pluck("name")
			//标签
			lids, _ := DB().Table("client_product_lid").Where("product_id", val["id"]).Pluck("lid")
			labelname, _ := DB().Table("client_product_label").WhereIn("id", lids.([]interface{})).Pluck("name")
			val["labelname"] = labelname
			val["catename"] = catename
			if val["images"] != "" {
				//多图
				var parameter []interface{}
				_ = json.Unmarshal([]byte(val["images"].(string)), &parameter)
				val["images"] = parameter
			} else {
				val["images"] = make([]interface{}, 0)
			}
			var myprolist []gorose.Data
			for _, pro := range prolist {
				pro_val, _ := DB().Table("client_product_manage_pro_val").Where("product_id", val["id"]).Where("pro_id", pro["id"]).Value("val")
				if pro_val != nil {
					pro["val"] = pro_val
				} else {
					pro["val"] = ""
				}
				myprolist = append(myprolist, map[string]interface{}{"id": pro["id"], "keyname": pro["keyname"], "name": pro["name"], "des": pro["des"], "weigh": pro["weigh"], "type": pro["type"], "val": pro["val"]})
			}
			val["prolist"] = myprolist
		}
		var totalCount int64
		totalCount, _ = whereMap2.Count()
		results.Success(context, "获取产品列表", map[string]interface{}{
			"items":    list,
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"getall":   getall,
		}, nil)
	}
}

// 获取详情内容
func GetProduct(context *gin.Context) {
	id := context.DefaultQuery("id", "")
	data, err := DB().Table("client_product_manage").Where("id", id).Fields("id,cuid,type,title,des,images,releasetime,content,createtime").First()
	if err != nil {
		results.Failed(context, "获取产品详情内容失败", err)
	} else {
		if data == nil {
			results.Success(context, "产品内容不存在！", nil, nil)
		} else {
			prolist, _ := DB().Table("client_product_manage_pro").Where("cuid", data["cuid"]).Where("status", 0).Fields("id,keyname,name,des,weigh,type").Order("weigh asc").Get()
			for _, pro := range prolist {
				pro_val, _ := DB().Table("client_product_manage_pro_val").Where("product_id", id).Where("pro_id", pro["id"]).Value("val")
				if pro_val != nil {
					pro["val"] = pro_val
				} else {
					pro["val"] = ""
				}
			}
			data["prolist"] = prolist
			//图片
			rooturl, _ := DB().Table("merchant_config").Where("keyname", "rooturl").Value("keyvalue")
			if data["images"] != "" {
				//多图
				var parameter []interface{}
				_ = json.Unmarshal([]byte(data["images"].(string)), &parameter)
				var newimg []interface{}
				for _, img := range parameter {
					img = fmt.Sprintf("%s%s", rooturl, img)
					newimg = append(newimg, img)
				}
				data["images"] = newimg
			} else {
				data["images"] = make([]interface{}, 0)
			}
			//标签
			lids, _ := DB().Table("client_product_lid").Where("product_id", id).Pluck("lid")
			labels, _ := DB().Table("client_product_label").WhereIn("id", lids.([]interface{})).Pluck("name")
			data["labels"] = labels
			results.Success(context, "获取产品详情内容", data, nil)
		}
	}
}

// 获取产品分类
func GetProductCate(context *gin.Context) {
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	data, err := DB().Table("client_product_cate").Where("cuid", user.ClientID).Fields("id,pid,name,weigh").Order("weigh asc , id asc").Get()
	if err != nil {
		results.Failed(context, "获取产品分类失败", err)
	} else {
		results.Success(context, "获取产品分类", data, nil)
	}
}
