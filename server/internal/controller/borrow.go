package controller

import (
	"time"

	"drone-server/internal/db"
	"drone-server/internal/models"
	"drone-server/utils"

	"github.com/gin-gonic/gin"
)

// ScanBorrow 扫描 RFID/电子标签借出无人机
func ScanBorrow(c *gin.Context) {
	var req struct {
		Code              string   `json:"code"` // RFID / 电子标签
		BorrowerName      string   `json:"borrowerName"`
		BorrowerID        uint     `json:"borrowerId"`
		ExpectedReturnHrs int      `json:"expectedReturnHrs"` // 预计归还时长（小时）
		Accessories       []string `json:"accessories"`       // 全套配件
		Remark            string   `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	if req.Code == "" {
		utils.Fail(c, "请扫描无人机标签")
		return
	}
	var d models.Drone
	if err := db.DB.Where("code = ?", req.Code).First(&d).Error; err != nil {
		utils.Fail(c, "未找到该标签对应的无人机")
		return
	}
	if d.StorageStatus != models.DroneStorageStatusInPosition {
		utils.Fail(c, "该无人机当前不在位，无法借出")
		return
	}
	now := time.Now()
	expected := now.Add(time.Duration(req.ExpectedReturnHrs) * time.Hour)
	rec := models.BorrowRecord{
		DroneID:            d.ID,
		DroneCode:          d.Code,
		DroneName:          d.Name,
		BorrowerID:         req.BorrowerID,
		BorrowerName:       req.BorrowerName,
		BorrowTime:         now,
		ExpectedReturnTime: &expected,
		Status:             models.BorrowStatusBorrowing,
		Accessories:        toJSONArray(req.Accessories),
		PositionCode:       d.PositionCode,
		Remark:             req.Remark,
	}
	db.DB.Create(&rec)
	// 更新无人机状态为借出，并释放库位
	db.DB.Model(&d).Updates(map[string]interface{}{"storage_status": models.DroneStorageStatusBorrowed})
	db.DB.Model(&models.RackPosition{}).Where("code = ?", d.PositionCode).
		Updates(map[string]interface{}{"drone_id": nil, "status": 0})
	writeLog(models.LogCategoryOperation, "drone", "借出无人机: "+d.Name+"("+d.Code+") 借用人:"+req.BorrowerName, "info", currentUsername(c), c.ClientIP())
	utils.OK(c, rec)
}

// ScanReturn 扫描 RFID/电子标签归还无人机
func ScanReturn(c *gin.Context) {
	var req struct {
		Code            string `json:"code"`
		ReturnCondition string `json:"returnCondition"` // 归还时设备状态 intact/damaged/scrapped
		Remark          string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	if req.Code == "" {
		utils.Fail(c, "请扫描无人机标签")
		return
	}
	var d models.Drone
	if err := db.DB.Where("code = ?", req.Code).First(&d).Error; err != nil {
		utils.Fail(c, "未找到该标签对应的无人机")
		return
	}
	var rec models.BorrowRecord
	if err := db.DB.Where("drone_id = ? AND status = ?", d.ID, models.BorrowStatusBorrowing).
		Order("id desc").First(&rec).Error; err != nil {
		utils.Fail(c, "未找到该无人机的借出记录")
		return
	}
	now := time.Now()
	condition := req.ReturnCondition
	if condition == "" {
		condition = models.DroneDeviceStatusIntact
	}
	db.DB.Model(&rec).Updates(map[string]interface{}{
		"actual_return_time": now,
		"status":             models.BorrowStatusReturned,
		"return_condition":   condition,
		"duration_min":       int(now.Sub(rec.BorrowTime).Minutes()),
		"remark":             req.Remark,
	})
	// 无人机回到原位，库位重新占用
	db.DB.Model(&d).Updates(map[string]interface{}{
		"storage_status": models.DroneStorageStatusInPosition,
		"device_status":  condition,
	})
	if rec.PositionCode != "" {
		db.DB.Model(&models.RackPosition{}).Where("code = ?", rec.PositionCode).
			Updates(map[string]interface{}{"drone_id": d.ID, "status": 1})
	}
	// 所在区域自动打开（此处记录开门日志，便于与门禁/密集架联动）
	writeLog(models.LogCategoryOperation, "drone", "归还无人机: "+d.Name+"("+d.Code+") 恢复库位:"+rec.PositionCode, "info", currentUsername(c), c.ClientIP())
	utils.OK(c, rec)
}

// ListBorrowRecords 借还记录（支持时间范围、状态、关键字过滤）
func ListBorrowRecords(c *gin.Context) {
	current, size := parsePage(c)
	droneCode := c.Query("droneCode")
	borrowerName := c.Query("borrowerName")
	status := c.Query("status")
	start := c.Query("start")
	end := c.Query("end")

	q := db.DB.Model(&models.BorrowRecord{})
	if droneCode != "" {
		q = q.Where("drone_code LIKE ?", "%"+droneCode+"%")
	}
	if borrowerName != "" {
		q = q.Where("borrower_name LIKE ?", "%"+borrowerName+"%")
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			q = q.Where("borrow_time >= ?", t)
		}
	}
	if end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			q = q.Where("borrow_time <= ?", t)
		}
	}
	var total int64
	q.Count(&total)
	var records []models.BorrowRecord
	q.Order("id desc").Offset((current - 1) * size).Limit(size).Find(&records)
	for i := range records {
		refreshBorrowStatus(&records[i])
	}
	utils.OK(c, utils.Page(records, current, size, total))
}

// GetBorrowRecord 借还记录详情
func GetBorrowRecord(c *gin.Context) {
	var rec models.BorrowRecord
	if err := db.DB.First(&rec, atoiParam(c, "id")).Error; err != nil {
		utils.Fail(c, "记录不存在")
		return
	}
	refreshBorrowStatus(&rec)
	utils.OK(c, rec)
}

// refreshBorrowStatus 计算实时状态与时长
func refreshBorrowStatus(rec *models.BorrowRecord) {
	if rec.Status == models.BorrowStatusReturned {
		if rec.ActualReturnTime != nil {
			rec.DurationMin = int(rec.ActualReturnTime.Sub(rec.BorrowTime).Minutes())
		}
		return
	}
	if rec.ExpectedReturnTime != nil && time.Now().After(*rec.ExpectedReturnTime) {
		rec.Status = models.BorrowStatusOverdue
	}
	rec.DurationMin = int(time.Now().Sub(rec.BorrowTime).Minutes())
}
