package service

import (
	"github.com/jinzhu/copier"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
	"github.com/lucyanddarlin/lucy-ez-admin/tools"
	"github.com/lucyanddarlin/lucy-ez-admin/tools/tree"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

// AllTeam 获取所有部门
func AllTeam(ctx *core.Context) (tree.Tree, error) {
	team := model.Team{}
	return team.Tree(ctx)
}

// AddTeam 添加部门
func AddTeam(ctx *core.Context, in *types.AddTeamRequest) error {
	team := model.Team{}
	team.OneByTeamName(ctx, in.Name)

	if team.Name != "" {
		return errors.ExistTeamError
	}

	team = model.Team{}

	// 获取用户管理的部门
	ids, err := CurrentAdminTeamIds(ctx)
	if err != nil {
		return err
	}

	if !tools.InList(ids, in.ParentID) {
		return errors.NotAddTeamError
	}

	if copier.Copy(&team, in) != nil {
		return errors.AssignError
	}

	return team.Create(ctx)
}
