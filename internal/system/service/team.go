package service

import (
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
	"github.com/lucyanddarlin/lucy-ez-admin/tools/tree"
)

func AllTeam(ctx *core.Context) (tree.Tree, error) {
	team := model.Team{}
	return team.Tree(ctx)
}
