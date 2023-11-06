package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/service"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

// Captcha 发送图片验证码
func Captcha(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.CaptchaRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if resp, err := service.ImageCaptcha(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(resp)
	}
}

// EmailCaptcha 发送邮件验证码
func EmailCaptcha(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.CaptchaRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if resp, err := service.EmailCaptcha(ctx, &in); err != nil {
		ctx.RespError(err)
		return
	} else {
		ctx.RespData(resp)
	}
}
