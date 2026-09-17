package controller

import (
	"time"

	"drone-server/internal/db"
	"drone-server/internal/models"
	"drone-server/utils"

	"github.com/gin-gonic/gin"
)

// ListAccess 人员进出记录（多条件过滤）
func ListAccess(c *gin.Context) {
	current, size := parsePage(c)
	accessType := c.Query("type")
	roomID := c.Query("roomId")
	person := c.Query("personName")
	cardNo := c.Query("cardNo")
	start := c.Query("start")
	end := c.Query("end")

	q := db.DB.Model(&models.AccessRecord{})
	if accessType != "" {
		q = q.Where("type = ?", accessType)
	}
	if roomID != "" {
		q = q.Where("room_id = ?", roomID)
	}
	if person != "" {
		q = q.Where("person_name LIKE ?", "%"+person+"%")
	}
	if cardNo != "" {
		q = q.Where("card_no LIKE ?", "%"+cardNo+"%")
	}
	if t, err := time.Parse(time.RFC3339, start); err == nil {
		q = q.Where("created_at >= ?", t)
	}
	if t, err := time.Parse(time.RFC3339, end); err == nil {
		q = q.Where("created_at <= ?", t)
	}
	var total int64
	q.Count(&total)
	var records []models.AccessRecord
	q.Order("id desc").Offset((current - 1) * size).Limit(size).Find(&records)
	utils.OK(c, utils.Page(records, current, size, total))
}

// CreateAccess 手动新增进出记录
func CreateAccess(c *gin.Context) {
	var a models.AccessRecord
	if err := c.ShouldBindJSON(&a); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	if a.Type == "" {
		a.Type = models.AccessTypeFacePass
	}
	db.DB.Create(&a)
	writeLog(models.LogCategoryOperation, "access", "人员进出记录: "+a.PersonName+"("+a.Type+")", "info", currentUsername(c), c.ClientIP())
	utils.OK(c, a)
}

// AccessStats 人员进出统计
func AccessStats(c *gin.Context) {
	// 今日零点
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var total, todayCount, cardFail, doorOpen, doorClose, facePass int64
	db.DB.Model(&models.AccessRecord{}).Count(&total)
	db.DB.Model(&models.AccessRecord{}).Where("created_at >= ?", today).Count(&todayCount)
	db.DB.Model(&models.AccessRecord{}).Where("type = ?", models.AccessTypeCardFail).Count(&cardFail)
	db.DB.Model(&models.AccessRecord{}).Where("type = ?", models.AccessTypeDoorOpen).Count(&doorOpen)
	db.DB.Model(&models.AccessRecord{}).Where("type = ?", models.AccessTypeDoorClose).Count(&doorClose)
	db.DB.Model(&models.AccessRecord{}).Where("type = ?", models.AccessTypeFacePass).Count(&facePass)

	utils.OK(c, gin.H{
		"total":      total,
		"today":      todayCount,
		"cardFail":   cardFail,
		"doorOpen":   doorOpen,
		"doorClose":  doorClose,
		"facePass":   facePass,
	})
}

// PushAccess 第三方对接接口：推送人员进出记录
func PushAccess(c *gin.Context) {
	var a models.AccessRecord
	if err := c.ShouldBindJSON(&a); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	if a.Type == "" {
		a.Type = models.AccessTypeFacePass
	}
	db.DB.Create(&a)
	utils.OK(c, a)
}
