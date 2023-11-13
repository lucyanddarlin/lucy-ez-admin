package tools

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
	"gorm.io/gorm"
)

// TransferErr 将数据的错误转换成中文
func TransferErr(m map[string]string, err error) error {
	if err == nil {
		return nil
	}

	if customErr, ok := err.(*types.Response); ok {
		return customErr
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.DBNotFoundError
	}

	if strings.Contains(err.Error(), "Duplicate") {
		str := err.Error()
		str = strings.ReplaceAll(str, "'", "")
		str = strings.TrimPrefix(str, "Error 1062: Duplicate Key ")
		arr := strings.Split(str, " for key ")
		return errors.NewF(`%v "%v" 已存在`, m[arr[1]], arr[0])
	}

	if strings.Contains(err.Error(), "FOREIGN KEY") {
		return errors.NewF(`数据正在被使用中, 无法删除`)
	}

	return errors.DBError
}

// DelayDelCache 数据延迟双删
func DelayDelCache(redis *redis.Client, key string) {
	redis.Del(context.Background(), key)
	go func() {
		time.Sleep(1 * time.Second)
		redis.Del(context.Background(), key)
	}()
}
