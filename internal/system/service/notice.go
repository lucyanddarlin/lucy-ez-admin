package service

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
	"github.com/lucyanddarlin/lucy-ez-admin/tools/proto"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
	"gorm.io/gorm"
)

// AddNotice 添加通知
func AddNotice(ctx *core.Context, in *types.AddNoticeRequest) error {
	notice := model.Notice{
		Status: proto.Bool(false),
	}

	if copier.Copy(&notice, in) != nil {
		return errors.AssignError
	}

	return notice.Create(ctx)
}

// UpdateNotice 更新通知
func UpdateNotice(ctx *core.Context, in *types.UpdateNoticeRequest) error {
	notice := model.Notice{}

	if copier.Copy(&notice, in) != nil {
		return errors.AssignError
	}

	return notice.Update(ctx)
}

// DeleteNotice 删除通知
func DeleteNotice(ctx *core.Context, in *types.DeleteNoticeRequest) error {
	notice := model.Notice{}

	return notice.DeleteByID(ctx, in.ID)
}

// PageNotice 分页获取通知列表
func PageNotice(ctx *core.Context, in *types.PageNoticeRequest) ([]*model.Notice, int64, error) {
	md := ctx.Metadata()
	if md == nil {
		return nil, 0, errors.MetadataError
	}

	notice := model.Notice{}
	noticeUser := model.NoticeUser{}

	// 返回数据
	return notice.Page(ctx, types.PageOptions{
		Page:     in.Page,
		PageSize: in.PageSize,
		Model:    in,
		Scopes: func(db *gorm.DB) *gorm.DB {
			join := fmt.Sprintf("left join %s on %s.notice_id=%s.id and %s.user_id=%d",
				noticeUser.TableName(),
				noticeUser.TableName(),
				notice.TableName(),
				noticeUser.TableName(),
				md.UserID,
			)
			db = db.Joins(join)

			if in.IsRead == nil {

				return db
			}
			if *in.IsRead {
				return db.Where(fmt.Sprintf("%s.user_id is not null", noticeUser.TableName()))
			} else {
				return db.Where(fmt.Sprintf("%s.user_id is null", noticeUser.TableName()))
			}
		},
	})
}

// GetNotice 获取通知信息
func GetNotice(ctx *core.Context, in *types.GetNoticeRequest) (*model.Notice, error) {
	notice := &model.Notice{}
	return notice, notice.OneByID(ctx, in.ID)
}

// GetNoticeUnreadNum 获取未读通知数量
func GetNoticeUnReadNum(ctx *core.Context) (int64, error) {
	notice := model.Notice{}
	return notice.UnreadNum(ctx)
}
