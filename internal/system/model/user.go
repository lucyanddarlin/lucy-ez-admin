package model

import (
	"time"

	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/tools"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
	"google.golang.org/protobuf/proto"
)

type User struct {
	types.BaseModel
	TeamID      int64   `json:"team_id" gorm:"not null;size:32;comment:部门id"`
	RoleID      int64   `json:"role_id" gorm:"not null;size:32;comment:角色id"`
	Name        string  `json:"name" gorm:"not null;size:32;comment:用户姓名"`
	Nickname    string  `json:"nickname" gorm:"not null;size:128;comment:用户昵称"`
	Sex         *bool   `json:"sex,omitempty" gorm:"not null;comment:用户性别"`
	Phone       string  `json:"phone" gorm:"not null;size:32;comment:用户电话"`
	Password    string  `json:"password,omitempty" gorm:"not null;->:false;<-:create,update;comment:用户密码"`
	Avatar      string  `json:"avatar" gorm:"not null;size:128;comment:用户头像"`
	Email       string  `json:"email,omitempty" gorm:"not null;type:varbinary(128);comment:用户邮箱"`
	Status      *bool   `json:"status,omitempty" gorm:"not null;comment:用户状态"`
	DisableDesc *string `json:"disable_desc" gorm:"not null;size:128;comment:禁用原因"`
	LastLogin   int64   `json:"last_login" gorm:"comment:最后登陆时间"`
	Operator    string  `json:"operator" gorm:"not null;size:128;comment:操作人员名称"`
	OperatorID  int64   `json:"operator_id" gorm:"not null;size:32;comment:操作人员id"`
	Role        Role    `json:"role" gorm:"->"`
	Team        Team    `json:"team" gorm:"->"`
}

func (u *User) TableName() string {
	return "tb_system_user"
}

// OnePyPhone 通过 phone 查询用户信息
func (u *User) OneByPhone(ctx *core.Context, phone string) error {
	db := database(ctx).Preload("Role").Preload("Team")
	return transferErr(db.First(u, "phone=?", phone).Error)
}

// PasswordByPhone 查询全部字段信息包括密码
func (u *User) PasswordByPhone(ctx *core.Context, phone string) (string, error) {
	m := map[string]any{}
	if err := database(ctx).First(u, "phone=?", phone).Scan(&m).Error; err != nil {
		return "", transferErr(err)
	}
	return m["password"].(string), nil
}

func (u *User) InitData(ctx *core.Context) error {
	ins := User{
		BaseModel:   types.BaseModel{ID: 1, CreatedAt: time.Now().Unix()},
		TeamID:      1,
		RoleID:      1,
		Name:        "超级管理员",
		Nickname:    "superAdmin",
		Sex:         proto.Bool(true),
		Phone:       "18888888888",
		Password:    tools.HashPwd("123456"),
		Email:       "18888888888@qq.com",
		Status:      proto.Bool(true),
		DisableDesc: proto.String(""),
	}
	return database(ctx).Create(&ins).Error
}
