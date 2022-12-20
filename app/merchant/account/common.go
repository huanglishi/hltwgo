package account

import (
	"huling/app/model"

	"github.com/gohouse/gorose/v2"
	jsoniter "github.com/json-iterator/go"
)

func DB() gorose.IOrm {
	return model.DB.NewOrm()
}

// 获取菜单子树结构
func GetMenuChildrenArray(pdata []gorose.Data, parent_id int64) []gorose.Data {
	var returnList []gorose.Data
	for _, v := range pdata {
		if v["pid"].(int64) == parent_id {
			children := GetMenuChildrenArray(pdata, v["id"].(int64))
			if children != nil {
				v["children"] = children
			}
			returnList = append(returnList, v)
		}
	}
	return returnList
}

// JSONMarshalToString JSON编码为字符串
func JSONMarshalToString(v interface{}) string {
	s, err := jsoniter.MarshalToString(v)
	if err != nil {
		return ""
	}
	return s
}

// StringToJSONMarshal 字符串为JSON编码
func StringToJSONMarshal(str string) interface{} {
	var s interface{}
	err := jsoniter.UnmarshalFromString(str, s)
	if err != nil {
		return nil
	}
	return s
}
