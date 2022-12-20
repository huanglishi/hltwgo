package microweb

import (
	"huling/app/model"
	"strconv"

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

// 　合并数组
func mergeArr(a, b []interface{}) []interface{} {
	var arr []interface{}
	for _, i := range a {
		arr = append(arr, i)
	}
	for _, j := range b {
		arr = append(arr, j)
	}
	return arr
}

// 去重
func uniqueArr(m []interface{}) []interface{} {
	d := make([]interface{}, 0)
	tempMap := make(map[int]bool, len(m))
	for _, v := range m { // 以值作为键名
		keyv := GetInterfaceToInt(v)
		if tempMap[keyv] == false {
			tempMap[keyv] = true
			d = append(d, v)
		}
	}
	return d
}

// interface{}转int
func GetInterfaceToInt(t1 interface{}) int {
	var t2 int
	switch t1.(type) {
	case uint:
		t2 = int(t1.(uint))
		break
	case int8:
		t2 = int(t1.(int8))
		break
	case uint8:
		t2 = int(t1.(uint8))
		break
	case int16:
		t2 = int(t1.(int16))
		break
	case uint16:
		t2 = int(t1.(uint16))
		break
	case int32:
		t2 = int(t1.(int32))
		break
	case uint32:
		t2 = int(t1.(uint32))
		break
	case int64:
		t2 = int(t1.(int64))
		break
	case uint64:
		t2 = int(t1.(uint64))
		break
	case float32:
		t2 = int(t1.(float32))
		break
	case float64:
		t2 = int(t1.(float64))
		break
	case string:
		t2, _ = strconv.Atoi(t1.(string))
		break
	default:
		t2 = t1.(int)
		break
	}
	return t2
}

// interface{}转float64
func InterfaceToFloat64(t1 interface{}) float64 {
	var t2 float64
	switch t1.(type) {
	case float64:
		t2 = t1.(float64)
		break
	case string:
		t2, _ = strconv.ParseFloat(t1.(string), 64)
		break
	default:
		t2 = t1.(float64)
		break
	}
	return t2
}
