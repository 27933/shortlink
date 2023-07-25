package logic

import (
	"context"
	"database/sql"
	"errors"

	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	Err404 = errors.New("404")
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 自己写缓存 		  surl -> lurl
// go-zero自带的缓存  surl -> 数据行  缓存整个数据行

func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	// 查看短链接 输入27933.cn/jt -> 重定向到真实的连接
	// req.ShortUrl = jt
	// 1. 根据短链接查询原始的长链接
	// 	1.0 布隆过滤器
	//		不存在的短链接直接返回404, 不需要后续处理
	//		1.0.1. 基于内存版本
	//			缺点: 服务重启之后所有短链接数据都没了, 所以每次重启重启都要重新加载一下已有的短链接
	//		1.0.2 基于redis版本
	//			相比于基于内存版本 服务重启后短链接数据依然存在 只要Redis没有挂就数据就一直存在
	exist, err := l.svcCtx.Filter.Exists([]byte(req.ShortUrl))
	if err != nil {
		logx.Errorw("l.svcCtx.Filter.Exists Bloom Filter failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	// 布隆过滤器中判断数据是否存在
	if !exist {
		return nil, Err404
	}
	// go-zero 自带缓存原生的支持singleflight
	u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("404")
		}
		logx.Errorw("l.svcCtx.ShortUrlModel.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	// 2. 返回重定向响应 在调用Show函数处返回重定向响应
	return &types.ShowResponse{LongUrl: u.Lurl.String}, nil
}
