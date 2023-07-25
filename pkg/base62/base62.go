package base62

import (
	"math"
	"strings"
)

// 62进制转换的模块

// 0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ

// 0-9: 0-9
// a-z: 10-35
// A-Z: 36-61

// 10进制数     ->     62进制数
// 	  0					 0
//    10				 a
//    61				 Z
//    62				 10
//    63				 11
//    6347				 ? -> 1En

// 代码如何实现62进制转换?

// 为了避免被人恶意请求可以将62进制对应的字符串给打乱
// const base62Str = `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
var (
	base62Str string
	// baseStrLen uint64
)

// MustInit 要使用base62这个包必须要完成该函数的初始化
func MustInit(bs string) {
	if len(bs) == 0 {
		panic("need base string!")
	}
	base62Str = bs
	// baseStrLen = uint64(len(bs))
}

// Int2Srting 输入10进制的数返回一个62进制的数
func Int2Srting(seq uint64) string {
	if seq == 0 {
		return string(base62Str[0])
	}
	bl := make([]byte, 0)
	for seq > 0 {
		mod := seq % 62
		div := seq / 62
		bl = append(bl, []byte(base62Str)[mod])
		seq = div
	}
	return string(reverse(bl))
}

func String2Int(s string) (res uint64) {
	bl := []byte(s)
	// 从左往右遍历
	for i := len(bl) - 1; i >= 0; i-- {
		base := strings.Index(base62Str, string(s[i]))
		res += uint64(math.Pow(62, float64(len(bl)-1-i)) * float64(base))
	}
	return
}

func reverse(s []byte) []byte {
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
	return s
}
