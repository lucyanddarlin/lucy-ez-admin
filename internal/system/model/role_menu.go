package model

import (
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
	"gorm.io/gorm"
)

type RoleMenu struct {
	types.BaseModel
	RoleID     int64  `json:"role_id" gorm:"not null;size:32;comment:角色id"`
	MenuID     int64  `json:"menu_id" gorm:"not null;size:32;comment:菜单id"`
	Operator   string `json:"operator" gorm:"not null;size:128;comment:操作人名称"`
	OperatorID int64  `json:"operator_id" gorm:"not null;size:32;comment:操作人id"`
	Role       Role   `json:"role" gorm:"->;constraint:onDelete:cascade"`
	Menu       Menu   `json:"menu" gorm:"->;constraint:onDelete:cascade"`
}

func (*RoleMenu) TableName() string {
	return "tb_system_role_menu"
}

// Update 批量更新角色所属菜单
func (rm *RoleMenu) Update(ctx *core.Context, roleId int64, menuIds []int64) error {
	// 操作者信息
	md := ctx.Metadata()
	if md == nil {
		return errors.MetadataError
	}

	rm.Operator = md.Username
	rm.OperatorID = md.UserID

	// 组装新的菜单数据
	list := make([]RoleMenu, 0)
	for _, menuId := range menuIds {
		list = append(list, RoleMenu{
			RoleID:     roleId,
			MenuID:     menuId,
			Operator:   md.Username,
			OperatorID: md.UserID,
		})
	}

	// 删除之后再重新创建
	err := database(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleId).Delete(rm).Error; err != nil {
			return nil
		}
		return tx.Create(&list).Error
	})

	return transferErr(err)
}

// RoleMenus 通过角色 id 获取角色菜单
func (rm *RoleMenu) RoleMenus(ctx *core.Context, roleId int64) ([]*RoleMenu, error) {
	var list []*RoleMenu

	return list, transferErr(database(ctx).Find(&list, "role_id = ?", roleId).Error)
}

// MenuRoles 通过菜单 ID 获取角色菜单列表
func (rm *RoleMenu) MenuRoles(ctx *core.Context, menuId int64) ([]*RoleMenu, error) {
	var list []*RoleMenu

	return list, transferErr(database(ctx).Find(&list, "menu_id = ?", menuId).Error)
}

// DeleteByRoleID 通过角色 ID 删除角色所属菜单
func (rm *RoleMenu) DeleteByRoleID(ctx *core.Context, roleId int64) error {
	return transferErr(database(ctx).Delete(rm, "role_id = ?", roleId).Error)
}

func (rm *RoleMenu) InitData(ctx *core.Context) error {
	return nil
}
