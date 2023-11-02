package model

import (
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/tools"
	"gorm.io/gorm"
)

const (
	_orm = "system"
)

func DBName() string {
	return _orm
}

// DataMap 数据字典
var DataMap = map[string]string{
	"phone":   "手机号码",
	"email":   "电子邮箱",
	"keyword": "标志",
	"name":    "名称",
}

func database(ctx *core.Context) *gorm.DB {
	return ctx.Orm().GetDB(_orm).WithContext(ctx.SourceCtx())
}
func transferErr(err error) error {
	return tools.TransferErr(DataMap, err)
}
