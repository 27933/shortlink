package urltool

import (
	"errors"
	"net/url"
	"path"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetBasePath 解析短链接部分 -> 获取URL路径的最后一节(不包括query部分) -> 判断是否已经短链接转换过
// eg: www.27933.com/jv  -> jv
func GetBasePath(targetUrl string) (string, error) {
	// url.Parse()函数 -> 不会判断是否是url 他只是当从字符串处理(根据 / 进行分割)
	myUrl, err := url.Parse(targetUrl)
	if err != nil {
		logx.Errorw("url.Parse failed", logx.LogField{Key: "url", Value: targetUrl}, logx.LogField{Key: "err", Value: err.Error()})
		return "", err
	}
	if len(myUrl.Host) == 0 {
		return "", errors.New("no host in tartgetUrl")
	}
	return path.Base(myUrl.Path), nil
}
