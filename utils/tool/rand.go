package utils

import (
	"math/rand"

	"github.com/bwmarrin/snowflake"
)

var (
	chars = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// RandString
func RandString(l int) string {
	bs := []byte{}
	for i := 0; i < l; i++ {
		bs = append(bs, chars[rand.Intn(len(chars))])
	}
	return string(bs)
}

// 获取UUID 生成雪花ID
func Getuuid() int64 {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return 0
	}
	return node.Generate().Int64()
}
