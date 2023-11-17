package install

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
	"gorm.io/gorm"
)

type Model interface {
	TableName() string
	InitData(ctx *core.Context) error
}

type install struct {
	db *gorm.DB
	ms []Model
}

func Init() {
	ins := install{
		ms: []Model{
			&model.Team{},
			&model.Role{},
			&model.User{},
			&model.LoginLog{},
			&model.Menu{},
			&model.RoleMenu{},
			&model.Notice{},
			&model.NoticeUser{},
		},
		db: core.GlobalOrm().GetDB(model.DBName()),
	}

	// 判断是否安装
	if ins.IsInstall() {
		return
	}

	// 进行安装
	if err := ins.Install(); err != nil {
		panic("系统初始化失败" + err.Error())
	}

}

func (ins *install) Install() error {
	for _, tb := range ins.ms {
		if err := ins.db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().AutoMigrate(tb); err != nil {
			return err
		}
		if err := tb.InitData(core.New(&gin.Context{})); err != nil {
			return err
		}
	}
	return nil
}

func (ins *install) IsInstall() bool {
	is := true
	for _, tb := range ins.ms {
		is = is && ins.db.Migrator().HasTable(tb)
	}
	return is
}
