package router

import (
	"drone-server/internal/controller"
	"drone-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 注册全部路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		// 公开接口
		api.POST("/auth/login", controller.Login)
		api.POST("/auth/logout", controller.Logout)

		// 需鉴权接口
		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/user/info", controller.GetUserInfo)

			// 系统管理
			auth.GET("/users", controller.ListUsers)
			auth.GET("/users/:id", controller.GetUser)
			auth.POST("/users", controller.CreateUser)
			auth.PUT("/users/:id", controller.UpdateUser)
			auth.DELETE("/users/:id", controller.DeleteUser)
			auth.POST("/users/:id/reset-password", controller.ResetPassword)

			auth.GET("/roles", controller.ListRoles)
			auth.GET("/roles/all", controller.AllRoles)
			auth.GET("/roles/:id", controller.GetRole)
			auth.POST("/roles", controller.CreateRole)
			auth.PUT("/roles/:id", controller.UpdateRole)
			auth.DELETE("/roles/:id", controller.DeleteRole)

			// 无人机管理
			auth.GET("/drones", controller.ListDrones)
			auth.GET("/drones/stats", controller.DroneStats)
			auth.GET("/drones/:id", controller.GetDrone)
			auth.POST("/drones", controller.CreateDrone)
			auth.PUT("/drones/:id", controller.UpdateDrone)
			auth.DELETE("/drones/:id", controller.DeleteDrone)

			// 借还（RFID 扫描）
			auth.POST("/borrow/scan", controller.ScanBorrow)
			auth.POST("/return/scan", controller.ScanReturn)
			auth.GET("/borrow/records", controller.ListBorrowRecords)
			auth.GET("/borrow/records/:id", controller.GetBorrowRecord)

			// 储存室 / 密集架 / 3D 导航 / 环境设备
			auth.GET("/storage-rooms", controller.ListRooms)
			auth.GET("/storage-rooms/:id", controller.GetRoom)
			auth.GET("/storage-rooms/:id/map", controller.RoomMap3D)
			auth.POST("/storage-rooms", controller.CreateRoom)
			auth.PUT("/storage-rooms/:id", controller.UpdateRoom)
			auth.DELETE("/storage-rooms/:id", controller.DeleteRoom)

			auth.GET("/racks", controller.ListRacks)
			auth.POST("/racks", controller.CreateRack)
			auth.PUT("/racks/:id", controller.UpdateRack)
			auth.DELETE("/racks/:id", controller.DeleteRack)

			auth.GET("/env-devices", controller.ListEnvDevices)
			auth.GET("/env-devices/summary", controller.EnvDeviceSummary)
			auth.POST("/env-devices", controller.CreateEnvDevice)
			auth.PUT("/env-devices/:id", controller.UpdateEnvDevice)
			auth.DELETE("/env-devices/:id", controller.DeleteEnvDevice)

			// 告警
			auth.GET("/alarms", controller.ListAlarms)
			auth.GET("/alarms/stats", controller.AlarmStats)
			auth.POST("/alarms", controller.CreateAlarm)
			auth.POST("/alarms/:id/resolve", controller.ResolveAlarm)
			auth.POST("/alarms/:id/ignore", controller.IgnoreAlarm)
			auth.DELETE("/alarms/:id", controller.DeleteAlarm)

			// 人员进出
			auth.GET("/access", controller.ListAccess)
			auth.GET("/access/stats", controller.AccessStats)
			auth.POST("/access", controller.CreateAccess)

			// 视频 / 开架
			auth.GET("/videos", controller.ListVideos)
			auth.POST("/videos", controller.CreateVideo)
			auth.PUT("/videos/:id", controller.UpdateVideo)
			auth.DELETE("/videos/:id", controller.DeleteVideo)
			auth.POST("/videos/:id/open-rack", controller.OpenRack)
			auth.POST("/open-racks/:id/stop", controller.StopOpenRack)
			auth.GET("/open-racks", controller.ListOpenRacks)

			// 报表 / 日志 / 大屏
			auth.GET("/reports/generate", controller.GenerateReport)
			auth.GET("/reports/export", controller.ExportReport)
			auth.GET("/logs", controller.ListLogs)
			auth.GET("/logs/export", controller.ExportLogs)
			auth.GET("/dashboard/overview", controller.DashboardOverview)

			// 第三方对接接口（预留）
			auth.POST("/integration/alarm/push", controller.IntegrationAlarmPush)
			auth.POST("/integration/access/push", controller.IntegrationAccessPush)
			auth.POST("/integration/env/push", controller.IntegrationEnvPush)
			auth.POST("/integration/rfid/borrow", controller.IntegrationRFIDBorrow)
			auth.POST("/integration/rfid/return", controller.IntegrationRFIDReturn)
		}
	}
	return r
}
