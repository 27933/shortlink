package sequence

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// 基于Redis实现一个发号器

type Redis struct {
	conn redis.Redis
}

func NewRedis(dsn string) Sequence {
	return &Redis{
		conn: *redis.New(dsn),
	}
}

func (c *Redis) Next() (seq uint64, err error) {
	var lid int64
	lid, err = c.conn.Incr("stub")
	if err != nil {
		logx.Errorw("rest.LastInsertId failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}
	return uint64(lid), nil
}
