package controller

import (
	"encoding/json"
	"strconv"
	"time"

	"drone-server/internal/db"
	"drone-server/internal/models"

	"github.com/gin-gonic/gin"
)

func parsePage(c *gin.Context) (current int, size int) {
	current, _ = strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ = strconv.Atoi(c.DefaultQuery("size", "10"))
	if current < 1 {
		current = 1
	}
	if size < 1 {
		size = 10
	}
	return
}

func atoiParam(c *gin.Context, key string) uint {
	v, _ := strconv.ParseUint(c.Param(key), 10, 64)
	return uint(v)
}

func currentRoles(c *gin.Context) []string {
	if v, ok := c.Get("roles"); ok {
		if rs, ok := v.([]string); ok {
			return rs
		}
	}
	return []string{}
}

func currentUserID(c *gin.Context) uint {
	if v, ok := c.Get("uid"); ok {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}

func currentUsername(c *gin.Context) string {
	if v, ok := c.Get("username"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// writeLog 写入系统日志（按设备/分类管理）
func writeLog(category, deviceType, content, level, operator, ip string) {
	_ = db.DB.Create(&models.SystemLog{
		Category:   category,
		DeviceType: deviceType,
		Content:    content,
		Level:      level,
		Operator:   operator,
		IP:         ip,
	}).Error
}

func parseJSONArray(s string) []string {
	if s == "" {
		return []string{}
	}
	var arr []string
	_ = json.Unmarshal([]byte(s), &arr)
	return arr
}

func toJSONArray(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}
	b, _ := json.Marshal(arr)
	return string(b)
}

// gormNow 返回当前时间的指针，用于 *time.Time 字段
func gormNow() *time.Time {
	t := time.Now()
	return &t
}
