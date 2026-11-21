package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims ساختار داده‌هایی که در توکن ذخیره می‌شوند
type Claims struct {
	AdminID  uint   `json:"admin_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// getAccessSecret دریافت کلید امضای Access Token از env
func getAccessSecret() []byte {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		return []byte("my_super_secret_access_key_change_me_in_prod")
	}
	return secret
}

// getRefreshSecret دریافت کلید امضای Refresh Token از env
func getRefreshSecret() []byte {
	secret := []byte(os.Getenv("JWT_REFRESH_SECRET"))
	if len(secret) == 0 {
		return []byte("my_super_secret_refresh_key_change_me_in_prod")
	}
	return secret
}

// GenerateTokens تولید همزمان Access Token و Refresh Token
func GenerateTokens(adminID uint, username string) (string, string, error) {
	now := time.Now()

	// ۱. توکن دسترسی (Access Token) - انقضا معمولاً کوتاه‌مدت، مثلاً ۲ ساعت
	accessExpiration := now.Add(2 * time.Hour)
	accessClaims := &Claims{
		AdminID:  adminID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiration),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(getAccessSecret())
	if err != nil {
		return "", "", err
	}

	// ۲. توکن بازنشانی (Refresh Token) - انقضا بلندمدت، مثلاً ۷ روز
	refreshExpiration := now.Add(7 * 24 * time.Hour)
	refreshClaims := &Claims{
		AdminID:  adminID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiration),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(getRefreshSecret())
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// ValidateRefreshToken اعتبارسنجی و خواندن داده‌های رفرش‌توکن
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("روش امضای توکن نامعتبر است")
		}
		return getRefreshSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("توکن نامعتبر است")
}
