package admin

import (
	"errors"
	"net/http"

	"hypermenu-backend/internal/config"
	"hypermenu-backend/internal/models"
	"hypermenu-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginRequest ساختار ورودی لاگین
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest ساختار ورودی تمدید توکن
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse ساختار داده‌های خروجی توکن‌ها
type AuthResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	Admin        *AdminInfo `json:"admin,omitempty"`
}

type AdminInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// Login هندلر ورود ادمین
func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "نام کاربری و رمز عبور الزامی است", "INVALID_CREDENTIALS_PAYLOAD")
		return
	}

	var admin models.Admin
	err := config.DB.Where("username = ?", req.Username).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONError(c, http.StatusUnauthorized, "نام کاربری یا رمز عبور اشتباه است", "USER_NOT_FOUND")
			return
		}

		utils.JSONError(c, http.StatusInternalServerError, "خطای سرور در بازیابی اطلاعات ادمین", "DATABASE_ERROR")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	if err != nil {
		utils.JSONError(c, http.StatusUnauthorized, "نام کاربری یا رمز عبور اشتباه است", "INVALID_PASSWORD")
		return
	}

	// تولید همزمان Access Token و Refresh Token
	accessToken, refreshToken, err := utils.GenerateTokens(admin.ID, admin.Username)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "خطا در تولید توکن‌های دسترسی", "TOKEN_GENERATION_FAILED")
		return
	}

	utils.JSONSuccess(c, http.StatusOK, "ورود با موفقیت انجام شد", AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Admin: &AdminInfo{
			ID:       admin.ID,
			Username: admin.Username,
		},
	})
}

// RefreshToken هندلر تمدید نشست با رفرش توکن
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "ارسال refresh_token الزامی است", "REFRESH_TOKEN_REQUIRED")
		return
	}

	// ۱. اعتبارسنجی و دیکود کردن توکن
	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		utils.JSONError(c, http.StatusUnauthorized, "رفرش توکن نامعتبر یا منقضی شده است", "INVALID_REFRESH_TOKEN")
		return
	}

	// ۲. اطمینان از وجود و فعال بودن هنوزِ این ادمین در دیتابیس
	var admin models.Admin
	if err := config.DB.First(&admin, claims.AdminID).Error; err != nil {
		utils.JSONError(c, http.StatusUnauthorized, "حساب کاربری یافت نشد یا غیرفعال شده است", "ADMIN_NOT_FOUND")
		return
	}

	// ۳. تولید جفت توکن جدید (Token Rotation برای امنیت بالاتر)
	newAccessToken, newRefreshToken, err := utils.GenerateTokens(admin.ID, admin.Username)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "خطا در صدور توکن‌های جدید", "TOKEN_REFRESH_FAILED")
		return
	}

	// ۴. بازگرداندن توکن‌های جدید
	utils.JSONSuccess(c, http.StatusOK, "توکن با موفقیت تمدید شد", AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
}
