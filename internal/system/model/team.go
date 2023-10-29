package model

import (
	"time"

	"github.com/lucyanddarlin/lucy-ez-admin/core"
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

func (t *Team) TableName() string {
	return "tb_system_team"
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
