package mwebapproval

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

// 获取子菜单包含的父级ID-返回全部ID
func GetRulesID(menus interface{}) interface{} {
	menus_rang := menus.([]interface{})
	var fnemuid []interface{}
	for _, v := range menus_rang {
		fid := getParentID(v)
		if fid != nil {
			fnemuid = mergeArr(fnemuid, fid)
		}
	}
	r_nemu := mergeArr(menus_rang, fnemuid)
	uni_fnemuid := uniqueArr(r_nemu) //去重
	return uni_fnemuid
}

// 获取所有父级ID
func getParentID(id interface{}) []interface{} {
	var pids []interface{}
	pid, _ := DB().Table("merchant_auth_rule").Where("id", id).Value("parentMenu")
	if pid != nil {
		a_pid := pid.(int64)
		var zr_pid int64 = 0
		if a_pid != zr_pid {
			pids = append(pids, a_pid)
			getParentID(pid)
		}
	}
	return pids
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
