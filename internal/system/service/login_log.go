package service

import (
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
	"github.com/lucyanddarlin/lucy-ez-admin/tools/address"
	"github.com/lucyanddarlin/lucy-ez-admin/tools/ua"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

func AddLoginLog(ctx *core.Context, phone string, err error) error {
	ip := ctx.ClientIP()
	userAgent := ctx.Request.Header.Get("User-Agent")
	info := ua.Parse(userAgent)
	desc := ""
	code := 0

	if err != nil {
		customErr, _ := err.(*types.Response)
		code = customErr.Code
		desc = customErr.Msg
	}

	log := model.LoginLog{
		Phone:       phone,
		IP:          ip,
		Address:     address.New(ip).GetAddress(),
		Browser:     info.Name,
		Status:      err == nil,
		Description: desc,
		Code:        code,
		Device:      info.OS + " " + info.OSVersion,
	}
	return log.Create(ctx)
}
