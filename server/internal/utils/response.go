package utils

import (
	"github.com/gin-gonic/gin"
)

// APIResponse ساختار استاندارد برای تمام پاسخ‌های API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"` // کد خطا یا متن خطا
}

// Success برای پاسخ‌های موفق
func JSONSuccess(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// JSONError برای پاسخ‌های خطا
func JSONError(c *gin.Context, status int, message string, errorCode string) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Error:   errorCode,
	})
}
