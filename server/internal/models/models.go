package models

import (
	"time"

	"gorm.io/gorm"
)

// 通用状态枚举（字符串常量，便于与前端/第三方对接）
const (
	// 无人机存放状态
	DroneStorageStatusInPosition = "in_position" // 在位
	DroneStorageStatusBorrowed   = "borrowed"    // 借出

	// 无人机设备状态
	DroneDeviceStatusIntact   = "intact"   // 完好
	DroneDeviceStatusDamaged  = "damaged"  // 损坏待维修
	DroneDeviceStatusScrapped = "scrapped" // 报废

	// 人员进出记录类型
	AccessTypeCardFail = "card_auth_fail" // 卡号认证失败超次报警
	AccessTypeDoorOpen = "door_open"      // 门锁打开
	AccessTypeDoorClose = "door_close"    // 门锁关闭
	AccessTypeFacePass = "face_auth_pass" // 人脸认证通过

	// 告警类型
	AlarmTypeNetwork   = "network_disconnect" // 断网异常
	AlarmTypeWater     = "water_leak"         // 漏水异常
	AlarmTypeDoor      = "door_alarm"         // 门锁异常
	AlarmTypeDevice    = "device_fault"       // 设备故障
	AlarmTypeLowBattery = "low_battery"       // 低电量
	AlarmTypeEnv       = "env_abnormal"       // 环境异常

	// 告警级别
	AlarmLevelInfo     = "info"
	AlarmLevelWarning  = "warning"
	AlarmLevelCritical = "critical"

	// 告警状态
	AlarmStatusActive   = "active"
	AlarmStatusResolved = "resolved"
	AlarmStatusIgnored  = "ignored"

	// 环境设备类型
	EnvDeviceTypeTemp     = "temperature"
	EnvDeviceTypeHumidity = "humidity"
	EnvDeviceTypeSmoke    = "smoke"
	EnvDeviceTypeWater    = "water"
	EnvDeviceTypeDoor     = "door"
	EnvDeviceTypePower    = "power"

	// 环境设备状态
	EnvDeviceStatusNormal   = "normal"
	EnvDeviceStatusAbnormal = "abnormal"
	EnvDeviceStatusOffline  = "offline"

	// 借还状态
	BorrowStatusBorrowing = "borrowing"
	BorrowStatusReturned  = "returned"
	BorrowStatusOverdue   = "overdue"

	// 系统日志分类
	LogCategoryDevice    = "device"
	LogCategoryUser      = "user"
	LogCategoryAlarm     = "alarm"
	LogCategoryOperation = "operation"

	// 视频通道状态
	VideoStatusOnline  = "online"
	VideoStatusOffline = "offline"

	// 开架记录状态
	OpenRackStatusRecording = "recording"
	OpenRackStatusFinished  = "finished"
)

// User 系统用户
type User struct {
	gorm.Model
	Username string `gorm:"size:64;uniqueIndex" json:"username"`
	Password string `gorm:"size:128" json:"-"`
	Nickname string `gorm:"size:64" json:"nickname"`
	Email    string `gorm:"size:128" json:"email"`
	Phone    string `gorm:"size:32" json:"phone"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Status   int    `gorm:"default:1" json:"status"` // 1启用 0禁用
	RoleCodes string `gorm:"type:text" json:"-"` // JSON 数组：角色编码
	Remark   string `gorm:"type:text" json:"remark"`
}

// Role 角色（不同角色不同权限）
type Role struct {
	gorm.Model
	Code        string `gorm:"size:64;uniqueIndex" json:"code"`
	Name        string `gorm:"size:64" json:"name"`
	Permissions string `gorm:"type:text" json:"permissions"` // JSON 数组：权限标识
	Remark      string `gorm:"type:text" json:"remark"`
	Enabled     int    `gorm:"default:1" json:"enabled"`
}

// Drone 无人机
type Drone struct {
	gorm.Model
	Code          string     `gorm:"size:64;uniqueIndex" json:"code"` // RFID / 电子标签
	Name          string     `gorm:"size:128" json:"name"`
	ModelName     string     `gorm:"size:128" json:"model"`
	Category      string     `gorm:"size:64" json:"category"`
	StorageStatus string     `gorm:"size:32;default:in_position" json:"storageStatus"`
	DeviceStatus  string     `gorm:"size:32;default:intact" json:"deviceStatus"`
	Battery       int        `gorm:"default:100" json:"battery"`
	PositionCode  string     `gorm:"size:64" json:"positionCode"` // 库位编号
	RoomID        uint       `json:"roomId"`
	Manufacturer  string     `gorm:"size:128" json:"manufacturer"`
	PurchaseDate  *time.Time `json:"purchaseDate"`
	ImageURL      string     `gorm:"size:255" json:"imageUrl"`
	Remark        string     `gorm:"type:text" json:"remark"`
}

// StorageRoom 储存室
type StorageRoom struct {
	gorm.Model
	Code         string  `gorm:"size:64;uniqueIndex" json:"code"`
	Name         string  `gorm:"size:128" json:"name"`
	Location     string  `gorm:"size:255" json:"location"`
	Capacity     int     `json:"capacity"` // 最大库位数
	EnvTemp      float64 `json:"envTemp"`
	EnvHumidity  float64 `json:"envHumidity"`
	Status       int     `gorm:"default:1" json:"status"`
	Remark       string  `gorm:"type:text" json:"remark"`
}

// StorageRack 密集架
type StorageRack struct {
	gorm.Model
	RoomID uint   `json:"roomId"`
	Code   string `gorm:"size:64" json:"code"`
	Name   string `gorm:"size:128" json:"name"`
	Rows   int    `json:"rows"`
	Cols   int    `json:"cols"`
	Layers int    `json:"layers"`
	Remark string `gorm:"type:text" json:"remark"`
}

// RackPosition 库位（用于 3D 导航图与图形化盘点）
type RackPosition struct {
	gorm.Model
	RackID uint   `json:"rackId"`
	Code   string `gorm:"size:64" json:"code"` // 例如 A-01-1（列-行-层）
	Row    int    `json:"row"`
	Col    int    `json:"col"`
	Layer  int    `json:"layer"`
	DroneID *uint `json:"droneId"`
	Status int    `gorm:"default:0" json:"status"` // 0空置 1占用
}

// BorrowRecord 借还记录（RFID 扫描触发）
type BorrowRecord struct {
	gorm.Model
	DroneID            uint       `json:"droneId"`
	DroneCode          string     `gorm:"size:64" json:"droneCode"`
	DroneName          string     `gorm:"size:128" json:"droneName"`
	BorrowerID         uint       `json:"borrowerId"`
	BorrowerName       string     `gorm:"size:64" json:"borrowerName"`
	BorrowTime         time.Time  `json:"borrowTime"`
	ExpectedReturnTime *time.Time `json:"expectedReturnTime"`
	ActualReturnTime   *time.Time `json:"actualReturnTime"`
	Status             string     `gorm:"size:32" json:"status"` // borrowing/returned/overdue
	Accessories        string     `gorm:"type:text" json:"accessories"` // JSON 数组：配件
	ReturnCondition    string     `gorm:"size:32" json:"returnCondition"` // 归还时设备状态
	PositionCode       string     `gorm:"size:64" json:"positionCode"`
	DurationMin        int        `json:"durationMin"` // 实时使用时长（分钟）
	Remark             string     `gorm:"type:text" json:"remark"`
}

// AccessRecord 人员进出记录
type AccessRecord struct {
	gorm.Model
	CardNo     string `gorm:"size:64" json:"cardNo"`
	PersonName string `gorm:"size:64" json:"personName"`
	Type       string `gorm:"size:32" json:"type"` // card_auth_fail/door_open/door_close/face_auth_pass
	RoomID     uint   `json:"roomId"`
	DeviceID   uint   `json:"deviceId"`
	ImageURL   string `gorm:"size:255" json:"imageUrl"`
}

// AlarmEvent 告警事件
type AlarmEvent struct {
	gorm.Model
	Type      string     `gorm:"size:32" json:"type"`
	Level     string     `gorm:"size:16;default:warning" json:"level"`
	Title     string     `gorm:"size:255" json:"title"`
	Content   string     `gorm:"type:text" json:"content"`
	Source    string     `gorm:"size:32;default:system" json:"source"` // system/third-party/device
	DeviceID  uint       `json:"deviceId"`
	RoomID    uint       `json:"roomId"`
	Status    string     `gorm:"size:16;default:active" json:"status"`
	ResolvedAt *time.Time `json:"resolvedAt"`
	ResolvedBy string     `gorm:"size:64" json:"resolvedBy"`
}

// VideoChannel 视频通道（库房摄像头）
type VideoChannel struct {
	gorm.Model
	Name   string `gorm:"size:128" json:"name"`
	RoomID uint   `json:"roomId"`
	RackID uint   `json:"rackId"`
	URL    string `gorm:"size:512" json:"url"` // rtsp / hls
	Status string `gorm:"size:16;default:online" json:"status"`
	Remark string `gorm:"type:text" json:"remark"`
}

// OpenRackRecord 开架记录（密集架摄像头录像）
type OpenRackRecord struct {
	gorm.Model
	VideoChannelID uint       `json:"videoChannelId"`
	RackID         uint       `json:"rackId"`
	RackCode       string     `gorm:"size:64" json:"rackCode"`
	Operator       string     `gorm:"size:64" json:"operator"`
	StartTime      time.Time  `json:"startTime"`
	EndTime        *time.Time `json:"endTime"`
	VideoURL       string     `gorm:"size:512" json:"videoUrl"` // 回放地址
	Status         string     `gorm:"size:16;default:recording" json:"status"`
	Remark         string     `gorm:"type:text" json:"remark"`
}

// EnvDevice 环境设备运行状态（温湿度、烟感、水浸、门禁、供电等）
type EnvDevice struct {
	gorm.Model
	Name         string     `gorm:"size:128" json:"name"`
	Type         string     `gorm:"size:32" json:"type"`
	RoomID       uint       `json:"roomId"`
	Status       string     `gorm:"size:16;default:normal" json:"status"`
	Value        float64    `json:"value"`
	Unit         string     `gorm:"size:16" json:"unit"`
	LastReportAt *time.Time `json:"lastReportAt"`
	Remark       string     `gorm:"type:text" json:"remark"`
}

// SystemLog 系统日志（按设备分类管理）
type SystemLog struct {
	gorm.Model
	Category   string `gorm:"size:32" json:"category"` // device/user/alarm/operation
	DeviceType string `gorm:"size:32" json:"deviceType"`
	Content    string `gorm:"type:text" json:"content"`
	Operator   string `gorm:"size:64" json:"operator"`
	Level      string `gorm:"size:16;default:info" json:"level"`
	IP         string `gorm:"size:64" json:"ip"`
}
