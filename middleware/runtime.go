package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strings"
	"unsafe"

	"github.com/gin-gonic/gin"
)

// RequestInfo 获取请求参数
func RequestInfo(c *gin.Context) map[string]any {
	// 获取 Body
	getBody := func(c *gin.Context) any {
		// 处理 form/data 上传
		if strings.Contains(c.Request.Header.Get("Context-Type"), "multipart/form-data") {
			return "getFormData"
		}

		//读取数据
		data, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))

		// 转换格式
		if strings.Contains(c.Request.Header.Get("Context-Type"), "application/json") {
			m := map[string]any{}
			_ = json.Unmarshal(data, &m)
			return m
		}

		// 返回字符串
		return *(*string)(unsafe.Pointer(&data))
	}

	getParam := func(c *gin.Context) string {
		return c.Request.URL.Query().Encode()
	}

	return map[string]any{
		"params": getParam(c),
		"body":   getBody(c),
	}
}

// PanicErr 获取 panic 错误堆栈
func PanicErr() []string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	var arr []string
	for _, pc := range pcs[:n-4] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		arr = append(arr, fmt.Sprintf("%s:%d", file, line))
	}
	return arr
}
