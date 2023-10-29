package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"go.uber.org/zap"
)

type ResponseWriterWrapper struct {
	gin.ResponseWriter
	Body *bytes.Buffer // 缓存
}

func (w ResponseWriterWrapper) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriterWrapper) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func RequestLog() gin.HandlerFunc {
	white := core.GlobalConfig().Middleware.RequestLog.WhiteList

	return func(c *gin.Context) {
		ctx := core.New(c)
		defer ctx.Release()

		method := strings.ToLower(ctx.Request.Method)
		path := ctx.Request.URL.String()
		if strings.Contains(path, "?") {
			path = strings.Split(path, "?")[0]
		}

		// 判断白名单
		if white[method+":"+path] {
			return
		}

		now := time.Now()
		blw := ResponseWriterWrapper{Body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		req := RequestInfo(ctx.Gin())

		ctx.Next()

		end := time.Now()
		res := any(nil)
		if strings.Contains(ctx.Writer.Header().Get("Content-Type"), "application/json") {
			m := map[string]any{}
		_:
			json.Unmarshal(blw.Body.Bytes(), &m)
			res = m
		} else {
			res = blw.Body.String()
		}

		ctx.Logger().WithOptions(zap.WithCaller(false)).Info("request-info",
			zap.Any("path", ctx.Request.URL.Path),
			zap.Any("method", ctx.Request.Method),
			zap.Any("user_agent", ctx.Request.Header.Get("User-Agent")),
			zap.Any("timestamp", fmt.Sprintf("%vms", end.Sub(now).Milliseconds())),
			zap.Any("req", req),
			zap.Any("res", res),
			zap.Any("status", ctx.Writer.Status()),
		)
	}
}
