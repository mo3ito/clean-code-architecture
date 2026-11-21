package main

import (
	"fmt"
	"hypermenu-backend/internal/config"
	"hypermenu-backend/internal/models"
	routers "hypermenu-backend/internal/routes"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// ۱. اتصال به دیتابیس
	config.ConnectDatabase()

	// ۲. اجرای مایگریشن جداول با استفاده از config.DB
	err := config.DB.AutoMigrate(
		&models.Business{},
		&models.BusinessAdmin{},
	)
	if err != nil {
		log.Fatalf("خطا در ایجاد جداول دیتابیس: %v", err)
	}

	// ۳. راه‌اندازی Gin
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ۴. ثبت مسیرها
	routers.SetupRoutes(r)

	// مسیر تست سلامت سیستم
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "سیستم هایپر منو فعال است 🚀",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port:", port)
	r.Run(":" + port)
}
