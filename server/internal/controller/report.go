package controller

import (
	"time"

	"drone-server/internal/db"
	"drone-server/internal/models"
	"drone-server/utils"

	"github.com/gin-gonic/gin"
)

// reportPeriod 依据报表类型计算统计周期
func reportPeriod(reportType string) (start, end time.Time) {
	now := time.Now()
	end = now
	switch reportType {
	case "daily":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "weekly":
		start = now.AddDate(0, 0, -7)
	case "monthly":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "yearly":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	default: // 默认近一天
		start = now.AddDate(0, 0, -1)
	}
	return
}

// GenerateReport 生成日/周/月/年报表数据（数据与图形共用）
func GenerateReport(c *gin.Context) {
	reportType := c.DefaultQuery("type", "daily")
	start, end := reportPeriod(reportType)

	var droneTotal, droneBorrowed int64
	db.DB.Model(&models.Drone{}).Count(&droneTotal)
	db.DB.Model(&models.Drone{}).Where("storage_status = ?", models.DroneStorageStatusBorrowed).Count(&droneBorrowed)

	var borrowCount, overdueCount int64
	db.DB.Model(&models.BorrowRecord{}).Where("borrow_time >= ? AND borrow_time <= ?", start, end).Count(&borrowCount)
	db.DB.Model(&models.BorrowRecord{}).Where("status = ?", models.BorrowStatusOverdue).Count(&overdueCount)

	var alarmCount int64
	db.DB.Model(&models.AlarmEvent{}).Where("created_at >= ? AND created_at <= ?", start, end).Count(&alarmCount)

	var accessCount int64
	db.DB.Model(&models.AccessRecord{}).Where("created_at >= ? AND created_at <= ?", start, end).Count(&accessCount)

	// 按日趋势（用于折线图）
	type trendItem struct {
		Date  string `json:"date"`
		Borrow int64 `json:"borrow"`
		Alarm  int64 `json:"alarm"`
		Access int64 `json:"access"`
	}
	var trend []trendItem
	days := 0
	switch reportType {
	case "daily":
		days = 1
	case "weekly":
		days = 7
	case "monthly":
		days = 30
	case "yearly":
		days = 12
	}
	if reportType == "yearly" {
		for m := 1; m <= 12; m++ {
			ms := time.Date(start.Year(), time.Month(m), 1, 0, 0, 0, 0, start.Location())
			me := ms.AddDate(0, 1, 0)
			var b, a, ac int64
			db.DB.Model(&models.BorrowRecord{}).Where("borrow_time >= ? AND borrow_time < ?", ms, me).Count(&b)
			db.DB.Model(&models.AlarmEvent{}).Where("created_at >= ? AND created_at < ?", ms, me).Count(&a)
			db.DB.Model(&models.AccessRecord{}).Where("created_at >= ? AND created_at < ?", ms, me).Count(&ac)
			trend = append(trend, trendItem{Date: ms.Format("01月"), Borrow: b, Alarm: a, Access: ac})
		}
	} else {
		for i := days - 1; i >= 0; i-- {
			d := start.AddDate(0, 0, i)
			ds := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
			de := ds.AddDate(0, 0, 1)
			var b, a, ac int64
			db.DB.Model(&models.BorrowRecord{}).Where("borrow_time >= ? AND borrow_time < ?", ds, de).Count(&b)
			db.DB.Model(&models.AlarmEvent{}).Where("created_at >= ? AND created_at < ?", ds, de).Count(&a)
			db.DB.Model(&models.AccessRecord{}).Where("created_at >= ? AND created_at < ?", ds, de).Count(&ac)
			trend = append(trend, trendItem{Date: ds.Format("01-02"), Borrow: b, Alarm: a, Access: ac})
		}
	}

	utils.OK(c, gin.H{
		"type":          reportType,
		"start":         start,
		"end":           end,
		"droneTotal":    droneTotal,
		"droneBorrowed": droneBorrowed,
		"borrowCount":   borrowCount,
		"overdueCount":  overdueCount,
		"alarmCount":    alarmCount,
		"accessCount":   accessCount,
		"trend":         trend,
	})
}

// ExportReport 报表导出为 Excel
func ExportReport(c *gin.Context) {
	reportType := c.DefaultQuery("type", "daily")

	// 直接基于库统计，便于导出
	start, end := reportPeriod(reportType)
	var droneTotal, droneBorrowed, borrowCount, alarmCount, accessCount int64
	db.DB.Model(&models.Drone{}).Count(&droneTotal)
	db.DB.Model(&models.Drone{}).Where("storage_status = ?", models.DroneStorageStatusBorrowed).Count(&droneBorrowed)
	db.DB.Model(&models.BorrowRecord{}).Where("borrow_time >= ? AND borrow_time <= ?", start, end).Count(&borrowCount)
	db.DB.Model(&models.AlarmEvent{}).Where("created_at >= ? AND created_at <= ?", start, end).Count(&alarmCount)
	db.DB.Model(&models.AccessRecord{}).Where("created_at >= ? AND created_at <= ?", start, end).Count(&accessCount)

	headers := []string{"指标", "数值"}
	rows := [][]interface{}{
		{"报表类型", reportType},
		{"统计开始", start.Format(time.RFC3339)},
		{"统计结束", end.Format(time.RFC3339)},
		{"无人机总数", droneTotal},
		{"借出中数量", droneBorrowed},
		{"周期内借还次数", borrowCount},
		{"周期内告警数", alarmCount},
		{"周期内进出人次", accessCount},
	}
	exportExcel(c, "report_"+reportType+".xlsx", headers, rows)
}
