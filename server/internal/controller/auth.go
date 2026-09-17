package controller

import (
	"drone-server/internal/db"
	"drone-server/internal/models"
	"drone-server/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求
type LoginRequest struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	var user models.User
	if err := db.DB.Where("username = ?", req.UserName).First(&user).Error; err != nil {
		utils.Fail(c, "用户不存在")
		return
	}
	if user.Status != 1 {
		utils.Fail(c, "账号已被禁用")
		return
	}
	if !utils.CheckPassword(req.Password, user.Password) {
		utils.Fail(c, "密码错误")
		return
	}
	roles := parseJSONArray(user.RoleCodes)
	token, err := utils.GenerateToken(user.ID, user.Username, roles)
	if err != nil {
		utils.Fail(c, "令牌生成失败")
		return
	}
	refresh, _ := utils.GenerateToken(user.ID, user.Username, roles)
	writeLog(models.LogCategoryUser, "", "用户登录: "+user.Username, "info", user.Username, c.ClientIP())
	utils.OK(c, gin.H{"token": token, "refreshToken": refresh})
}

// Logout 退出登录
func Logout(c *gin.Context) {
	utils.OK(c, nil)
}

// GetUserInfo 获取当前用户信息
func GetUserInfo(c *gin.Context) {
	uid := currentUserID(c)
	var user models.User
	if err := db.DB.First(&user, uid).Error; err != nil {
		utils.Unauthorized(c, "用户不存在")
		return
	}
	roles := parseJSONArray(user.RoleCodes)
	utils.OK(c, gin.H{
		"userId":   user.ID,
		"userName": user.Username,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"nickname": user.Nickname,
		"roles":    roles,
		"buttons":  []string{},
	})
}
