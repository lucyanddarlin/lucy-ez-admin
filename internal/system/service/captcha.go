package service

import (
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
	"go.uber.org/zap"
)

func ImageCaptcha(ctx *core.Context, in *types.CaptchaRequest) (any, error) {
	return ctx.ImageCaptcha(in.Name).New()
}

func EmailCaptcha(ctx *core.Context, in *types.CaptchaRequest) (any, error) {
	md := ctx.Metadata()

	if md == nil {
		return nil, errors.MetadataError
	}

	// 获取用户邮箱信息
	user := model.User{}
	if err := user.OneByID(ctx, md.UserID); err != nil {
		return nil, err
	}

	// 发送邮箱验证码
	res, err := ctx.EmailCaptcha(in.Name).New(user.Email)
	if err != nil {
		ctx.Logger().Error("验证码发送失败", zap.Error(err))
		return nil, errors.CaptchaSendError
	}

	return res, nil
}
