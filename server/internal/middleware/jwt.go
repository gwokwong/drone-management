package middleware

import (
	"strings"

	"drone-server/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuth 鉴权中间件，从 Authorization 头解析令牌并写入上下文
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			utils.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}
		auth = strings.TrimPrefix(auth, "Bearer ")
		claims, err := utils.ParseToken(auth)
		if err != nil {
			utils.Unauthorized(c, "令牌无效或已过期")
			c.Abort()
			return
		}
		c.Set("uid", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

// RequireRoles 角色权限校验中间件
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("roles")
		if !ok {
			utils.Forbidden(c, "无访问权限")
			c.Abort()
			return
		}
		userRoles, ok := v.([]string)
		if !ok {
			utils.Forbidden(c, "无访问权限")
			c.Abort()
			return
		}
		for _, r := range userRoles {
			for _, need := range roles {
				if r == need {
					c.Next()
					return
				}
			}
		}
		utils.Forbidden(c, "无访问权限")
		c.Abort()
	}
}
