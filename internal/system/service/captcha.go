package service

import (
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

func ImageCaptcha(ctx *core.Context, in *types.CaptchaRequest) (any, error) {
	return ctx.ImageCaptcha(in.Name).New()
}
