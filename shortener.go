package main

import (
	"flag"
	"fmt"

	"shortener/internal/config"
	"shortener/internal/handler"
	"shortener/internal/svc"
	"shortener/pkg/base62"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

// 在编译和运行的时候 通过指定配置文件所在位置运行
// eg: go run -f ./xx/xxx.yaml -> 指定配置xxx.yaml配置文件
var configFile = flag.String("f", "etc/shortener-api.yaml", "the config file")

// 主程序入口 使用go-zero框架开发
func main() {
	// 解析命令行请求
	flag.Parse()

	// 配置文件加载
	// 默认是加载./etc/shortener-api.yaml 文件
	var c config.Config
	conf.MustLoad(*configFile, &c)
	fmt.Printf("--------load conf::\n%#v\n", c)

	//初始化base62模块
	base62.MustInit(c.BaseString)

	// 创建一个server
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	/* 初始化自定义的接口或者结构体和一些配置数据 用于在程序后续运行中方便的调用
	 eg: type ServiceContext struct {
			Config        config.Config
			// 初始化ShortUrlModel 后可以方便的调用其定义的函数
			// eg: l.svcCtx.ShortUrlModel.FindOneBySurl(...)
			ShortUrlModel model.ShortUrlMapModel //对应short_url_map这张表

			Sequence          sequence.Sequence //对应的是sequence表
			ShortUrlBlackList map[string]struct{}

			// bloom filter
			// 初始化布隆过滤器 后续可以直接调用
			Filter *bloom.Filter
		}
	*/
	ctx := svc.NewServiceContext(c)
	// 将Handler注册进入server中 即路由和Handler相匹配
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
