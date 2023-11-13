package service

import (
	"github.com/jinzhu/copier"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

// AddMenu 添加菜单
func AddMenu(ctx *core.Context, in *types.AddMenuRequest) error {
	menu := model.Menu{}
	if copier.Copy(&menu, in) != nil {
		return errors.AssignError
	}

	if in.Name != "" && menu.OneByName(ctx, in.Name) == nil {
		return errors.DulMenuNameError
	}

	if err := menu.Create(ctx); err != nil {
		return err
	}

	// 更新菜单首页
	if in.IsHome {
		return menu.UpdateMenuHome(ctx, menu.ID())
	}

	return nil
}
