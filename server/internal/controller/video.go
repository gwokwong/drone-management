package controller

import (
	"time"

	"drone-server/internal/db"
	"drone-server/internal/models"
	"drone-server/utils"

	"github.com/gin-gonic/gin"
)

// ListVideos 视频通道列表
func ListVideos(c *gin.Context) {
	roomID := c.Query("roomId")
	q := db.DB.Model(&models.VideoChannel{})
	if roomID != "" {
		q = q.Where("room_id = ?", roomID)
	}
	var videos []models.VideoChannel
	q.Order("id asc").Find(&videos)
	utils.OK(c, videos)
}

// CreateVideo 新增视频通道
func CreateVideo(c *gin.Context) {
	var v models.VideoChannel
	if err := c.ShouldBindJSON(&v); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	if v.Status == "" {
		v.Status = models.VideoStatusOnline
	}
	db.DB.Create(&v)
	utils.OK(c, v)
}

// UpdateVideo 更新视频通道
func UpdateVideo(c *gin.Context) {
	id := atoiParam(c, "id")
	var v models.VideoChannel
	if err := db.DB.First(&v, id).Error; err != nil {
		utils.Fail(c, "通道不存在")
		return
	}
	_ = c.ShouldBindJSON(&v)
	v.ID = id
	db.DB.Save(&v)
	utils.OK(c, v)
}

// DeleteVideo 删除视频通道
func DeleteVideo(c *gin.Context) {
	db.DB.Delete(&models.VideoChannel{}, atoiParam(c, "id"))
	utils.OK(c, nil)
}

// OpenRack 开架录制：调用密集架摄像头录制视频并记录开架
func OpenRack(c *gin.Context) {
	id := atoiParam(c, "id")
	var v models.VideoChannel
	if err := db.DB.First(&v, id).Error; err != nil {
		utils.Fail(c, "视频通道不存在")
		return
	}
	var req struct {
		Operator string `json:"operator"`
		Remark   string `json:"remark"`
	}
	_ = c.ShouldBindJSON(&req)
	rec := models.OpenRackRecord{
		VideoChannelID: v.ID,
		RackID:         v.RackID,
		RackCode:       rackCodeByID(v.RackID),
		Operator:       req.Operator,
		StartTime:      time.Now(),
		Status:         models.OpenRackStatusRecording,
		Remark:         req.Remark,
		VideoURL:       v.URL, // 实际系统此处拼接回放地址
	}
	db.DB.Create(&rec)
	// 记录开架动作日志
	writeLog(models.LogCategoryOperation, "video", "开架录制: 通道"+v.Name+" 操作人:"+req.Operator, "info", currentUsername(c), c.ClientIP())
	utils.OK(c, rec)
}

func rackCodeByID(id uint) string {
	var r models.StorageRack
	if err := db.DB.Select("code").First(&r, id).Error; err != nil {
		return ""
	}
	return r.Code
}

// StopOpenRack 结束开架录制（生成回放记录）
func StopOpenRack(c *gin.Context) {
	id := atoiParam(c, "id")
	var rec models.OpenRackRecord
	if err := db.DB.First(&rec, id).Error; err != nil {
		utils.Fail(c, "开架记录不存在")
		return
	}
	now := time.Now()
	db.DB.Model(&rec).Updates(map[string]interface{}{
		"end_time": now,
		"status":   models.OpenRackStatusFinished,
	})
	utils.OK(c, nil)
}

// ListOpenRacks 开架记录列表（支持回放浏览）
func ListOpenRacks(c *gin.Context) {
	current, size := parsePage(c)
	rackID := c.Query("rackId")
	q := db.DB.Model(&models.OpenRackRecord{})
	if rackID != "" {
		q = q.Where("rack_id = ?", rackID)
	}
	var total int64
	q.Count(&total)
	var records []models.OpenRackRecord
	q.Order("id desc").Offset((current - 1) * size).Limit(size).Find(&records)
	utils.OK(c, utils.Page(records, current, size, total))
}
