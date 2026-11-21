package config

import (
	"fmt"
	"log"
	"os"

	"hypermenu-backend/internal/models" 
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("هشدار: فایل .env پیدا نشد.")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("خطا در اتصال به پایگاه داده: ", err)
	}

	DB = database
	fmt.Println("اتصال به دیتابیس با موفقیت برقرار شد!")

	// 1. اجرای خودکار مهاجرت (ساخت جدول ادمین اگر وجود ندارد)
	if err := DB.AutoMigrate(&models.Admin{}); err != nil {
		log.Fatal("خطا در اجرای Migrations: ", err)
	}

	// 2. فراخوانی Seeder برای ساخت ادمین پیش‌فرض
	SeedAdminSuperuser() 
}
