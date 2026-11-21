package businessowner

import (
	"net/http"
	"time"

	"hypermenu-backend/internal/config"
	"hypermenu-backend/internal/models"
	"hypermenu-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ساختار داده ورودی برای ثبت‌نام کافه و مالک آن
type RegisterBusinessRequest struct {
	BusinessName string `json:"business_name" binding:"required,min=3,max=150"`
	Slug         string `json:"slug" binding:"required,alphanum,min=3,max=100"`
	Phone        string `json:"phone" binding:"required"`
	Address      string `json:"address"`

	// مشخصات مالک / ادمین
	FullName string `json:"full_name" binding:"required,min=3"`
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

// ساختار داده ورودی برای لاگین
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterBusiness هندلر ثبت نام کسب‌وکار جدید
func RegisterBusiness(c *gin.Context) {
	var req RegisterBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "اطلاعات ورودی نامعتبر است", "INVALID_INPUT")
		return
	}

	// ۱. بررسی تکراری نبودن اسلاگ کسب‌وکار
	var existingBusiness models.Business
	if err := config.DB.Where("slug = ?", req.Slug).First(&existingBusiness).Error; err == nil {
		utils.JSONError(c, http.StatusConflict, "این نام انگلیسی (اسلاگ) قبلاً ثبت شده است", "SLUG_EXISTS")
		return
	}

	// ۲. بررسی تکراری نبودن نام کاربری ادمین
	var existingAdmin models.BusinessAdmin
	if err := config.DB.Where("username = ?", req.Username).First(&existingAdmin).Error; err == nil {
		utils.JSONError(c, http.StatusConflict, "این نام کاربری قبلاً استفاده شده است", "USERNAME_EXISTS")
		return
	}

	// ۳. هش کردن پسورد
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "خطا در پردازش رمز عبور", "SERVER_ERROR")
		return
	}

	// ۴. ثبت در قالب یک Transaction (تا اگر یکی ناموفق بود، هیچ‌کدام ثبت نشود)
	tx := config.DB.Begin()

	// مهلت تست رایگان ۱۴ روزه
	trialExpire := time.Now().AddDate(0, 0, 14)
	business := models.Business{
		Name:                  req.BusinessName,
		Slug:                  req.Slug,
		Phone:                 req.Phone,
		Address:               req.Address,
		Plan:                  "trial",
		SubscriptionExpiresAt: &trialExpire,
		IsActive:              true,
	}

	if err := tx.Create(&business).Error; err != nil {
		tx.Rollback()
		utils.JSONError(c, http.StatusInternalServerError, "خطا در ثبت کسب‌وکار", "DB_ERROR")
		return
	}

	admin := models.BusinessAdmin{
		BusinessID:   business.ID,
		FullName:     req.FullName,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         "owner",
		IsActive:     true,
	}

	if err := tx.Create(&admin).Error; err != nil {
		tx.Rollback()
		utils.JSONError(c, http.StatusInternalServerError, "خطا در ثبت اطلاعات مدیر", "DB_ERROR")
		return
	}

	tx.Commit()

	// ۵. تولید توکن‌های دسترسی
	accessToken, refreshToken, err := utils.GenerateTokens(admin.ID, admin.Role)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "خطا در ایجاد توکن دسترسی", "TOKEN_ERROR")
		return
	}

	utils.JSONSuccess(c, http.StatusCreated, "ثبت‌نام با موفقیت انجام شد", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"business": gin.H{
			"id":   business.ID,
			"name": business.Name,
			"slug": business.Slug,
			"plan": business.Plan,
		},
		"admin": gin.H{
			"id":        admin.ID,
			"full_name": admin.FullName,
			"username":  admin.Username,
			"role":      admin.Role,
		},
	})
}

// LoginBusinessAdmin هندلر ورود ادمین کافه
func LoginBusinessAdmin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "نام کاربری و رمز عبور الزامی است", "INVALID_INPUT")
		return
	}

	var admin models.BusinessAdmin
	if err := config.DB.Preload("Business").Where("username = ?", req.Username).First(&admin).Error; err != nil {
		utils.JSONError(c, http.StatusUnauthorized, "نام کاربری یا رمز عبور اشتباه است", "AUTH_FAILED")
		return
	}

	if !admin.IsActive {
		utils.JSONError(c, http.StatusForbidden, "حساب کاربری شما مسدود شده است", "ACCOUNT_DISABLED")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		utils.JSONError(c, http.StatusUnauthorized, "نام کاربری یا رمز عبور اشتباه است", "AUTH_FAILED")
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(admin.ID, admin.Role)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "خطا در ایجاد توکن", "TOKEN_ERROR")
		return
	}

	utils.JSONSuccess(c, http.StatusOK, "ورود با موفقیت انجام شد", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"admin": gin.H{
			"id":        admin.ID,
			"full_name": admin.FullName,
			"username":  admin.Username,
			"role":      admin.Role,
		},
		"business": gin.H{
			"id":                      admin.Business.ID,
			"name":                    admin.Business.Name,
			"slug":                    admin.Business.Slug,
			"plan":                    admin.Business.Plan,
			"subscription_expires_at": admin.Business.SubscriptionExpiresAt,
			"is_active":               admin.Business.IsActive,
		},
	})
}
