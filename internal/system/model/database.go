package model

import (
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"gorm.io/gorm"
)

const (
	_orm = "system"
)

func DBName() string {
	return _orm
}

func database(ctx *core.Context) *gorm.DB {
	return ctx.Orm().GetDB(_orm).WithContext(ctx.SourceCtx())
}
