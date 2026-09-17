package controller

import (
	"time"

	"drone-server/internal/db"
	"drone-server/internal/models"
	"drone-server/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// buildLogQuery 构造系统日志查询条件（支持设备分类、模糊查询、时间范围）
func buildLogQuery(c *gin.Context) *gorm.DB {
	q := db.DB.Model(&models.SystemLog{})
	if v := c.Query("category"); v != "" {
		q = q.Where("category = ?", v)
	}
	if v := c.Query("deviceType"); v != "" {
		q = q.Where("device_type = ?", v)
	}
	if v := c.Query("level"); v != "" {
		q = q.Where("level = ?", v)
	}
	if v := c.Query("keyword"); v != "" {
		q = q.Where("content LIKE ?", "%"+v+"%")
	}
	if v := c.Query("start"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			q = q.Where("created_at >= ?", t)
		}
	}
	if v := c.Query("end"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			q = q.Where("created_at <= ?", t)
		}
	}
	return q
}

// ListLogs 系统日志查询（支持按设备类型、自定义时间、模糊查询）
func ListLogs(c *gin.Context) {
	current, size := parsePage(c)
	q := buildLogQuery(c)
	var total int64
	q.Count(&total)
	var logs []models.SystemLog
	q.Order("id desc").Offset((current - 1) * size).Limit(size).Find(&logs)
	utils.OK(c, utils.Page(logs, current, size, total))
}

// ExportLogs 系统日志导出为 Excel
func ExportLogs(c *gin.Context) {
	q := buildLogQuery(c)
	var logs []models.SystemLog
	q.Order("id desc").Find(&logs)
	headers := []string{"时间", "分类", "设备类型", "级别", "操作人", "内容", "IP"}
	rows := make([][]interface{}, 0, len(logs))
	for _, l := range logs {
		rows = append(rows, []interface{}{
			l.CreatedAt.Format("2006-01-02 15:04:05"),
			l.Category, l.DeviceType, l.Level, l.Operator, l.Content, l.IP,
		})
	}
	exportExcel(c, "system_logs.xlsx", headers, rows)
}
