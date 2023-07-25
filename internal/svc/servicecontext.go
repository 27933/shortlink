package svc

import (
	"shortener/internal/config"
	"shortener/sequence"

	"shortener/model"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// svc层 -> 数据库相关操作 -> 数据库缓存、依赖等数据库服务

type ServiceContext struct {
	Config        config.Config
	ShortUrlModel model.ShortUrlMapModel //对应short_url_map这张表

	Sequence          sequence.Sequence //对应的是sequence表
	ShortUrlBlackList map[string]struct{}

	// bloom filter
	Filter *bloom.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 1. 初始化MySQL连接
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	// 2. 初始化短链接黑名单
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	// 把配置文件中配置的黑名单加载到map, 方便后续判断
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}
	// 3. 初始化布隆过滤器
	// 初始化RedisBitSet
	store := redis.New(c.BloomRedis.Host, func(r *redis.Redis) {
		r.Type = redis.NodeType
	})
	// 声明一个bitSet key="test_key"名且bits是1024位
	filter := bloom.New(store, "test_Key", 20*(1<<20))
	return &ServiceContext{
		Config:        c,
		ShortUrlModel: model.NewShortUrlMapModel(conn, c.CacheRedis),
		// Sequence:      sequence.NewMySQL(c.Sequence.DSN),
		Sequence:          sequence.NewMySQL(c.Sequence.DSN),
		ShortUrlBlackList: m,
		Filter:            filter,
	}
}
