package businessowner

import (
	"github.com/gin-gonic/gin"
	"hypermenu-backend/internal/handlers/businessOwner"
	"hypermenu-backend/internal/middleware"
)

func SetupBusinessOwnerRoutes(r *gin.Engine) {
	businessOwnerGroup := r.Group("/business-owners")
	{

		businessOwnerGroup.POST("/auth/register", businessowner.RegisterBusiness)
		businessOwnerGroup.POST("/auth/login", businessowner.LoginBusinessAdmin)
	}

	// گروه روت‌های محافظت‌شده
	protected := businessOwnerGroup.Group("/dashboard")
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
