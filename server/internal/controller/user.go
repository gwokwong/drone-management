package controller

import (
	"drone-server/internal/db"
	"drone-server/internal/models"
	"drone-server/utils"

	"github.com/gin-gonic/gin"
)

// ListUsers 用户列表（分页 + 过滤）
func ListUsers(c *gin.Context) {
	current, size := parsePage(c)
	username := c.Query("userName")
	status := c.Query("status")
	nickname := c.Query("nickName")

	q := db.DB.Model(&models.User{})
	if username != "" {
		q = q.Where("username LIKE ?", "%"+username+"%")
	}
	if nickname != "" {
		q = q.Where("nickname LIKE ?", "%"+nickname+"%")
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	q.Count(&total)
	var users []models.User
	q.Order("id desc").Offset((current - 1) * size).Limit(size).Find(&users)

	// models.User.RoleCodes 标记为 json:"-"，转视图对象回传角色编码
	type userView struct {
		models.User
		RoleCodes []string `json:"roleCodes"`
	}
	views := make([]userView, 0, len(users))
	for _, u := range users {
		views = append(views, userView{User: u, RoleCodes: parseJSONArray(u.RoleCodes)})
	}
	utils.OK(c, utils.Page(views, current, size, total))
}

// GetUser 用户详情
func GetUser(c *gin.Context) {
	var u models.User
	if err := db.DB.First(&u, atoiParam(c, "id")).Error; err != nil {
		utils.Fail(c, "用户不存在")
		return
	}
	utils.OK(c, u)
}

// CreateUser 新增用户
func CreateUser(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	// Password / RoleCodes 在 models.User 上标记为 json:"-"，需单独解析
	var ext struct {
		Password  string   `json:"password"`
		RoleCodes []string `json:"roleCodes"`
	}
	_ = c.ShouldBindJSON(&ext)

	if u.Username == "" {
		utils.Fail(c, "用户名不能为空")
		return
	}
	var cnt int64
	db.DB.Model(&models.User{}).Where("username = ?", u.Username).Count(&cnt)
	if cnt > 0 {
		utils.Fail(c, "用户名已存在")
		return
	}
	if ext.Password == "" {
		u.Password = utils.HashPassword("123456")
	} else {
		u.Password = utils.HashPassword(ext.Password)
	}
	u.RoleCodes = toJSONArray(ext.RoleCodes)
	if u.Status == 0 {
		u.Status = 1
	}
	if err := db.DB.Create(&u).Error; err != nil {
		utils.Fail(c, "创建失败")
		return
	}
	writeLog(models.LogCategoryUser, "", "新增用户: "+u.Username, "info", currentUsername(c), c.ClientIP())
	utils.OK(c, u)
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	id := atoiParam(c, "id")
	var req struct {
		Nickname  string   `json:"nickname"`
		Email     string   `json:"email"`
		Phone     string   `json:"phone"`
		Avatar    string   `json:"avatar"`
		Status    int      `json:"status"`
		RoleCodes []string `json:"roleCodes"`
		Remark    string   `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	var u models.User
	if err := db.DB.First(&u, id).Error; err != nil {
		utils.Fail(c, "用户不存在")
		return
	}
	u.Nickname = req.Nickname
	u.Email = req.Email
	u.Phone = req.Phone
	u.Avatar = req.Avatar
	u.Status = req.Status
	u.RoleCodes = toJSONArray(req.RoleCodes)
	u.Remark = req.Remark
	db.DB.Save(&u)
	writeLog(models.LogCategoryUser, "", "更新用户: "+u.Username, "info", currentUsername(c), c.ClientIP())
	utils.OK(c, u)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id := atoiParam(c, "id")
	if err := db.DB.Delete(&models.User{}, id).Error; err != nil {
		utils.Fail(c, "删除失败")
		return
	}
	writeLog(models.LogCategoryUser, "", "删除用户 ID", "info", currentUsername(c), c.ClientIP())
	utils.OK(c, nil)
}

// ResetPassword 重置密码
func ResetPassword(c *gin.Context) {
	id := atoiParam(c, "id")
	var req struct {
		Password string `json:"password"`
	}
	_ = c.ShouldBindJSON(&req)
	pw := req.Password
	if pw == "" {
		pw = "123456"
	}
	db.DB.Model(&models.User{}).Where("id = ?", id).Update("password", utils.HashPassword(pw))
	utils.OK(c, nil)
}
