package model

import (
	"time"

	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/tools/tree"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

type Team struct {
	types.BaseModel
	Name        string  `json:"name" gorm:"not null;size:128;comment:部门名称"`
	Description string  `json:"description,omitempty" gorm:"size:256;comment:部门备注"`
	ParentID    int64   `json:"parent_id" gorm:"not null;size:32;comment:父级部门"`
	Operator    string  `json:"operator" gorm:"not null;size:128;comment:操作人员名称"`
	OperatorID  int64   `json:"operator_id" gorm:"not null;size:32;comment:操作人员id"`
	Children    []*Team `json:"children,omitempty" gorm:"-"`
}

func (t *Team) ID() int64 {
	return t.BaseModel.ID
}

func (t *Team) Parent() int64 {
	return t.ParentID
}

func (t *Team) AppendChildren(child any) {
	team := child.(*Team)
	t.Children = append(t.Children, team)
}

func (t *Team) ChildrenNode() []tree.Tree {
	var list []tree.Tree
	for _, item := range t.Children {
		list = append(list, item)
	}
	return list
}

func (t *Team) TableName() string {
	return "tb_system_team"
}

// All 获取所有部门
func (t *Team) All(ctx *core.Context) ([]*Team, error) {
	list := make([]*Team, 0)
	if err := database(ctx).Find(&list).Error; err != nil {
		return nil, transferErr(err)
	}
	return list, nil
}

// Tree 获取部门树
func (t *Team) Tree(ctx *core.Context) (tree.Tree, error) {
	// 获取部门列表
	list := make([]*Team, 0)
	if err := database(ctx).Find(&list).Error; err != nil {
		return nil, err
	}

	// 根据列表构建部门树
	var trees []tree.Tree
	for _, item := range list {
		trees = append(trees, item)
	}
	return tree.BuildTree(trees), nil
}

// OneByTeamName 根据部门名称查询部门
func (t *Team) OneByTeamName(ctx *core.Context, name string) error {
	return transferErr(database(ctx).First(t, "name = ?", name).Error)
}

// IsExistTeamByName
func (t *Team) IsExistTeamByName(ctx *core.Context, name string) {
	t.OneByTeamName(ctx, name)
}

// Create 添加部门
func (t *Team) Create(ctx *core.Context) error {
	md := ctx.Metadata()
	if md == nil {
		return errors.MetadataError
	}
	t.Operator = md.Username
	t.OperatorID = md.UserID

	return transferErr(database(ctx).Create(&t).Error)
}

// Update 更新部门
func (t *Team) Update(ctx *core.Context) error {
	md := ctx.Metadata()
	if md == nil {
		return errors.MetadataError
	}

	t.Operator = md.Username
	t.OperatorID = md.UserID

	return transferErr(database(ctx).Updates(t).Error)
}

// DeleteByID 通过 ID 删除指定部门, 以及该部门下的所有部门
func (t *Team) DeleteByID(ctx *core.Context, id int64) error {
	list, err := t.All(ctx)
	if err != nil {
		return err
	}

	var treeList []tree.Tree
	for _, item := range list {
		treeList = append(treeList, item)
	}
	team := tree.BuildTreeByID(treeList, id)
	ids := tree.GetTreeID(team)

	// 进行数据删除
	return transferErr(database(ctx).Where("id in ?", ids).Delete(&t).Error)
}

func (t *Team) InitData(ctx *core.Context) error {
	db := database(ctx)
	ins := Team{
		BaseModel:   types.BaseModel{ID: 1, CreatedAt: time.Now().Unix()},
		Name:        "青橙科技有限责任公司",
		ParentID:    0,
		Description: "青橙科技有限责任公司",
	}
	return db.Create(&ins).Error
}
