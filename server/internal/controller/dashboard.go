package controller

import (
	"time"

	"drone-server/internal/db"
	"drone-server/internal/models"
	"drone-server/utils"

	"github.com/gin-gonic/gin"
)

// DashboardOverview 大屏数据展示总览（实时聚合）
func DashboardOverview(c *gin.Context) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var droneTotal, inPos, borrowed, intact, damaged, scrapped, activeBorrow int64
	db.DB.Model(&models.Drone{}).Count(&droneTotal)
	db.DB.Model(&models.Drone{}).Where("storage_status = ?", models.DroneStorageStatusInPosition).Count(&inPos)
	db.DB.Model(&models.Drone{}).Where("storage_status = ?", models.DroneStorageStatusBorrowed).Count(&borrowed)
	db.DB.Model(&models.Drone{}).Where("device_status = ?", models.DroneDeviceStatusIntact).Count(&intact)
	db.DB.Model(&models.Drone{}).Where("device_status = ?", models.DroneDeviceStatusDamaged).Count(&damaged)
	db.DB.Model(&models.Drone{}).Where("device_status = ?", models.DroneDeviceStatusScrapped).Count(&scrapped)
	db.DB.Model(&models.BorrowRecord{}).Where("status = ?", models.BorrowStatusBorrowing).Count(&activeBorrow)

	var alarmActive, alarmCritical int64
	db.DB.Model(&models.AlarmEvent{}).Where("status = ?", models.AlarmStatusActive).Count(&alarmActive)
	db.DB.Model(&models.AlarmEvent{}).Where("level = ? AND status = ?", models.AlarmLevelCritical, models.AlarmStatusActive).Count(&alarmCritical)

	var accessToday int64
	db.DB.Model(&models.AccessRecord{}).Where("created_at >= ?", today).Count(&accessToday)

	var envTotal, envAbnormal, envOffline int64
	db.DB.Model(&models.EnvDevice{}).Count(&envTotal)
	db.DB.Model(&models.EnvDevice{}).Where("status = ?", models.EnvDeviceStatusAbnormal).Count(&envAbnormal)
	db.DB.Model(&models.EnvDevice{}).Where("status = ?", models.EnvDeviceStatusOffline).Count(&envOffline)

	// 最近告警
	var recentAlarms []models.AlarmEvent
	db.DB.Where("status = ?", models.AlarmStatusActive).Order("id desc").Limit(5).Find(&recentAlarms)

	// 最近进出
	var recentAccess []models.AccessRecord
	db.DB.Order("id desc").Limit(5).Find(&recentAccess)

	// 各储存室占用
	type roomStat struct {
		RoomID uint   `json:"roomId"`
		Name   string `json:"name"`
		Used   int64  `json:"used"`
		Total  int64  `json:"total"`
	}
	var rooms []models.StorageRoom
	db.DB.Find(&rooms)
	roomStats := make([]roomStat, 0, len(rooms))
	for _, r := range rooms {
		var used, total int64
		db.DB.Model(&models.RackPosition{}).Where("rack_id IN (?)", rackIDsOfRoom(r.ID)).Count(&total)
		db.DB.Model(&models.RackPosition{}).Where("rack_id IN (?) AND status = 1", rackIDsOfRoom(r.ID)).Count(&used)
		roomStats = append(roomStats, roomStat{RoomID: r.ID, Name: r.Name, Used: used, Total: total})
	}

	utils.OK(c, gin.H{
		"drone": gin.H{
			"total":      droneTotal,
			"inPosition": inPos,
			"borrowed":   borrowed,
			"intact":     intact,
			"damaged":    damaged,
			"scrapped":   scrapped,
			"activeBorrow": activeBorrow,
		},
		"alarm": gin.H{
			"active":   alarmActive,
			"critical": alarmCritical,
		},
		"accessToday": accessToday,
		"env": gin.H{
			"total":    envTotal,
			"abnormal": envAbnormal,
			"offline":  envOffline,
		},
		"recentAlarms":  recentAlarms,
		"recentAccess":  recentAccess,
		"roomStats":     roomStats,
	})
}

func rackIDsOfRoom(roomID uint) []uint {
	var racks []models.StorageRack
	db.DB.Select("id").Where("room_id = ?", roomID).Find(&racks)
	ids := make([]uint, 0, len(racks))
	for _, r := range racks {
		ids = append(ids, r.ID)
	}
	return ids
}
