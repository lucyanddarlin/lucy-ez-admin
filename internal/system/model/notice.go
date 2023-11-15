package model

import (
	"time"

	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

type Notice struct {
	types.BaseModel
	Type       string `json:"type" gorm:"not null;size:128;comment:通知类型"`
	Title      string `json:"title" gorm:"not null;size:256;comment:通知标题"`
	Status     *bool  `json:"status" gorm:"not null;comment:通知状态"`
	Content    string `json:"content,omitempty" gorm:"not null;comment:通知内容"`
	Operator   string `json:"operator" gorm:"not null;size:128;comment:操作人名称"`
	OperatorID int64  `json:"operator_id" gorm:"not null;size:64;comment:操作人 ID"`
	ReadAt     int64  `json:"read_at" gorm:"-"`
}

func (n *Notice) TableName() string {
	return "tb_system_notice"
}

// Create 创建通知
func (n *Notice) Create(ctx *core.Context) error {
	md := ctx.Metadata()
	if md == nil {
		return errors.MetadataError
	}

	n.Operator = md.Username
	n.OperatorID = md.UserID
	n.UpdatedAt = time.Now().Unix()

	return transferErr(database(ctx).Create(n).Error)
}

// Update 更新通知
func (n *Notice) Update(ctx *core.Context) error {
	md := ctx.Metadata()
	if md == nil {
		return errors.MetadataError
	}

	return transferErr(database(ctx).Updates(n).Error)
}

func (n *Notice) InitData(ctx *core.Context) error {
	return nil
}
