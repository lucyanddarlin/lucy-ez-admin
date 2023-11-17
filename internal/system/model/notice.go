package model

import (
	"fmt"
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

// DeleteByID 根据 ID 删除通知
func (n *Notice) DeleteByID(ctx *core.Context, id int64) error {
	return transferErr(database(ctx).Delete(n, id).Error)
}

// Page 查询分页数据
func (n *Notice) Page(ctx *core.Context, options types.PageOptions) ([]*Notice, int64, error) {
	list, total := make([]*Notice, 0), int64(0)

	db := database(ctx).Model(n)

	if options.Model != nil {
		db = ctx.Orm().GormWhere(db, n.TableName(), options.Model)
	}

	if options.Scopes != nil {
		db = db.Scopes(options.Scopes)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Select("id,type,title,status,operator,operator_id,created_at,updated_at,read_at")
	db = db.Offset((options.Page - 1) * options.PageSize).Limit(options.PageSize)

	return list, total, db.Find(&list).Error
}

// OneByID 根据 ID 获取通知信息
func (n *Notice) OneByID(ctx *core.Context, id int64) error {
	md := ctx.Metadata()
	if md == nil {
		return errors.MetadataError
	}

	nu := NoticeUser{
		UserID:   md.UserID,
		NoticeID: id,
		ReadAt:   time.Now().Unix(),
	}

	_ = nu.Create(ctx)

	return transferErr(database(ctx).Find(n, id).Error)
}

// UnreadNum 获取未读通知数量
func (n *Notice) UnreadNum(ctx *core.Context) (int64, error) {
	md := ctx.Metadata()
	if md == nil {
		return 0, errors.MetadataError
	}

	nu := NoticeUser{}
	db := database(ctx).Model(n)

	join := fmt.Sprintf("left join %s on %s.notice_id=%s.id and %s.user_id=%d",
		nu.TableName(),
		nu.TableName(),
		n.TableName(),
		nu.TableName(),
		md.UserID,
	)

	total := int64(0)
	return total, db.Joins(join).Where(fmt.Sprintf("status=true and %v.user_id is null", nu.TableName())).Count(&total).Error
}

func (n *Notice) InitData(ctx *core.Context) error {
	return nil
}
