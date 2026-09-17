package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 全局配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	CORS     CORSConfig     `yaml:"cors"`
	App      AppConfig      `yaml:"app"`
}

// AppConfig 应用级配置
type AppConfig struct {
	// SeedDemoData 是否在启动时写入演示业务数据（无人机、储存室、告警、进出记录等）。
	// 默认 false：仅初始化角色与管理员账号，系统为空系统，不含任何测试数据。
	// 需要演示/联调时置为 true，或设置环境变量 SEED_DEMO_DATA=true。
	SeedDemoData bool `yaml:"seedDemoData"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Driver string       `yaml:"driver"` // sqlite | mysql
	SQLite SQLiteConfig `yaml:"sqlite"`
	MySQL  MySQLConfig  `yaml:"mysql"`
}

type SQLiteConfig struct {
	Path string `yaml:"path"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire string `yaml:"expire"`
}

type CORSConfig struct {
	AllowOrigins []string `yaml:"allowOrigins"`
}

// GlobalConfig 运行时全局配置
var GlobalConfig Config

// InitConfig 加载配置文件，未找到时使用内置默认值。
// 切换数据库：将 database.driver 改为 "mysql" 并填写 mysql 配置即可，无需改动代码。
func InitConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		// 未提供配置文件时回退到默认配置（本地 sqlite）
		GlobalConfig = defaultConfig()
		return nil
	}
	if err := yaml.Unmarshal(data, &GlobalConfig); err != nil {
		return err
	}
	applyDefaults()
	return nil
}

func applyDefaults() {
	if GlobalConfig.Server.Port == 0 {
		GlobalConfig.Server.Port = 8080
	}
	if GlobalConfig.Server.Mode == "" {
		GlobalConfig.Server.Mode = "debug"
	}
	if GlobalConfig.JWT.Secret == "" {
		GlobalConfig.JWT.Secret = "drone-server-default-secret"
	}
	if GlobalConfig.JWT.Expire == "" {
		GlobalConfig.JWT.Expire = "72h"
	}
	if GlobalConfig.Database.Driver == "" {
		GlobalConfig.Database.Driver = "sqlite"
	}
	if GlobalConfig.Database.SQLite.Path == "" {
		GlobalConfig.Database.SQLite.Path = "data/drone.db"
	}
	if GlobalConfig.Database.MySQL.Charset == "" {
		GlobalConfig.Database.MySQL.Charset = "utf8mb4"
	}
	// 环境变量优先级高于配置文件，便于临时开启演示数据而不改动 yaml
	if os.Getenv("SEED_DEMO_DATA") == "true" {
		GlobalConfig.App.SeedDemoData = true
	}
	// 确保 sqlite 文件所在目录存在
	if GlobalConfig.Database.Driver != "mysql" {
		if dir := filepath.Dir(GlobalConfig.Database.SQLite.Path); dir != "" && dir != "." {
			_ = os.MkdirAll(dir, 0o755)
		}
	}
}

func defaultConfig() Config {
	c := Config{}
	c.Server.Port = 8080
	c.Server.Mode = "debug"
	c.JWT.Secret = "drone-server-default-secret"
	c.JWT.Expire = "72h"
	c.Database.Driver = "sqlite"
	c.Database.SQLite.Path = "data/drone.db"
	c.Database.MySQL.Host = "127.0.0.1"
	c.Database.MySQL.Port = 3306
	c.Database.MySQL.User = "root"
	c.Database.MySQL.Password = "root"
	c.Database.MySQL.DBName = "drone"
	c.Database.MySQL.Charset = "utf8mb4"
	c.CORS.AllowOrigins = []string{"*"}
	return c
}
