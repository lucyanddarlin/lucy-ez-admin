package model

import "github.com/lucyanddarlin/lucy-ez-admin/core"

type NoticeUser struct {
	NoticeID int64  `json:"notice_id" gorm:"not null;size:32;comment:通知 id"`
	UserID   int64  `json:"user_id" gorm:"not null;size:32;comment:人员 id"`
	ReadAt   int64  `json:"read_at" gorm:"not null;size:32;comment:阅读时间"`
	User     User   `gorm:"->"`
	Notice   Notice `gorm:"->"`
}

func (nu *NoticeUser) TableName() string {
	return "tb_system_notice_user"
}

// Create 创建阅读消息
func (nu *NoticeUser) Create(ctx *core.Context) error {
	return transferErr(database(ctx).Create(nu).Error)
}

func (nu *NoticeUser) InitData(ctx *core.Context) error {
	return nil
}
