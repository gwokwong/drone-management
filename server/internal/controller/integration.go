package controller

import (
	"drone-server/internal/db"
	"drone-server/internal/models"
	"drone-server/utils"

	"github.com/gin-gonic/gin"
)

// 第三方对接接口：以下接口为系统预留的标准对接入口，第三方软件/设备可直接调用，
// 实现实时数据传输、告警推送、RFID 借还触发、环境设备状态上报等能力。

// IntegrationAlarmPush 第三方推送告警
func IntegrationAlarmPush(c *gin.Context) {
	PushAlarm(c)
}

// IntegrationAccessPush 第三方推送人员进出记录
func IntegrationAccessPush(c *gin.Context) {
	PushAccess(c)
}

// IntegrationEnvPush 第三方上报环境设备状态（按 ID 更新，无则创建）
func IntegrationEnvPush(c *gin.Context) {
	var d models.EnvDevice
	if err := c.ShouldBindJSON(&d); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	if d.Status == "" {
		d.Status = models.EnvDeviceStatusNormal
	}
	if d.ID != 0 {
		db.DB.Model(&models.EnvDevice{}).Where("id = ?", d.ID).Updates(map[string]interface{}{
			"value":          d.Value,
			"status":         d.Status,
			"last_report_at": gormNow(),
		})
		utils.OK(c, d)
		return
	}
	if d.Name != "" {
		var exist models.EnvDevice
		if err := db.DB.Where("name = ? AND room_id = ?", d.Name, d.RoomID).First(&exist).Error; err == nil {
			db.DB.Model(&exist).Updates(map[string]interface{}{
				"value":          d.Value,
				"status":         d.Status,
				"last_report_at": gormNow(),
			})
			utils.OK(c, exist)
			return
		}
	}
	d.LastReportAt = gormNow()
	db.DB.Create(&d)
	utils.OK(c, d)
}

// IntegrationRFIDBorrow 第三方 RFID 触发借出
func IntegrationRFIDBorrow(c *gin.Context) {
	ScanBorrow(c)
}

// IntegrationRFIDReturn 第三方 RFID 触发归还
func IntegrationRFIDReturn(c *gin.Context) {
	ScanReturn(c)
}
