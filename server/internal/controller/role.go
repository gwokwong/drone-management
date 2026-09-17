package controller

import (
	"drone-server/internal/db"
	"drone-server/internal/models"
	"drone-server/utils"

	"github.com/gin-gonic/gin"
)

// ListRoles 角色列表
func ListRoles(c *gin.Context) {
	current, size := parsePage(c)
	name := c.Query("name")
	code := c.Query("code")

	q := db.DB.Model(&models.Role{})
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		q = q.Where("code LIKE ?", "%"+code+"%")
	}
	var total int64
	q.Count(&total)
	var roles []models.Role
	q.Order("id desc").Offset((current - 1) * size).Limit(size).Find(&roles)
	utils.OK(c, utils.Page(roles, current, size, total))
}

// GetRole 角色详情
func GetRole(c *gin.Context) {
	var r models.Role
	if err := db.DB.First(&r, atoiParam(c, "id")).Error; err != nil {
		utils.Fail(c, "角色不存在")
		return
	}
	utils.OK(c, r)
}

// CreateRole 新增角色
func CreateRole(c *gin.Context) {
	var r models.Role
	if err := c.ShouldBindJSON(&r); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	if r.Code == "" || r.Name == "" {
		utils.Fail(c, "角色编码与名称不能为空")
		return
	}
	var cnt int64
	db.DB.Model(&models.Role{}).Where("code = ?", r.Code).Count(&cnt)
	if cnt > 0 {
		utils.Fail(c, "角色编码已存在")
		return
	}
	if r.Enabled == 0 {
		r.Enabled = 1
	}
	db.DB.Create(&r)
	utils.OK(c, r)
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
	id := atoiParam(c, "id")
	var req struct {
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
		Remark      string   `json:"remark"`
		Enabled     int      `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误")
		return
	}
	var r models.Role
	if err := db.DB.First(&r, id).Error; err != nil {
		utils.Fail(c, "角色不存在")
		return
	}
	r.Name = req.Name
	r.Permissions = toJSONArray(req.Permissions)
	r.Remark = req.Remark
	r.Enabled = req.Enabled
	db.DB.Save(&r)
	utils.OK(c, r)
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	db.DB.Delete(&models.Role{}, atoiParam(c, "id"))
	utils.OK(c, nil)
}

// ListPermissions 全部角色（用于分配）
func AllRoles(c *gin.Context) {
	var roles []models.Role
	db.DB.Find(&roles)
	utils.OK(c, roles)
}
