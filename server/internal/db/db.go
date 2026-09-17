package db

import (
	"fmt"

	"drone-server/config"
	"drone-server/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// InitDB 依据配置初始化数据库（sqlite 或 mysql），并执行自动迁移。
func InitDB() error {
	cfg := config.GlobalConfig.Database
	var dialector gorm.Dialector

	if cfg.Driver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.MySQL.User, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.DBName, cfg.MySQL.Charset)
		dialector = mysql.Open(dsn)
	} else {
		dialector = sqlite.Open(cfg.SQLite.Path)
	}

	gdb, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return err
	}
	DB = gdb
	return migrate()
}

func migrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Drone{},
		&models.StorageRoom{},
		&models.StorageRack{},
		&models.RackPosition{},
		&models.BorrowRecord{},
		&models.AccessRecord{},
		&models.AlarmEvent{},
		&models.VideoChannel{},
		&models.OpenRackRecord{},
		&models.EnvDevice{},
		&models.SystemLog{},
	)
}
