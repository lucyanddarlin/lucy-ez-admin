package model

import (
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

type LoginLog struct {
	types.CreateModel
	Phone       string `json:"phone" gorm:"not null;type:varbinary(32);comment:手机号"`
	IP          string `json:"ip" gorm:"not null;type:varbinary(64);comment:登陆IP"`
	Address     string `json:"address" gorm:"not null;size:128;comment:登陆地址"`
	Browser     string `json:"browser" gorm:"not null;size:128;comment:登陆浏览器"`
	Device      string `json:"device" gorm:"not null;size:128;comment:登录设备"`
	Status      bool   `json:"status" gorm:"not null;comment:登录状态"`
	Code        int    `json:"code" gorm:"not null;size:32;comment:错误码"`
	Description string `json:"description" gorm:"not null;size:256;comment:登录备注"`
}

func (u *LoginLog) TableName() string {
	return "tb_system_login_log"
}

func (u *LoginLog) Create(ctx *core.Context) error {
	return transferErr(database(ctx).Create(u).Error)
}

func (u *LoginLog) InitData(ctx *core.Context) error {
	return nil
}
