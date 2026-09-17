package main

import (
	"fmt"
	"os"
	"path/filepath"

	"drone-server/config"
	"drone-server/internal/db"
	"drone-server/internal/router"
)

// resolveBaseDir 确定配置与数据文件的基础目录。
// 优先使用可执行文件所在目录，避免从不同工作目录启动时
// 配置文件找不到、数据库文件散落到意外位置。
func resolveBaseDir() string {
	if exe, err := os.Executable(); err == nil {
		if real, err := filepath.EvalSymlinks(exe); err == nil {
			exe = real
		}
		return filepath.Dir(exe)
	}
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "."
}

// resolveConfigPath 依次尝试：环境变量 -> 当前目录 -> 可执行文件同级目录
func resolveConfigPath(baseDir string) string {
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		return p
	}
	const rel = "config/config.yaml"
	if _, err := os.Stat(rel); err == nil {
		return rel
	}
	return filepath.Join(baseDir, rel)
}

func main() {
	baseDir := resolveBaseDir()
	cfgPath := resolveConfigPath(baseDir)

	if err := config.InitConfig(cfgPath); err != nil {
		panic("配置加载失败: " + err.Error())
	}

	// 相对路径的 SQLite 文件固定落在基础目录下，避免随工作目录漂移
	if config.GlobalConfig.Database.Driver != "mysql" {
		p := config.GlobalConfig.Database.SQLite.Path
		if p != "" && !filepath.IsAbs(p) {
			config.GlobalConfig.Database.SQLite.Path = filepath.Join(baseDir, p)
			_ = os.MkdirAll(filepath.Dir(config.GlobalConfig.Database.SQLite.Path), 0o755)
		}
	}

	if err := db.InitDB(); err != nil {
		panic("数据库初始化失败: " + err.Error())
	}
	db.Seed()

	r := router.SetupRouter()
	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	fmt.Printf("无人机管理平台后端已启动: http://localhost%s  (数据库: %s)\n", addr, config.GlobalConfig.Database.Driver)
	_ = r.Run(addr)
}
