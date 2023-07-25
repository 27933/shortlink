package md5

import (
	"crypto/md5"
	"encoding/hex"
)

// Sum 根据传入的数据参数生成MD5值
func Sum(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil)) // 32位16进制数
}
