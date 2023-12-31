package model

import (
	"time"

	"github.com/lucyanddarlin/lucy-ez-admin/constants"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/tools/tree"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
	"google.golang.org/protobuf/proto"
)

type Role struct {
	types.BaseModel
	ParentID    int64   `json:"parent_id" gorm:"not null;size:32;comment:父角色id"`
	Name        string  `json:"name" gorm:"not null;size:64;comment:角色名称"`
	Keyword     string  `json:"keyword" gorm:"not null;type:varbinary(32);comment:角色关键字"`
	Status      *bool   `json:"status,omitempty" gorm:"not null;comment:角色状态"`
	Weight      *int    `json:"weight" gorm:"default:0;size:16;comment:角色权重"`
	Description *string `json:"description,omitempty" gorm:"size:128;comment:角色备注"`
	TeamIds     *string `json:"team_ids,omitempty" gorm:"type:text;comment:自定义权限部门id"`
	DataScope   string  `json:"data_scope,omitempty" gorm:"not null;size:128;comment:数据权限"`
	Operator    string  `json:"operator" gorm:"size:128;comment:操作人员名称"`
	OperatorID  int64   `json:"operator_id" gorm:"size:32;comment:操作人员id"`
	Children    []*Role `json:"children"  gorm:"-"`
}

func (r *Role) TableName() string {
	return "tb_system_role"
}

func (r *Role) ID() int64 {
	return r.BaseModel.ID
}

func (r *Role) Parent() int64 {
	return r.ParentID
}

func (r *Role) AppendChildren(child any) {
	menu := child.(*Role)
	r.Children = append(r.Children, menu)
}

func (r *Role) ChildrenNode() []tree.Tree {
	var list []tree.Tree
	for _, item := range r.Children {
		list = append(list, item)
	}
	return list
}

// OneByID 通过 ID 查询角色信息
func (r *Role) OneByID(ctx *core.Context, id int64) error {
	return transferErr(database(ctx).First(r, "id = ?", id).Error)
}

// RoleStatus 获取角色状态
func (r *Role) RoleStatus(ctx *core.Context, roleID int64) bool {
	team, err := r.Tree(ctx, 1)
	if err != nil {
		return false
	}
	res := false
	dfsRoleStatus(team.(*Role), roleID, true, &res)
	return res
}

func dfsRoleStatus(role *Role, roleID int64, status bool, res *bool) bool {
	if roleID == role.BaseModel.ID {
		is := *role.Status && status
		*res = is
	}
	for _, item := range role.Children {
		dfsRoleStatus(item, roleID, status && *item.Status, res)
	}
	return status
}

// Tree 查询指定角色为根节点角色树
func (r *Role) Tree(ctx *core.Context, roleID int64) (tree.Tree, error) {
	list, err := r.All(ctx)
	if err != nil {
		return nil, err
	}
	var treeList []tree.Tree
	for _, item := range list {
		treeList = append(treeList, item)
	}
	return tree.BuildTreeByID(treeList, roleID), nil
}

// All 获取所有的角色
func (r *Role) All(ctx *core.Context, cond ...any) ([]*Role, error) {
	var list []*Role
	return list, transferErr(database(ctx).Order("weight desc").Find(&list, cond...).Error)
}

// Create 创建角色
func (r *Role) Create(ctx *core.Context) error {
	md := ctx.Metadata()
	if md == nil {
		return errors.MetadataError
	}

	r.Operator = md.Username
	r.OperatorID = md.UserID

	return transferErr(database(ctx).Create(&r).Error)
}

// Update 更新角色
func (r *Role) Update(ctx *core.Context) error {

	md := ctx.Metadata()
	if md == nil {
		return errors.MetadataError
	}

	r.Operator = md.Username
	r.OperatorID = md.UserID

	return transferErr(database(ctx).Updates(&r).Error)
}

// DeleteByID 根据 roleID 删除角色
func (r *Role) DeleteByID(ctx *core.Context, id int64) error {
	return transferErr(database(ctx).Where("id = ?", id).Delete(&r).Error)
}

func (r *Role) InitData(ctx *core.Context) error {
	db := database(ctx)
	ins := []Role{
		{
			BaseModel: types.BaseModel{
				ID:        1,
				CreatedAt: time.Now().Unix(),
			}, ParentID: 0, Name: "超级管理员", Keyword: "superAdmin", Status: proto.Bool(true), DataScope: constants.ALLTEAM,
		},
	}
	return db.Create(&ins).Error
}
