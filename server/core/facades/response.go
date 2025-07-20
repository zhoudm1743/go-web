package facades

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`            // 状态码
	Message string      `json:"message"`         // 提示消息
	Data    interface{} `json:"data"`            // 响应数据
	Error   string      `json:"error,omitempty"` // 错误信息
}

// ResponseWithSuccess 成功响应
func ResponseWithSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// ResponseWithError 错误响应
func ResponseWithError(c *gin.Context, code int, message string, err error) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Error:   errMsg,
	})
}

// ResponseWithPage 分页响应
func ResponseWithPage(c *gin.Context, message string, data interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": message,
		"data":    data,
		"pagination": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}
