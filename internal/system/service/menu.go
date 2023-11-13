package service

import (
	"github.com/jinzhu/copier"
	"github.com/lucyanddarlin/lucy-ez-admin/constants"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
	"github.com/lucyanddarlin/lucy-ez-admin/tools/tree"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

// AllMenu 获取菜单树
func AllMenu(ctx *core.Context) (tree.Tree, error) {
	menu := model.Menu{}
	return menu.Tree(ctx)
}

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

// UpdateMenu 更新菜单
func UpdateMenu(ctx *core.Context, in *types.UpdateMenuRequest) error {
	inMenu := model.Menu{}

	if copier.Copy(&inMenu, in) != nil {
		return errors.AssignError
	}

	menu := model.Menu{}
	if err := menu.OneByID(ctx, in.ID); err != nil {
		return err
	}

	if in.ParentID != 0 && in.ID == in.ParentID {
		return errors.MenuParentIdError
	}

	// TODO: 测试
	if menu.Name != in.Name && menu.OneByName(ctx, in.Name) != nil {
		return errors.DulMenuNameError
	}

	// 之前为接口, 现在修改类型不为接口, 则删除之前的 rbac 权限
	if menu.Type == constants.MenuA && in.Type != constants.MenuA {
		_, _ = ctx.Enforcer().Instance().RemoveFilteredPolicy(1, menu.Path, menu.Method)
	}

	// 之前和现在都为接口, 且存在方法或者逻辑变更时, 更新 rbac 权限
	if menu.Type == constants.MenuA || in.Type == constants.MenuA && (menu.Method != in.Method || menu.Path != in.Path) {
		oldPolices := ctx.Enforcer().Instance().GetFilteredPolicy(1, menu.Method, menu.Path)
		if len(oldPolices) != 0 {
			var newPolices [][]string
			for _, val := range oldPolices {
				newPolices = append(newPolices, []string{val[0], in.Path, in.Method})
			}
			_, _ = ctx.Enforcer().Instance().UpdatePolicies(oldPolices, newPolices)
		}
	}

	// 之前不是接口, 现在修改为接口, 则进行 rbac 新增
	if menu.Type != constants.MenuA && in.Type == constants.MenuA {
		// 获取选中当前菜单的角色
		roleMenu := model.RoleMenu{}
		roleMenus, _ := roleMenu.MenuRoles(ctx, in.ID)
		if len(roleMenus) != 0 {
			var roleIds []int64
			for _, item := range roleMenus {
				roleIds = append(roleIds, item.RoleID)
			}

			// 获取当前菜单的全部角色信息
			role := model.Role{}
			roles, _ := role.All(ctx, "id in ?", roleIds)

			// 添加菜单到 rbac 权限表
			var newPolices [][]string
			for _, val := range roles {
				newPolices = append(newPolices, []string{val.Keyword, in.Path, in.Method})
			}
			_, _ = ctx.Enforcer().Instance().AddPolicies(newPolices)
		}
	}

	if err := inMenu.Update(ctx); err != nil {
		return err
	}

	// 更新首页菜单
	if inMenu.IsHome != menu.IsHome && inMenu.IsHome != nil && *inMenu.IsHome {
		return inMenu.UpdateMenuHome(ctx, inMenu.ID())
	}

	return nil
}

// DeleteMenu 删除菜单
func DeleteMenu(ctx *core.Context, in *types.DeleteMenuRequest) error {
	if in.ID == 1 {
		return errors.DeleteRootMenuError
	}

	// 获取指定 id 为根节点的菜单树
	menu := model.Menu{}
	list, _ := menu.All(ctx)
	var treeList []tree.Tree
	for _, item := range list {
		treeList = append(treeList, item)
	}
	t := tree.BuildTreeByID(treeList, in.ID)

	// 获取菜单树下所有的菜单 ID
	ids := tree.GetTreeID(t)

	// 删除当前 id 中的类型为 API 的 rbac 权限表
	apiList, _ := menu.All(ctx, "id in ? and type='A'", ids)
	for _, item := range apiList {
		ctx.Enforcer().Instance().RemoveFilteredPolicy(1, item.Path, item.Method)
	}

	// 从数据库删除菜单
	return menu.DeleteByIds(ctx, ids)
}
