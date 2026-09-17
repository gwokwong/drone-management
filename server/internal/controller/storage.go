package controller

import (
	"fmt"

	"drone-server/internal/db"
	"drone-server/internal/models"
	"drone-server/utils"

	"github.com/gin-gonic/gin"
)

// ---------------- 储存室 ----------------

func ListRooms(c *gin.Context) {
	current, size := parsePage(c)
	name := c.Query("name")
	q := db.DB.Model(&models.StorageRoom{})
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	var total int64
	q.Count(&total)
	var rooms []models.StorageRoom
	q.Order("id desc").Offset((current - 1) * size).Limit(size).Find(&rooms)
	utils.OK(c, utils.Page(rooms, current, size, total))
}

func GetRoom(c *gin.Context) {
	var r models.StorageRoom
	if err := db.DB.First(&r, atoiParam(c, "id")).Error; err != nil {
		utils.Fail(c, "储存室不存在")
		return
	}
	utils.OK(c, r)
}

func CreateRoom(c *gin.Context) {
	var r models.StorageRoom
	if err := c.ShouldBindJSON(&r); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	if r.Code == "" || r.Name == "" {
		utils.Fail(c, "编码与名称不能为空")
		return
	}
	db.DB.Create(&r)
	utils.OK(c, r)
}

func UpdateRoom(c *gin.Context) {
	id := atoiParam(c, "id")
	var r models.StorageRoom
	if err := db.DB.First(&r, id).Error; err != nil {
		utils.Fail(c, "储存室不存在")
		return
	}
	_ = c.ShouldBindJSON(&r)
	r.ID = id
	db.DB.Save(&r)
	utils.OK(c, r)
}

func DeleteRoom(c *gin.Context) {
	db.DB.Delete(&models.StorageRoom{}, atoiParam(c, "id"))
	utils.OK(c, nil)
}

// ---------------- 密集架 ----------------

func ListRacks(c *gin.Context) {
	roomID := c.Query("roomId")
	q := db.DB.Model(&models.StorageRack{})
	if roomID != "" {
		q = q.Where("room_id = ?", roomID)
	}
	var racks []models.StorageRack
	q.Order("id asc").Find(&racks)
	utils.OK(c, racks)
}

func CreateRack(c *gin.Context) {
	var r models.StorageRack
	if err := c.ShouldBindJSON(&r); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	if r.Rows <= 0 || r.Cols <= 0 || r.Layers <= 0 {
		utils.Fail(c, "行列层数量必须为正")
		return
	}
	db.DB.Create(&r)
	// 自动生成库位
	for col := 1; col <= r.Cols; col++ {
		for row := 1; row <= r.Rows; row++ {
			for layer := 1; layer <= r.Layers; layer++ {
				pos := models.RackPosition{
					RackID: r.ID,
					Code:   rackPosCode(r.Code, col, row, layer),
					Row:    row,
					Col:    col,
					Layer:  layer,
					Status: 0,
				}
				db.DB.Create(&pos)
			}
		}
	}
	utils.OK(c, r)
}

func rackPosCode(rackCode string, col, row, layer int) string {
	return fmt.Sprintf("%s-%02d-%02d-%d", rackCode, col, row, layer)
}

func UpdateRack(c *gin.Context) {
	id := atoiParam(c, "id")
	var r models.StorageRack
	if err := db.DB.First(&r, id).Error; err != nil {
		utils.Fail(c, "密集架不存在")
		return
	}
	_ = c.ShouldBindJSON(&r)
	r.ID = id
	db.DB.Save(&r)
	utils.OK(c, r)
}

func DeleteRack(c *gin.Context) {
	id := atoiParam(c, "id")
	db.DB.Where("rack_id = ?", id).Delete(&models.RackPosition{})
	db.DB.Delete(&models.StorageRack{}, id)
	utils.OK(c, nil)
}

// RoomMap3D 3D 导航图数据：储存室 + 各密集架库位占用情况
func RoomMap3D(c *gin.Context) {
	id := atoiParam(c, "id")
	var room models.StorageRoom
	if err := db.DB.First(&room, id).Error; err != nil {
		utils.Fail(c, "储存室不存在")
		return
	}
	var racks []models.StorageRack
	db.DB.Where("room_id = ?", id).Order("id asc").Find(&racks)

	type posView struct {
		Code      string `json:"code"`
		Row       int    `json:"row"`
		Col       int    `json:"col"`
		Layer     int    `json:"layer"`
		Occupied  bool   `json:"occupied"`
		DroneCode string `json:"droneCode"`
		DroneID   *uint  `json:"droneId"`
	}
	type rackView struct {
		ID    uint      `json:"id"`
		Code  string    `json:"code"`
		Name  string    `json:"name"`
		Rows  int       `json:"rows"`
		Cols  int       `json:"cols"`
		Layers int      `json:"layers"`
		Positions []posView `json:"positions"`
	}
	var rackViews []rackView
	for _, rk := range racks {
		var positions []models.RackPosition
		db.DB.Where("rack_id = ?", rk.ID).Order("col asc,row asc,layer asc").Find(&positions)
		pv := make([]posView, 0, len(positions))
		for _, p := range positions {
			pv = append(pv, posView{
				Code:      p.Code,
				Row:       p.Row,
				Col:       p.Col,
				Layer:     p.Layer,
				Occupied:  p.Status == 1 && p.DroneID != nil,
				DroneCode: droneCodeByID(p.DroneID),
				DroneID:   p.DroneID,
			})
		}
		rackViews = append(rackViews, rackView{
			ID: rk.ID, Code: rk.Code, Name: rk.Name, Rows: rk.Rows, Cols: rk.Cols, Layers: rk.Layers, Positions: pv,
		})
	}
	// 占用统计
	var totalPos, usedPos int64
	db.DB.Model(&models.RackPosition{}).Where("rack_id IN (?)", rackIDs(racks)).Count(&totalPos)
	db.DB.Model(&models.RackPosition{}).Where("rack_id IN (?) AND status = 1", rackIDs(racks)).Count(&usedPos)

	utils.OK(c, gin.H{
		"room":       room,
		"racks":      rackViews,
		"totalPos":   totalPos,
		"usedPos":    usedPos,
		"freePos":    totalPos - usedPos,
	})
}

func rackIDs(racks []models.StorageRack) []uint {
	ids := make([]uint, 0, len(racks))
	for _, r := range racks {
		ids = append(ids, r.ID)
	}
	return ids
}

func droneCodeByID(id *uint) string {
	if id == nil {
		return ""
	}
	var d models.Drone
	if err := db.DB.Select("code").First(&d, *id).Error; err != nil {
		return ""
	}
	return d.Code
}

// ---------------- 环境设备 ----------------

func ListEnvDevices(c *gin.Context) {
	roomID := c.Query("roomId")
	devType := c.Query("type")
	status := c.Query("status")
	q := db.DB.Model(&models.EnvDevice{})
	if roomID != "" {
		q = q.Where("room_id = ?", roomID)
	}
	if devType != "" {
		q = q.Where("type = ?", devType)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var devices []models.EnvDevice
	q.Order("id asc").Find(&devices)
	utils.OK(c, devices)
}

func CreateEnvDevice(c *gin.Context) {
	var d models.EnvDevice
	if err := c.ShouldBindJSON(&d); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	if d.Status == "" {
		d.Status = models.EnvDeviceStatusNormal
	}
	db.DB.Create(&d)
	utils.OK(c, d)
}

func UpdateEnvDevice(c *gin.Context) {
	id := atoiParam(c, "id")
	var d models.EnvDevice
	if err := db.DB.First(&d, id).Error; err != nil {
		utils.Fail(c, "设备不存在")
		return
	}
	_ = c.ShouldBindJSON(&d)
	d.ID = id
	db.DB.Save(&d)
	utils.OK(c, d)
}

func DeleteEnvDevice(c *gin.Context) {
	db.DB.Delete(&models.EnvDevice{}, atoiParam(c, "id"))
	utils.OK(c, nil)
}

// EnvDeviceSummary 环境设备运行状态汇总（设备自检展示）
func EnvDeviceSummary(c *gin.Context) {
	var total, normal, abnormal, offline int64
	db.DB.Model(&models.EnvDevice{}).Count(&total)
	db.DB.Model(&models.EnvDevice{}).Where("status = ?", models.EnvDeviceStatusNormal).Count(&normal)
	db.DB.Model(&models.EnvDevice{}).Where("status = ?", models.EnvDeviceStatusAbnormal).Count(&abnormal)
	db.DB.Model(&models.EnvDevice{}).Where("status = ?", models.EnvDeviceStatusOffline).Count(&offline)
	utils.OK(c, gin.H{
		"total":    total,
		"normal":   normal,
		"abnormal": abnormal,
		"offline":  offline,
	})
}
