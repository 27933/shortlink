package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/model"
	"shortener/pkg/base62"
	"shortener/pkg/connect"
	"shortener/pkg/md5"
	"shortener/pkg/urltool"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// logic层 -> 业务逻辑相关的操作

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Convert 转链 -> 输入一个长链接转为一个短链接
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 1. 解析请求 获取长链接 校验长链接数据
	// 	1.1 获取长链接
	// 		long := req.LongUrl
	// 	1.2 校验
	// 		1.2.1 不能为空
	// 			方法一:  常规方法校验
	//				if len(req.LongUrl) == 0 {}
	// 			方法二:  使用validator库来做参数校验 -> 在Handler已经注册了validator服务, 故在这里不需要再进行校验
	// 		1.2.2 长链接是一个有效连接 -> 能请求通的网站
	// 				if r, err := http.Get(req.LongUrl); err != nil{ ... } -> 封装成包 实现复用
	if ok := connect.Get(req.LongUrl); !ok {
		// 说明是无效链接 网络不可达
		return nil, errors.New("req中长链接是无效链接")
	}
	// 		1.2.3 判断之前是否已经转链 -> 查询数据库中是否已经存在该长链接
	// 			1.2.3.1 使用生成md5的方法给长链接缩小长度 -> 不能在数据库中存储太长的值
	md5Value := md5.Sum([]byte(req.LongUrl))
	// 			1.2.3.2 判断数据库中是否已经存在该长链接
	u, err := l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, fmt.Errorf("该链接已被转为%s", u.Surl.String)
		}
		logx.Errorw("l.svcCtx.ShortUrlModel.FindOneByMd5 failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	// 		1.2.4 避免循环转链 -> 输入的不能是一个短链接 -> 查询数据库是否存在相同短链接
	// 			27933.cn/lydad  -> shortUrl=lydad
	basePath, err := urltool.GetBasePath(req.LongUrl)
	if err != nil {
		logx.Errorw("urltool.GetBasePath failed", logx.LogField{Key: "url", Value: req.LongUrl}, logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	_, err = l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, fmt.Errorf("该链接已是短链了")
		}
		logx.Errorw("l.svcCtx.ShortUrlModel.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	var short string
	for {
		// 2. 使用转号器 基于MySQL实现的发号器   取号 -> 转链
		// 	每来一个转链请求, 就使用REPLACE INTO语句往 sequence 表插入一条数据, 取出主键id作为号码
		var seq uint64
		seq, err = l.svcCtx.Sequence.Next()
		if err != nil {
			logx.Errorw("l.svcCtx.Sequence.Next failed", logx.LogField{Key: "err", Value: err.Error()})
			return nil, err
		}
		// 使用seq数据完成转链操作
		fmt.Printf("---------seq:%v\n", seq)
		// 3. 存储数据库 -> 存储长链接短链接的映射关系
		//		3.1 安全性
		//			打乱62进制每一位对应的数
		short = base62.Int2Srting(seq)
		fmt.Println("short:", short)
		//		3.2 短域名避免某些特殊词 -> 黑名单机制 health、api、fuck等等
		//			使用配置文件定义某些词不能出现
		if _, ok := l.svcCtx.ShortUrlBlackList[short]; !ok {
			// 生成的短链接没在黑名单里面就跳出循环
			break
		}
	}
	// 4. 存储长短链接映射
	if _, err := l.svcCtx.ShortUrlModel.Insert(
		l.ctx,
		&model.ShortUrlMap{
			Lurl: sql.NullString{String: req.LongUrl, Valid: true},
			Md5:  sql.NullString{String: md5Value, Valid: true},
			Surl: sql.NullString{String: short, Valid: true},
		}); err != nil {
		logx.Errorw("l.svcCtx.ShortUrlModel.Insert failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	//		4.1 将生成的短链接加到布隆过滤器中 -> 避免出现缓存穿透的问题
	if err := l.svcCtx.Filter.Add([]byte(short)); err != nil {
		logx.Errorw("l.svcCtx.Filter.Add Bloom Add failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	// 5. 返回响应 返回的是 短域名+短链接
	shortUrl := l.svcCtx.Config.ShortDoamin + "/" + short
	return &types.ConvertResponse{ShortUrl: shortUrl}, nil
}
