package service

import (
	"github.com/lucyanddarlin/lucy-ez-admin/constants"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
	"github.com/lucyanddarlin/lucy-ez-admin/tools/tree"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

// RoleMenuIds 获取角色菜单的所有 id
func RoleMenuIds(ctx *core.Context, in *types.RoleMenuIdsRequest) ([]int64, error) {
	// 获取当前角色的所有菜单
	rm := model.RoleMenu{}
	rmList, err := rm.RoleMenus(ctx, in.RoleID)
	if err != nil {
		return nil, err
	}

	// 组装所有的菜单 id
	var ids []int64
	for _, item := range rmList {
		ids = append(ids, item.MenuID)
	}

	return ids, nil
}

func RoleMenu(ctx *core.Context, in *types.RoleMenuRequest) (tree.Tree, error) {
	// 查询角色信息
	role := model.Role{}
	if err := role.OneByID(ctx, in.RoleID); err != nil {
		return nil, err
	}

	var menus []*model.Menu
	var menu model.Menu

	if role.Keyword == constants.JwtSuperAdmin {
		menus, _ = menu.All(ctx, "type!=?", constants.MenuBA)
	} else {
		// 查询角色所属菜单
		rm := model.RoleMenu{}
		rmList, err := rm.RoleMenus(ctx, in.RoleID)
		if err != nil {
			return nil, err
		}

		// 获取菜单的所有 id
		var ids []int64
		for _, item := range rmList {
			ids = append(ids, item.MenuID)
		}
		if len(ids) == 0 {
			return nil, nil
		}

		// 获取指定 id 的所有菜单
		menus, _ = menu.All(ctx, "id in ? and type!= ?", ids, constants.MenuBA)
	}

	var listTree []tree.Tree
	for _, item := range menus {
		listTree = append(listTree, item)
	}

	return tree.BuildTree(listTree), nil
}

// UpdateRoleMenu 修改角色所属菜单
func UpdateRoleMenu(ctx *core.Context, in *types.AddRoleMenuRequest) error {
	// 超级管理员不存在菜单权限, 自动获取所有菜单, 因此禁止修改
	if in.RoleID == 1 {
		return errors.SuperAdminEditError
	}

	// 获取当前 role 的数据
	role := model.Role{}
	if err := role.OneByID(ctx, in.RoleID); err != nil {
		return err
	}

	// 进行菜单修改
	rm := model.RoleMenu{}
	if err := rm.Update(ctx, in.RoleID, in.MenuIDs); err != nil {
		return nil
	}

	// 删除当前用户的全部 rbac 权限
	_, _ = ctx.Enforcer().Instance().RemoveFilteredPolicy(0, role.Keyword)

	// 获取当前修改菜单的信息
	menu := model.Menu{}
	var policies [][]string
	apiList, _ := menu.All(ctx, "id in ? and type = A", in.MenuIDs)
	for _, item := range apiList {
		policies = append(policies, []string{role.Keyword, item.Path, item.Method})
	}

	// 将新的策略的策略写入 rbac
	_, _ = ctx.Enforcer().Instance().AddPolicies(policies)

	return nil
}
