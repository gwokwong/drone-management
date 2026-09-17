package db

import (
	"encoding/json"
	"fmt"
	"time"

	"drone-server/config"
	"drone-server/internal/models"
	"drone-server/utils"
)

// Seed 写入初始数据。
//
// 基础数据（角色、管理员账号）始终写入，保证系统可登录可用；
// 演示业务数据（操作员账号、无人机、储存室、环境设备、摄像头、告警、进出记录等）
// 由配置 app.seedDemoData 控制，默认不写入，交付环境为不含测试数据的干净系统。
func Seed() {
	seedRoles()
	seedAdmin()

	if !config.GlobalConfig.App.SeedDemoData {
		fmt.Println("[seed] 演示数据已关闭（app.seedDemoData=false）：系统为干净空系统，仅含角色与管理员账号")
		return
	}

	seedOperator()
	roomID := seedRoomsAndRacks()
	seedDrones(roomID)
	seedEnvDevices(roomID)
	seedVideos(roomID)
	seedSampleEvents(roomID)
	fmt.Println("[seed] 已写入演示业务数据（无人机 / 储存室 / 告警 / 进出记录 / 摄像头）")
}

func jsonArr(items ...string) string {
	b, _ := json.Marshal(items)
	return string(b)
}

func seedRoles() {
	var cnt int64
	DB.Model(&models.Role{}).Count(&cnt)
	if cnt > 0 {
		return
	}
	roles := []models.Role{
		{Code: "R_SUPER", Name: "超级管理员", Permissions: jsonArr("*"), Enabled: 1, Remark: "拥有全部权限"},
		{Code: "R_ADMIN", Name: "系统管理员", Permissions: jsonArr("dashboard", "drone", "storage", "alarm", "access", "video", "report", "log", "user", "role"), Enabled: 1, Remark: "系统配置与用户管理"},
		{Code: "R_OPERATOR", Name: "操作员", Permissions: jsonArr("dashboard", "drone", "borrow", "storage", "alarm", "access", "video"), Enabled: 1, Remark: "日常借还与监控"},
	}
	DB.Create(&roles)
}

// seedAdmin 初始化管理员账号（系统可用性的必要数据，始终写入）
func seedAdmin() {
	var cnt int64
	DB.Model(&models.User{}).Where("username = ?", "admin").Count(&cnt)
	if cnt > 0 {
		return
	}
	DB.Create(&models.User{
		Username: "admin", Password: utils.HashPassword("admin123"), Nickname: "系统管理员",
		Email: "admin@drone.local", Status: 1, RoleCodes: jsonArr("R_SUPER"), Remark: "默认超级管理员",
	})
}

// seedOperator 演示用的操作员账号（仅演示数据开启时写入）
func seedOperator() {
	var cnt int64
	DB.Model(&models.User{}).Where("username = ?", "operator").Count(&cnt)
	if cnt > 0 {
		return
	}
	DB.Create(&models.User{
		Username: "operator", Password: utils.HashPassword("operator123"), Nickname: "操作员小李",
		Email: "operator@drone.local", Status: 1, RoleCodes: jsonArr("R_OPERATOR"), Remark: "日常操作员",
	})
}

func seedRoomsAndRacks() uint {
	var cnt int64
	DB.Model(&models.StorageRoom{}).Count(&cnt)
	if cnt > 0 {
		var first models.StorageRoom
		DB.First(&first)
		return first.ID
	}
	room := models.StorageRoom{
		Code: "ROOM-01", Name: "一号无人机储存室", Location: "A栋1层东侧", Capacity: 200,
		EnvTemp: 22.5, EnvHumidity: 45.0, Status: 1, Remark: "主储存室",
	}
	DB.Create(&room)
	rack := models.StorageRack{RoomID: room.ID, Code: "RACK-A", Name: "A区密集架", Rows: 4, Cols: 5, Layers: 3}
	DB.Create(&rack)
	for col := 1; col <= rack.Cols; col++ {
		for row := 1; row <= rack.Rows; row++ {
			for layer := 1; layer <= rack.Layers; layer++ {
				DB.Create(&models.RackPosition{
					RackID: rack.ID,
					Code:   rackPosCode(rack.Code, col, row, layer),
					Row:    row, Col: col, Layer: layer, Status: 0,
				})
			}
		}
	}
	return room.ID
}

func rackPosCode(rackCode string, col, row, layer int) string {
	return fmt.Sprintf("%s-%02d-%02d-%d", rackCode, col, row, layer)
}

func seedDrones(roomID uint) {
	var cnt int64
	DB.Model(&models.Drone{}).Count(&cnt)
	if cnt > 0 {
		return
	}
	specs := []struct {
		code, name, model, cat, pos string
	}{
		{"RFID0001", "大疆M300-RTK-01", "M300 RTK", "多旋翼", "RACK-A-01-01-1"},
		{"RFID0002", "大疆M30T-01", "M30T", "多旋翼", "RACK-A-02-01-1"},
		{"RFID0003", "纵横CW-25-01", "CW-25", "垂直起降固定翼", "RACK-A-03-01-1"},
		{"RFID0004", "大疆M300-RTK-02", "M300 RTK", "多旋翼", "RACK-A-01-02-1"},
		{"RFID0005", "大疆Mavic3E-01", "Mavic 3E", "多旋翼", "RACK-A-02-02-1"},
		{"RFID0006", "极飞V40-01", "V40", "多旋翼", "RACK-A-03-02-1"},
		{"RFID0007", "大疆M30T-02", "M30T", "多旋翼", ""},
		{"RFID0008", "纵横CW-25-02", "CW-25", "垂直起降固定翼", ""},
	}
	now := time.Now()
	for i, s := range specs {
		ds := models.DroneStorageStatusInPosition
		dst := models.DroneDeviceStatusIntact
		if i == 6 {
			ds = models.DroneStorageStatusBorrowed
		}
		if i == 7 {
			dst = models.DroneDeviceStatusDamaged
		}
		d := models.Drone{
			Code: s.code, Name: s.name, ModelName: s.model, Category: s.cat,
			StorageStatus: ds, DeviceStatus: dst, Battery: 100 - i*5,
			PositionCode: s.pos, RoomID: roomID, Manufacturer: "厂商",
			PurchaseDate: &now,
		}
		DB.Create(&d)
		if s.pos != "" {
			DB.Model(&models.RackPosition{}).Where("code = ?", s.pos).
				Updates(map[string]interface{}{"drone_id": d.ID, "status": 1})
		}
		if ds == models.DroneStorageStatusBorrowed {
			exp := now.Add(2 * time.Hour)
			DB.Create(&models.BorrowRecord{
				DroneID: d.ID, DroneCode: d.Code, DroneName: d.Name,
				BorrowerName: "外勤一组", BorrowTime: now.Add(-3 * time.Hour),
				ExpectedReturnTime: &exp, Status: models.BorrowStatusBorrowing,
				PositionCode: d.PositionCode,
			})
		}
	}
	// 报废样例
	DB.Create(&models.Drone{
		Code: "RFID0099", Name: "老旧机型-待报废", ModelName: "X100", Category: "多旋翼",
		StorageStatus: models.DroneStorageStatusInPosition, DeviceStatus: models.DroneDeviceStatusScrapped,
		Battery: 0, RoomID: roomID, PositionCode: "RACK-A-04-03-1",
	})
	DB.Model(&models.RackPosition{}).Where("code = ?", "RACK-A-04-03-1").
		Updates(map[string]interface{}{"drone_id": 0, "status": 1})
}

func seedEnvDevices(roomID uint) {
	var cnt int64
	DB.Model(&models.EnvDevice{}).Count(&cnt)
	if cnt > 0 {
		return
	}
	now := time.Now()
	devices := []models.EnvDevice{
		{Name: "温湿度传感器-01", Type: models.EnvDeviceTypeTemp, RoomID: roomID, Status: models.EnvDeviceStatusNormal, Value: 22.5, Unit: "℃", LastReportAt: &now},
		{Name: "温湿度传感器-02", Type: models.EnvDeviceTypeHumidity, RoomID: roomID, Status: models.EnvDeviceStatusNormal, Value: 45.0, Unit: "%", LastReportAt: &now},
		{Name: "烟感探测器-01", Type: models.EnvDeviceTypeSmoke, RoomID: roomID, Status: models.EnvDeviceStatusNormal, Value: 0, Unit: "", LastReportAt: &now},
		{Name: "水浸传感器-01", Type: models.EnvDeviceTypeWater, RoomID: roomID, Status: models.EnvDeviceStatusNormal, Value: 0, Unit: "", LastReportAt: &now},
		{Name: "门禁控制器-01", Type: models.EnvDeviceTypeDoor, RoomID: roomID, Status: models.EnvDeviceStatusNormal, Value: 1, Unit: "", LastReportAt: &now},
		{Name: "供电监测-01", Type: models.EnvDeviceTypePower, RoomID: roomID, Status: models.EnvDeviceStatusAbnormal, Value: 220, Unit: "V", LastReportAt: &now},
	}
	DB.Create(&devices)
}

func seedVideos(roomID uint) {
	var cnt int64
	DB.Model(&models.VideoChannel{}).Count(&cnt)
	if cnt > 0 {
		return
	}
	videos := []models.VideoChannel{
		{Name: "储存室主摄", RoomID: roomID, URL: "rtsp://192.168.1.101/live/main", Status: models.VideoStatusOnline},
		{Name: "密集架A摄", RoomID: roomID, RackID: 1, URL: "rtsp://192.168.1.102/live/rackA", Status: models.VideoStatusOnline},
		{Name: "门口摄", RoomID: roomID, URL: "rtsp://192.168.1.103/live/door", Status: models.VideoStatusOffline},
	}
	DB.Create(&videos)
}

func seedSampleEvents(roomID uint) {
	var cnt int64
	DB.Model(&models.AlarmEvent{}).Count(&cnt)
	if cnt > 0 {
		return
	}
	alarms := []models.AlarmEvent{
		{Type: models.AlarmTypeNetwork, Level: models.AlarmLevelWarning, Title: "储存室网络中断", Content: "门口摄像机离线超过 60s", Source: "system", RoomID: roomID, Status: models.AlarmStatusActive},
		{Type: models.AlarmTypeWater, Level: models.AlarmLevelCritical, Title: "水浸告警", Content: "水浸传感器检测到积水", Source: "system", RoomID: roomID, Status: models.AlarmStatusActive},
		{Type: models.AlarmTypeDevice, Level: models.AlarmLevelInfo, Title: "供电电压异常", Content: "供电监测-01 电压波动", Source: "system", RoomID: roomID, Status: models.AlarmStatusResolved},
	}
	DB.Create(&alarms)

	access := []models.AccessRecord{
		{Type: models.AccessTypeFacePass, PersonName: "张工", CardNo: "C1001", RoomID: roomID},
		{Type: models.AccessTypeDoorOpen, PersonName: "系统", CardNo: "C1001", RoomID: roomID},
		{Type: models.AccessTypeCardFail, PersonName: "未知", CardNo: "C9999", RoomID: roomID},
		{Type: models.AccessTypeFacePass, PersonName: "李工", CardNo: "C1002", RoomID: roomID},
		{Type: models.AccessTypeDoorClose, PersonName: "系统", CardNo: "C1002", RoomID: roomID},
	}
	for i := range access {
		DB.Create(&access[i])
	}
}
