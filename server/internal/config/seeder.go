package config

import (
	"log"
	"hypermenu-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdminSuperuser() {
	var count int64

	if err := DB.Model(&models.Admin{}).Count(&count).Error; err != nil {
		log.Printf("خطا در بررسی ادمین‌ها: %v", err)
		return
	}

	if count > 0 {
		log.Println("حساب ادمین از قبل موجود است.")
		return
	}

	// برای توسعه مناسب است؛ پیش از انتشار، این رمز را عوض کنید.
	const defaultPassword = "admin123456"

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(defaultPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("خطا در هش‌کردن رمز ادمین: %v", err)
		return
	}

	admin := models.Admin{
		Username: "admin",
		Password: string(hashedPassword),
	}

	if err := DB.Create(&admin).Error; err != nil {
		log.Printf("خطا در ساخت ادمین: %v", err)
		return
	}

	log.Println("ادمین اولیه ساخته شد؛ نام کاربری: admin")
}
