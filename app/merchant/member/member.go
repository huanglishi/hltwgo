package member

import (
	"encoding/json"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 1获取数据列表
func Getlist(context *gin.Context) {
	cid := context.DefaultQuery("cid", "0")
	title := context.DefaultQuery("title", "")
	page := context.DefaultQuery("page", "1")
	_pageSize := context.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(page)
	pageSize, _ := strconv.Atoi(_pageSize)
	MDB := DB().Table("client_user").Fields("id,cid,avatar,accountID,clientid,name,username,mobile,email,status,validtime,remark,company,createtime")
	MDBCON := DB().Table("client_user")
	if cid != "0" {
		sub_cids := GetAllChilId(cid)
		allsubids := append(sub_cids, cid)
		MDB = MDB.WhereIn("cid", allsubids)
		MDBCON = MDBCON.WhereIn("cid", allsubids)
	}
	if title != "" {
		MDB = MDB.Where("name", "like", "%"+title+"%")
		MDBCON = MDBCON.Where("name", "like", "%"+title+"%")
	}
	list, err := MDB.Limit(pageSize).Page(pageNo).Order("id desc").Get()
	if err != nil {
		results.Failed(context, "加载数据失败", err)
	} else {
		clien_url, _ := DB().Table("merchant_config").Where("keyname", "clien_url").Value("keyvalue")
		for _, val := range list {
			//分组
			groupname, _ := DB().Table("merchant_micweb_group").Where("id", val["cid"]).Value("name")
			val["groupname"] = groupname
			//获取轻站id
			micweb_data, _ := DB().Table("client_micweb").Where("cuid", val["id"]).Fields("id,status,title").First()
			//发布状态
			val["publish_status"] = micweb_data["status"]
			val["title"] = micweb_data["title"]
			val["webid"] = micweb_data["id"]
			//直接跳转地址
			if clien_url != nil {
				//拼接账号toke
				token := utils.GenerateToken(&utils.UserClaims{
					ID:             val["id"].(int64),
					Accountid:      val["accountID"].(int64),
					ClientID:       val["clientid"].(int64),
					StandardClaims: jwt.StandardClaims{},
				})
				val["webhref"] = clien_url.(string)
				val["token"] = token
			} else {
				val["webhref"] = ""
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

// 2添加用户
func SaveMemberData(context *gin.Context) {
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
	if parameter["validtime"] != nil {
		parameter["validtime"] = utils.StrToTime(parameter["validtime"].(string))
	} else {
		parameter["validtime"] = 0
	}
	if parameter["password"] != nil && parameter["password"] != "" {
		rnd := rand.New(rand.NewSource(6))
		salt := strconv.Itoa(rnd.Int())
		mdpass := parameter["password"].(string) + salt
		parameter["password"] = utils.Md5(mdpass)
		parameter["salt"] = salt
	}
	parameter["createtime"] = time.Now().Unix()
	if f_id == 0 {
		parameter["uid"] = user.ID
		parameter["accountID"] = user.Accountid
		parameter["avatar"] = "resource/staticfile/avatar.png"
		addId, err := DB().Table("client_user").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			//添加一条轻站数据
			micwebdata := map[string]interface{}{
				"cuid":       addId,
				"uid":        user.ID,
				"accountID":  user.Accountid,
				"is_select":  1,
				"title":      parameter["name"],
				"des":        parameter["remark"],
				"createtime": parameter["createtime"],
				"updatetime": parameter["createtime"],
			}
			DB().Table("client_micweb").Data(micwebdata).Insert()
			DB().Table("client_user").Data(map[string]interface{}{"clientID": addId}).Where("id", addId).Update()
			DB().Table("client_user_config").Data(map[string]interface{}{"accountID": user.Accountid, "cuid": addId, "fileSize": 5}).Insert()
			DB().Table("client_auth_role_access").Data(map[string]interface{}{"uid": addId, "role_id": 12}).Insert()
			results.Success(context, "添加成功！", addId, nil)
		}
	} else {
		res, err := DB().Table("client_user").
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

// 3检查账号是否存在
func IsAccountExist(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["id"] != nil {
		res1, err := DB().Table("client_user").Where("id", "!=", parameter["id"]).Where("username", parameter["account"]).Value("id")
		if err != nil {
			results.Failed(context, "验证失败", err)
		} else if res1 != nil {
			results.Failed(context, "账号已存在", err)
		} else {
			results.Success(context, "验证通过", res1, nil)
		}
	} else {
		res2, err := DB().Table("client_user").Where("username", parameter["account"]).Value("id")
		if err != nil {
			results.Failed(context, "验证失败", err)
		} else if res2 != nil {
			results.Failed(context, "账号已存在", err)
		} else {
			results.Success(context, "验证通过", res2, nil)
		}
	}
}

// 更新状态
func UpStatus(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := DB().Table("client_user").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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
