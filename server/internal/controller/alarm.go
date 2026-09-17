package controller

import (
	"time"

	"drone-server/internal/db"
	"drone-server/internal/models"
	"drone-server/utils"

	"github.com/gin-gonic/gin"
)

// ListAlarms 告警列表（多条件过滤）
func ListAlarms(c *gin.Context) {
	current, size := parsePage(c)
	alarmType := c.Query("type")
	level := c.Query("level")
	status := c.Query("status")
	roomID := c.Query("roomId")
	source := c.Query("source")
	start := c.Query("start")
	end := c.Query("end")

	q := db.DB.Model(&models.AlarmEvent{})
	if alarmType != "" {
		q = q.Where("type = ?", alarmType)
	}
	if level != "" {
		q = q.Where("level = ?", level)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if roomID != "" {
		q = q.Where("room_id = ?", roomID)
	}
	if source != "" {
		q = q.Where("source = ?", source)
	}
	if t, err := time.Parse(time.RFC3339, start); err == nil {
		q = q.Where("created_at >= ?", t)
	}
	if t, err := time.Parse(time.RFC3339, end); err == nil {
		q = q.Where("created_at <= ?", t)
	}
	var total int64
	q.Count(&total)
	var alarms []models.AlarmEvent
	q.Order("id desc").Offset((current - 1) * size).Limit(size).Find(&alarms)
	utils.OK(c, utils.Page(alarms, current, size, total))
}

// CreateAlarm 手动创建告警
func CreateAlarm(c *gin.Context) {
	var a models.AlarmEvent
	if err := c.ShouldBindJSON(&a); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	if a.Status == "" {
		a.Status = models.AlarmStatusActive
	}
	if a.Source == "" {
		a.Source = "system"
	}
	db.DB.Create(&a)
	writeLog(models.LogCategoryAlarm, a.Type, "产生告警: "+a.Title, a.Level, currentUsername(c), c.ClientIP())
	utils.OK(c, a)
}

// ResolveAlarm 处理/解除告警
func ResolveAlarm(c *gin.Context) {
	id := atoiParam(c, "id")
	var a models.AlarmEvent
	if err := db.DB.First(&a, id).Error; err != nil {
		utils.Fail(c, "告警不存在")
		return
	}
	now := time.Now()
	db.DB.Model(&a).Updates(map[string]interface{}{
		"status":     models.AlarmStatusResolved,
		"resolved_at": now,
		"resolved_by": currentUsername(c),
	})
	utils.OK(c, nil)
}

// IgnoreAlarm 忽略告警
func IgnoreAlarm(c *gin.Context) {
	db.DB.Model(&models.AlarmEvent{}).Where("id = ?", atoiParam(c, "id")).
		Update("status", models.AlarmStatusIgnored)
	utils.OK(c, nil)
}

// DeleteAlarm 删除告警
func DeleteAlarm(c *gin.Context) {
	db.DB.Delete(&models.AlarmEvent{}, atoiParam(c, "id"))
	utils.OK(c, nil)
}

// AlarmStats 告警统计（按类型/级别/状态汇总 + 实时活跃数）
func AlarmStats(c *gin.Context) {
	var active, resolved, ignored, critical int64
	db.DB.Model(&models.AlarmEvent{}).Where("status = ?", models.AlarmStatusActive).Count(&active)
	db.DB.Model(&models.AlarmEvent{}).Where("status = ?", models.AlarmStatusResolved).Count(&resolved)
	db.DB.Model(&models.AlarmEvent{}).Where("status = ?", models.AlarmStatusIgnored).Count(&ignored)
	db.DB.Model(&models.AlarmEvent{}).Where("level = ? AND status = ?", models.AlarmLevelCritical, models.AlarmStatusActive).Count(&critical)

	type typeItem struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	var byType []typeItem
	db.DB.Model(&models.AlarmEvent{}).Select("type, count(*) as count").Group("type").Scan(&byType)

	type levelItem struct {
		Level string `json:"level"`
		Count int64  `json:"count"`
	}
	var byLevel []levelItem
	db.DB.Model(&models.AlarmEvent{}).Select("level, count(*) as count").Group("level").Scan(&byLevel)

	utils.OK(c, gin.H{
		"active":   active,
		"resolved": resolved,
		"ignored":  ignored,
		"critical": critical,
		"byType":   byType,
		"byLevel":  byLevel,
	})
}

// PushAlarm 第三方对接接口：推送告警事件
func PushAlarm(c *gin.Context) {
	var a models.AlarmEvent
	if err := c.ShouldBindJSON(&a); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	a.Source = "third-party"
	if a.Status == "" {
		a.Status = models.AlarmStatusActive
	}
	db.DB.Create(&a)
	writeLog(models.LogCategoryAlarm, a.Type, "第三方推送告警: "+a.Title, a.Level, "third-party", c.ClientIP())
	utils.OK(c, a)
}
