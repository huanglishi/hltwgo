package user

import (
	"huling/app/model"
	"strings"

	"github.com/gohouse/gorose/v2"
)

func DB() gorose.IOrm {
	return model.DB.NewOrm()
}

// 获取权限菜单
func GetMenuArray(pdata []gorose.Data, parent_id int64) []gorose.Data {
	var returnList []gorose.Data
	var one int64 = 1
	for _, v := range pdata {
		// fmt.Printf("ID: %v,值：%v\n", v["id"], parent_id)
		if v["parentMenu"].(int64) == parent_id {

			mid_item := map[string]interface{}{
				"path":      v["routePath"],
				"name":      v["routeName"],
				"component": v["component"],
			}
			children := GetMenuArray(pdata, v["id"].(int64))
			if children != nil {
				mid_item["children"] = children
			}
			//1.标题
			var Menu_title interface{}
			if v["title"] != nil && v["title"] != "" {
				Menu_title = v["title"]
			} else {
				Menu_title = v["menuName"]
			}
			meta := map[string]interface{}{
				"title": Menu_title,
			}
			//2.重定向
			if v["redirect"] != nil && v["redirect"] != "" {
				mid_item["redirect"] = v["redirect"]
			}
			//3.重定向
			if v["hideChildrenInMenu"] != nil && v["hideChildrenInMenu"] != "" {
				var hideChildrenInMenu bool
				if v["hideChildrenInMenu"].(int64) == one {
					hideChildrenInMenu = true
				} else {
					hideChildrenInMenu = false
				}
				meta["hideChildrenInMenu"] = hideChildrenInMenu
			}
			//3.图标
			if v["icon"] != nil && v["icon"] != "" {
				meta["icon"] = v["icon"]
			}
			//4.缓存
			if v["keepalive"] != nil && v["keepalive"].(int64) == one {
				meta["ignoreKeepAlive"] = false
			} else {
				meta["ignoreKeepAlive"] = true
			}
			//5.隐藏菜单
			if v["hideMenu"] != nil && v["hideMenu"].(int64) == one {
				meta["hideMenu"] = true
			}
			//6.在标签隐藏
			if v["hideTab"] != nil && v["hideTab"].(int64) == one {
				meta["hideTab"] = true
			}
			//7.详情页在本业打开-用于配置详情页时左侧激活的菜单路径
			if v["currentActiveMenu"] != nil && v["currentActiveMenu"] != "" {
				meta["currentActiveMenu"] = v["currentActiveMenu"]
			}
			//赋值
			mid_item["meta"] = meta
			returnList = append(returnList, mid_item)
		}
	}
	return returnList
}

// tool-获取树状数组
func GetTreeArray(num []gorose.Data, pid int64) []gorose.Data {
	childs := ToolFar(num, pid) //获取pid下的所有数据
	var chridnum []gorose.Data
	if childs != nil {
		for _, v := range childs {
			v["children"] = GetTreeArray(num, v["id"].(int64))
			chridnum = append(chridnum, v)
		}
	}
	return chridnum
}

// base_tool-获取pid下所有数组
func ToolFar(data []gorose.Data, pid int64) []gorose.Data {
	var mapString []gorose.Data
	for _, v := range data {
		if v["pid"].(int64) == pid {
			mapString = append(mapString, v)
		}
	}
	return mapString
}

// 判断元素是否存在数组中
func IsContain(items []interface{}, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// 多维数组合并
func ArrayMerge(data []interface{}) []interface{} {
	var rule_ids_arr []interface{}
	for _, mainv := range data {
		ids_arr := strings.Split(mainv.(string), `,`)
		for _, intv := range ids_arr {
			rule_ids_arr = append(rule_ids_arr, intv)
		}
	}
	return rule_ids_arr
}
