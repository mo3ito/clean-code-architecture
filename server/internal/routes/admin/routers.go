package admin

import (
	"hypermenu-backend/internal/handlers/admin"
	"hypermenu-backend/internal/middleware" 
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(r *gin.Engine) {
	adminGroup := r.Group("/admin")
	{
		
		adminGroup.POST("/auth/login", admin.Login)
		adminGroup.POST("/auth/refresh_token", admin.RefreshToken)
	}

	// گروه روت‌های محافظت‌شده
	protected := adminGroup.Group("/dashboard")
	protected.Use(middleware.AuthMiddleware()) // <--- اعمال میدلور
	{
		// حالا هر چیزی که اینجا اضافه کنید نیاز به توکن دارد
		protected.GET("/", func(c *gin.Context) {
			// می‌توانید اطلاعات ادمین را از کانتکست بگیرید
			username, _ := c.Get("username")
			c.JSON(200, gin.H{
				"message": "خوش آمدید به پنل مدیریت",
				"admin":   username,
			})
		})
	}
}
