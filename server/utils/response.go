package utils

import "github.com/gin-gonic/gin"

// Response 统一接口响应结构，与前端 BaseResponse<T> 对齐：{code,msg,data}
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// OK 业务成功（HTTP 200，code=200）
func OK(c *gin.Context, data interface{}) {
	c.JSON(200, Response{Code: 200, Msg: "success", Data: data})
}

// Fail 业务失败（HTTP 200，code=400，前端按 code 判断）
func Fail(c *gin.Context, msg string) {
	c.JSON(200, Response{Code: 400, Msg: msg, Data: nil})
}

// Unauthorized 未授权（HTTP 401）
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(401, Response{Code: 401, Msg: msg, Data: nil})
}

// Forbidden 无权限（HTTP 403）
func Forbidden(c *gin.Context, msg string) {
	c.JSON(403, Response{Code: 403, Msg: msg, Data: nil})
}

// Page 分页响应结构，与前端 Api.Common.PaginatedResponse 对齐
func Page(records interface{}, current, size int, total int64) gin.H {
	return gin.H{
		"records": records,
		"current": current,
		"size":    size,
		"total":   total,
	}
}
