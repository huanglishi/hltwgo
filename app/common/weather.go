package common

import (
	"encoding/json"
	"fmt"
	"huling/utils/results"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取天气预报
func GetWeather(context *gin.Context) {
	var errData string = ""
	contentpoint := getlocalpoint()
	url := "https://api.map.baidu.com/weather/v1/?district_id=" + getDistrict(contentpoint) + "&data_type=now&ak=R2ZwnNw2cs2rWQhRI1DAshrTGBFTCTUY"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	//得到返回结果
	body, _ := ioutil.ReadAll(res.Body)
	bodystr := string(body)
	//对返回的json数据做解析
	var dataAttr map[string]interface{}
	if err := json.Unmarshal([]byte(bodystr), &dataAttr); err != nil {
		errData = "json解析出错"
		fmt.Println(err.Error())
	}
	if errData == "" {
		results.Success(context, "获取列表", dataAttr, nil)
	} else {
		results.Failed(context, errData, nil)
	}
}

// 获取当前城市名称
func getlocalpoint() string {
	var errData string = ""
	url := "https://api.map.baidu.com/location/ip?ak=R2ZwnNw2cs2rWQhRI1DAshrTGBFTCTUY&ip=" + get_external() + "&coor=bd09ll"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	//得到返回结果
	body, _ := ioutil.ReadAll(res.Body)
	bodystr := string(body)
	//对返回的json数据做解析
	var dataAttr map[string]interface{}
	var rebackAttr string = ""
	if err := json.Unmarshal([]byte(bodystr), &dataAttr); err == nil {
		for idx, value := range dataAttr {
			if idx == "content" {
				mapTmp := value.(map[string]interface{})
				if mapTmp["address_detail"] != nil {
					mapPoint := mapTmp["address_detail"].(map[string]interface{})
					rebackAttr = mapPoint["city"].(string)
				}
			}
		}
	} else {
		errData = "json解析出错:"
	}
	if errData == "" {
		return rebackAttr
	} else {
		return errData
	}
}

// 获取城市对应-行政区划编码
func getDistrict(name_str string) string {
	var errData string = ""
	url := "https://quhua.ipchaxun.com/api/areas/data?name=" + name_str
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	//得到返回结果
	body, _ := ioutil.ReadAll(res.Body)
	bodystr := string(body)
	//对返回的json数据做解析
	var dataAttr map[string]interface{}
	var rebackAttr string = ""
	if err := json.Unmarshal([]byte(bodystr), &dataAttr); err == nil {
		for idx, value := range dataAttr {
			if idx == "data" {
				mapTmp := value.(map[string]interface{})
				if mapTmp["results"] != nil {
					mapResults := mapTmp["results"].([]interface{})
					if mapResults != nil && len(mapResults) > 0 {
						mapCity := mapResults[0].(map[string]interface{})
						rebackAttr = mapCity["code"].(string)
					}
				}
			}
		}
	} else {
		errData = "json解析出错:"
	}
	if errData == "" {
		return rebackAttr
	} else {
		return errData
	}
}
func get_external() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}
