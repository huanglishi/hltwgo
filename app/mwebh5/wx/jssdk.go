package wx

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"huling/utils/results"
	"io"
	"io/ioutil"
	mathRand "math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取微信签名-前端分享sdk
func GetWxSign(context *gin.Context) {
	micweb_id := context.DefaultQuery("micweb_id", "")
	cuid, _ := DB().Table("client_micweb").Where("id", micweb_id).Value("cuid")
	var wxconfig map[string]interface{}
	system_wxconfig_table := "client_system_wxconfig"
	cwxconfig, _ := DB().Table("client_system_wxconfig").Where("cuid", cuid).Fields("id,AppID,AppSecret,expires_access_token,expires_jsapi_ticket,jsapi_ticket,access_token").First()
	if cwxconfig == nil { //C端没有账号可以调用平台
		adminwxconfig, _ := DB().Table("admin_system_wxconfig").Where("cuid", cuid).Fields("id,AppID,AppSecret,expires_access_token,expires_jsapi_ticket,jsapi_ticket,access_token").First()
		wxconfig = adminwxconfig
		system_wxconfig_table = "admin_system_wxconfig"
	} else {
		wxconfig = cwxconfig
	}
	if wxconfig == nil {
		results.Failed(context, "您的账号为配置微信公众号API及秘钥", nil)
		return
	}
	var (
		AppID           string = wxconfig["AppID"].(string)
		AppSecret       string = wxconfig["AppSecret"].(string)
		AccessTokenHost string = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + AppID + "&secret=" + AppSecret
		JsAPITicketHost string = "https://api.weixin.qq.com/cgi-bin/ticket/getticket"
	)

	var (
		noncestr, jsapi_ticket, timestamp, url, signature, signatureStr, access_token string
		wxAccessToken                                                                 WxAccessToken
		wxJsApiTicket                                                                 WxJsApiTicket
		wxSignature                                                                   WxSignature
		wxSignRtn                                                                     WxSignRtn
	)
	url = context.DefaultQuery("url", "")
	noncestr = RandStringBytes(16)
	timestamp = strconv.FormatInt(time.Now().Unix(), 10) //10位时间戳
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
			wxSignRtn.Code = 1
			wxSignRtn.Msg = err.Error()
			results.Failed(context, "请求AccessToken失败1", wxSignRtn)
			return
		}
		err = json.Unmarshal(body, &wxAccessToken)
		if err != nil {
			wxSignRtn.Code = 1
			wxSignRtn.Msg = err.Error()
			results.Failed(context, "解析AccessToken失败", wxSignRtn)
			return
		}
		if wxAccessToken.Errcode == 0 {
			access_token = wxAccessToken.Access_token
		} else {
			wxSignRtn.Code = 1
			wxSignRtn.Msg = wxAccessToken.Errmsg
			results.Failed(context, "获取AccessToken失败", wxSignRtn)
			return
		}
		//添加access_tokens时间
		DB().Table(system_wxconfig_table).Where("id", wxconfig["id"]).Data(map[string]interface{}{"access_token": access_token, "expires_access_token": timestamp}).Update()

		//获取 jsapi_ticket
		requestJs, _ := http.NewRequest("GET", JsAPITicketHost+"?access_token="+access_token+"&type=jsapi", nil)
		responseJs, _ := client.Do(requestJs)
		defer responseJs.Body.Close()
		bodyJs, err := ioutil.ReadAll(responseJs.Body)
		if err != nil {
			wxSignRtn.Code = 1
			wxSignRtn.Msg = err.Error()
			results.Failed(context, "请求jsapi_ticket失败", wxSignRtn)
			return
		}
		err = json.Unmarshal(bodyJs, &wxJsApiTicket)
		if err != nil {
			wxSignRtn.Code = 1
			wxSignRtn.Msg = err.Error()
			results.Failed(context, "解析jsapi_ticket失败", wxSignRtn)
			return
		}
		if wxJsApiTicket.Errcode == 0 {
			jsapi_ticket = wxJsApiTicket.Ticket
		} else {
			wxSignRtn.Code = 1
			wxSignRtn.Msg = wxJsApiTicket.Errmsg
			results.Failed(context, "获取jsapi_ticket失败", wxSignRtn)
			return
		}
		//更新数据库jsapi_ticket时间
		DB().Table(system_wxconfig_table).Where("id", wxconfig["id"]).Data(map[string]interface{}{"jsapi_ticket": jsapi_ticket, "expires_jsapi_ticket": timestamp}).Update()
	} else {
		//缓存中存在access_token，直接读取
		access_token = wxconfig["access_token"].(string)
		jsapi_ticket = wxconfig["jsapi_ticket"].(string)
	}

	// 获取 signature
	signatureStr = "jsapi_ticket=" + jsapi_ticket + "&noncestr=" + noncestr + "&timestamp=" + timestamp + "&url=" + url
	signature = GetSha1(signatureStr)

	wxSignature.Url = url
	wxSignature.Noncestr = noncestr
	wxSignature.Timestamp = timestamp
	wxSignature.Signature = signature
	wxSignature.AppID = AppID

	// 返回前端需要的数据
	wxSignRtn.Code = 0
	wxSignRtn.Msg = "success"
	wxSignRtn.Data = wxSignature
	results.Success(context, "签名", wxSignRtn, nil)
}

// 生成指定长度的字符串
func RandStringBytes(n int) string {
	const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[mathRand.Intn(len(letterBytes))]
	}
	return string(b)
}

// SHA1加密
func GetSha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
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
