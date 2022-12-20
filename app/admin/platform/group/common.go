package platformgroup

import (
	"huling/app/model"

	"github.com/gohouse/gorose/v2"
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
