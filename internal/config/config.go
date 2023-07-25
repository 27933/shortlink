package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	//连接转发mysql -> 转链的表 -> 存储长短链接映射
	ShortUrlDB struct {
		DSN string
	}

	//连接发号器mysql -> 存储发送器的表
	Sequence struct {
		DSN string
	}

	// 定义根据发号器生成的id来生成短链接的基础信息 base62指定基础字符串
	BaseString string

	// 定义生成的短链接中的敏感词汇 并且放入黑名单
	ShortUrlBlackList []string
	// 定义基本的短域名 因为最后响应的是 短域名+短链接的方式
	ShortDoamin string

	CacheRedis cache.CacheConf
	/* cache.CacheConf 结构体中的字段
	type (
		// A ClusterConf is the config of a redis cluster that used as cache.
		ClusterConf []NodeConf

		// A NodeConf is the config of a redis node that used as cache.
		NodeConf struct {
			redis.RedisConf
			Weight int `json:",default=100"`
		}
	)
	*/

	BloomRedis redis.RedisConf
}

// 不能加这个 会导致生成ctx出现问题
// type MySQLConfig struct {
// }
