package controller

import (
	"drone-server/internal/db"
	"drone-server/internal/models"
	"drone-server/utils"

	"github.com/gin-gonic/gin"
)

// ListDrones 无人机列表（分页 + 多条件过滤）
func ListDrones(c *gin.Context) {
	current, size := parsePage(c)
	keyword := c.Query("keyword") // 名称/编号/型号
	storageStatus := c.Query("storageStatus")
	deviceStatus := c.Query("deviceStatus")
	roomID := c.Query("roomId")
	category := c.Query("category")

	q := db.DB.Model(&models.Drone{})
	if keyword != "" {
		q = q.Where("name LIKE ? OR code LIKE ? OR model LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if storageStatus != "" {
		q = q.Where("storage_status = ?", storageStatus)
	}
	if deviceStatus != "" {
		q = q.Where("device_status = ?", deviceStatus)
	}
	if roomID != "" {
		q = q.Where("room_id = ?", roomID)
	}
	if category != "" {
		q = q.Where("category = ?", category)
	}
	var total int64
	q.Count(&total)
	var drones []models.Drone
	q.Order("id desc").Offset((current - 1) * size).Limit(size).Find(&drones)
	utils.OK(c, utils.Page(drones, current, size, total))
}

// GetDrone 无人机详情
func GetDrone(c *gin.Context) {
	var d models.Drone
	if err := db.DB.First(&d, atoiParam(c, "id")).Error; err != nil {
		utils.Fail(c, "无人机不存在")
		return
	}
	utils.OK(c, d)
}

// CreateDrone 新增无人机
func CreateDrone(c *gin.Context) {
	var d models.Drone
	if err := c.ShouldBindJSON(&d); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	if d.Code == "" {
		utils.Fail(c, "RFID/电子标签编号不能为空")
		return
	}
	var cnt int64
	db.DB.Model(&models.Drone{}).Where("code = ?", d.Code).Count(&cnt)
	if cnt > 0 {
		utils.Fail(c, "该编号已存在")
		return
	}
	if d.StorageStatus == "" {
		d.StorageStatus = models.DroneStorageStatusInPosition
	}
	if d.DeviceStatus == "" {
		d.DeviceStatus = models.DroneDeviceStatusIntact
	}
	// 占用对应库位
	if d.PositionCode != "" {
		db.DB.Model(&models.RackPosition{}).Where("code = ?", d.PositionCode).
			Updates(map[string]interface{}{"drone_id": d.ID, "status": 1})
	}
	db.DB.Create(&d)
	writeLog(models.LogCategoryOperation, "drone", "登记无人机: "+d.Name+"("+d.Code+")", "info", currentUsername(c), c.ClientIP())
	utils.OK(c, d)
}

// UpdateDrone 更新无人机
func UpdateDrone(c *gin.Context) {
	id := atoiParam(c, "id")
	var d models.Drone
	if err := db.DB.First(&d, id).Error; err != nil {
		utils.Fail(c, "无人机不存在")
		return
	}
	var req struct {
		Name         string `json:"name"`
		Model        string `json:"model"`
		Category     string `json:"category"`
		DeviceStatus string `json:"deviceStatus"`
		Battery      int    `json:"battery"`
		PositionCode string `json:"positionCode"`
		RoomID       uint   `json:"roomId"`
		Manufacturer string `json:"manufacturer"`
		Remark       string `json:"remark"`
	}
	_ = c.ShouldBindJSON(&req)
	d.Name = req.Name
	d.ModelName = req.Model
	d.Category = req.Category
	if req.DeviceStatus != "" {
		d.DeviceStatus = req.DeviceStatus
	}
	d.Battery = req.Battery
	d.PositionCode = req.PositionCode
	d.RoomID = req.RoomID
	d.Manufacturer = req.Manufacturer
	d.Remark = req.Remark
	// 同步库位占用
	db.DB.Model(&models.RackPosition{}).Where("drone_id = ?", d.ID).Updates(map[string]interface{}{"drone_id": nil, "status": 0})
	if d.PositionCode != "" && d.StorageStatus == models.DroneStorageStatusInPosition {
		db.DB.Model(&models.RackPosition{}).Where("code = ?", d.PositionCode).
			Updates(map[string]interface{}{"drone_id": d.ID, "status": 1})
	}
	db.DB.Save(&d)
	utils.OK(c, d)
}

// DeleteDrone 删除无人机
func DeleteDrone(c *gin.Context) {
	id := atoiParam(c, "id")
	db.DB.Model(&models.RackPosition{}).Where("drone_id = ?", id).Updates(map[string]interface{}{"drone_id": nil, "status": 0})
	db.DB.Delete(&models.Drone{}, id)
	writeLog(models.LogCategoryOperation, "drone", "删除无人机 ID", "info", currentUsername(c), c.ClientIP())
	utils.OK(c, nil)
}

// DroneStats 无人机统计（在位/借出、完好/损坏/报废）
func DroneStats(c *gin.Context) {
	var total, inPos, borrowed, intact, damaged, scrapped int64
	db.DB.Model(&models.Drone{}).Count(&total)
	db.DB.Model(&models.Drone{}).Where("storage_status = ?", models.DroneStorageStatusInPosition).Count(&inPos)
	db.DB.Model(&models.Drone{}).Where("storage_status = ?", models.DroneStorageStatusBorrowed).Count(&borrowed)
	db.DB.Model(&models.Drone{}).Where("device_status = ?", models.DroneDeviceStatusIntact).Count(&intact)
	db.DB.Model(&models.Drone{}).Where("device_status = ?", models.DroneDeviceStatusDamaged).Count(&damaged)
	db.DB.Model(&models.Drone{}).Where("device_status = ?", models.DroneDeviceStatusScrapped).Count(&scrapped)

	// 按类别统计数量
	type catItem struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}
	var byCategory []catItem
	db.DB.Model(&models.Drone{}).Select("category, count(*) as count").Group("category").Scan(&byCategory)

	utils.OK(c, gin.H{
		"total":        total,
		"inPosition":   inPos,
		"borrowed":     borrowed,
		"intact":       intact,
		"damaged":      damaged,
		"scrapped":     scrapped,
		"byCategory":   byCategory,
	})
}
