package wxpay

import (
	"encoding/json"
	"fmt"
	"huling/app/model"
	"math"
	"math/rand"
	"time"

	"github.com/gohouse/gorose/v2"
	jsoniter "github.com/json-iterator/go"
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

// JSONMarshalToString JSON编码为字符串
func JSONMarshalToString(v interface{}) string {
	s, err := jsoniter.MarshalToString(v)
	if err != nil {
		return ""
	}
	return s
}

// 获取时间部分
func GetFormatTime(time time.Time) string {
	return time.Format("20060102")
}

// 获取
func GetTimeTick64() int64 {
	return time.Now().UnixNano() / 1e6
}

// 获取编号
func GenerateCode() string {
	date := GetFormatTime(time.Now())
	r := rand.Intn(1000)
	return fmt.Sprintf("%s%d%03d", date, GetTimeTick64(), r)
}

// interface转float64
func Interface2Type(i interface{}) float64 {
	var t2 float64
	switch i.(type) {
	case float64:
		fmt.Println(i.(float64))
	default:
		fmt.Println(0)
	}
	return t2
}

// 将float64转成精确的int64
func Wrap(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}
