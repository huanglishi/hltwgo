package mwebh5

import (
	"encoding/json"
	"huling/app/model"

	"github.com/gohouse/gorose/v2"
)

func DB() gorose.IOrm {
	return model.DB.NewOrm()
}

// 字符串转JSON编码
func StingToJSON(str interface{}) []interface{} {
	var parameter []interface{}
	_ = json.Unmarshal([]byte(str.(string)), &parameter)
	return parameter
}
