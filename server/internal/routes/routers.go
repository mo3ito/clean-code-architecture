package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"

	// ایمپورت فایلِ روت‌های ادمین
	"hypermenu-backend/internal/routes/admin"
	businessowner "hypermenu-backend/internal/routes/businessOwner"
)

func SetupRoutes(r *gin.Engine) {
	// فراخوانی تابع ادمین و پاس دادنِ Engine اصلی
	admin.SetupAdminRoutes(r)
	businessowner.SetupBusinessOwnerRoutes(r)

	// برای اطمینان از ثبت شدن، روت‌ها را چاپ کن
	fmt.Println("Registered Routes:")
	for _, route := range r.Routes() {
		fmt.Printf("%s\t%s\n", route.Method, route.Path)
	}
}
