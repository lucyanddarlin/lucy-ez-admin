package service

import (
	"github.com/jinzhu/copier"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
	"github.com/lucyanddarlin/lucy-ez-admin/tools/proto"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
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
