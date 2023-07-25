package handler

import (
	"net/http"

	"shortener/internal/logic"
	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// ConvertHandler 处理来自路由转发的长链接请求 并返回生成的短链接响应
func ConvertHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req types.ConvertRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// validator 参数规则校验 -> 解析传入的请求是否事先定义的validator tag要求
		if err := validator.New().StructExceptCtx(r.Context(), &req); err != nil {
			logx.Errorw("validator check failed", logx.LogField{Key: "err", Value: err.Error()})
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 执行业务逻辑
		// logic.NewConvertLogic 获取调用逻辑的结构体
		l := logic.NewConvertLogic(r.Context(), svcCtx)
		// 执行获取逻辑结构体实现的方法 -> 实现逻辑处理方法的调用 -> 长链转短链
		resp, err := l.Convert(&req)
		// 根据返回结构进行判断 -> 返回短链
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
