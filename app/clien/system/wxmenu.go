package system

import (
	"encoding/json"
	"fmt"
	"huling/utils/results"
	utils "huling/utils/tool"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取微信菜单-微信服务器上
func Getmenu(context *gin.Context) {
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	//获取公众号配置
	wxconfig, _ := DB().Table("client_system_wxconfig").Where("cuid", user.ClientID).Fields("id,name,AppID,AppSecret,expires_access_token,access_token").First()
	//更新access_token
	AccessTokenHost := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", wxconfig["AppID"], wxconfig["AppSecret"])
	timestamp := strconv.FormatInt(time.Now().Unix(), 10) //10位时间戳
	var (
		access_token  string
		wxAccessToken WxAccessToken
	)
	//当前时间戳转int
	intNum, _ := strconv.Atoi(timestamp)
	timestampint := int64(intNum)
	expires_access_token_int := wxconfig["expires_access_token"].(int64) //数据库的时间传戳
	//获取access_token，如果缓存中有，则直接取出数据使用；否则重新调用微信端接口获取
	client := &http.Client{}
	//判断access_token是否过期
	if wxconfig["access_token"] == "" || expires_access_token_int == 0 || (timestampint-expires_access_token_int) > 7200 { //重新请求access_token
		request, _ := http.NewRequest("GET", AccessTokenHost, nil)
		response, _ := client.Do(request)
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			results.Failed(context, "请求AccessToken失败1", err.Error())
			return
		}
		err = json.Unmarshal(body, &wxAccessToken)
		if err != nil {
			results.Failed(context, "解析AccessToken失败", err.Error())
			return
		}
		if wxAccessToken.Errcode == 0 {
			access_token = wxAccessToken.Access_token
		} else {
			results.Failed(context, "获取AccessToken失败", wxAccessToken.Errmsg)
			return
		}
		//添加access_tokens时间
		DB().Table("client_system_wxconfig").Where("id", wxconfig["id"]).Data(map[string]interface{}{"access_token": access_token, "expires_access_token": timestamp}).Update()
	} else {
		//缓存中存在access_token，直接读取
		access_token = wxconfig["access_token"].(string)
	}
	//获取 菜单接口
	wxmenu_data, err := Get_x(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token=%s", access_token))
	if err != nil {
		results.Failed(context, "获取微信openid失败", err)
	} else {
		var data_parameter map[string]interface{}
		if err := json.Unmarshal([]byte(wxmenu_data), &data_parameter); err == nil {
			results.Success(context, "获取微信菜单", map[string]interface{}{
				"name":   wxconfig["name"],
				"wxmenu": data_parameter,
			}, nil)
		}
	}

}

// 创建微信菜单
func SaveMenu(context *gin.Context) {
	//获取post传过来的data
	body, _ := ioutil.ReadAll(context.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	//当前用户
	getuser, _ := context.Get("user")
	user := getuser.(*utils.UserClaims)
	//获取公众号配置
	wxconfig, _ := DB().Table("client_system_wxconfig").Where("cuid", user.ClientID).Fields("id,name,AppID,AppSecret,expires_access_token,access_token").First()
	//更新access_token
	AccessTokenHost := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", wxconfig["AppID"], wxconfig["AppSecret"])
	timestamp := strconv.FormatInt(time.Now().Unix(), 10) //10位时间戳
	var (
		access_token  string
		wxAccessToken WxAccessToken
	)
	//当前时间戳转int
	intNum, _ := strconv.Atoi(timestamp)
	timestampint := int64(intNum)
	expires_access_token_int := wxconfig["expires_access_token"].(int64) //数据库的时间传戳
	//获取access_token，如果缓存中有，则直接取出数据使用；否则重新调用微信端接口获取
	client := &http.Client{}
	//判断access_token是否过期
	if wxconfig["access_token"] == "" || expires_access_token_int == 0 || (timestampint-expires_access_token_int) > 7200 { //重新请求access_token
		request, _ := http.NewRequest("GET", AccessTokenHost, nil)
		response, _ := client.Do(request)
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			results.Failed(context, "请求AccessToken失败1", err.Error())
			return
		}
		err = json.Unmarshal(body, &wxAccessToken)
		if err != nil {
			results.Failed(context, "解析AccessToken失败", err.Error())
			return
		}
		if wxAccessToken.Errcode == 0 {
			access_token = wxAccessToken.Access_token
		} else {
			results.Failed(context, "获取AccessToken失败", wxAccessToken.Errmsg)
			return
		}
		//添加access_tokens时间
		DB().Table("client_system_wxconfig").Where("id", wxconfig["id"]).Data(map[string]interface{}{"access_token": access_token, "expires_access_token": timestamp}).Update()
	} else {
		//缓存中存在access_token，直接读取
		access_token = wxconfig["access_token"].(string)
	}
	//获取 菜单接口
	wxmenu_data, err := Post(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s", access_token), parameter, "")
	if err != nil {
		results.Failed(context, "获取微信openid失败", err)
	} else {
		var data_parameter map[string]interface{}
		if err := json.Unmarshal([]byte(wxmenu_data), &data_parameter); err == nil {
			results.Success(context, "创建微信菜单", data_parameter, parameter)
		}
	}
}

// 保存微信菜单
func SaveMenuOnly(context *gin.Context) {
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
	if parameter["menu"] != nil {
		parameter["menu"] = JSONMarshalToString(parameter["menu"])
	}
	if f_id != 0 {
		res, err := DB().Table("client_system_wxmenu").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(context, "更新失败", err)
		} else {
			results.Success(context, "更新成功！", res, nil)
		}
	} else {
		parameter["cuid"] = user.ClientID
		parameter["accountID"] = user.Accountid
		addres, err := DB().Table("client_system_wxmenu").Data(parameter).Insert()
		if err != nil {
			results.Failed(context, "添加失败", err)
		} else {
			results.Success(context, "添加成功！", addres, nil)
		}
	}
}

type WxAccessToken struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
}
type WxJsApiTicket struct {
	Ticket     string `json:"ticket"`
	Expires_in int    `json:"expires_in"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}
type WxSignature struct {
	Noncestr  string `json:"noncestr"`
	Timestamp string `json:"timestamp"`
	Url       string `json:"url"`
	Signature string `json:"signature"`
	AppID     string `json:"appId"`
}

type WxSignRtn struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data WxSignature `json:"data"`
}
