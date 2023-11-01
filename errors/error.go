package errors

import (
	"errors"
	"fmt"

	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

const (
	code = 100001 // 通用错误码
)

var (
	New = func(msg string) error {
		return &types.Response{
			Code: code,
			Msg:  msg,
		}
	}

	NewF = func(msg string, arg ...interface{}) error {
		return &types.Response{
			Code: code,
			Msg:  fmt.Sprintf(msg, arg...),
		}
	}

	Is = func(err, tar error) bool {
		return errors.Is(err, tar)
	}
)
