package model

import (
	"time"

	"github.com/lucyanddarlin/lucy-ez-admin/constants"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
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
