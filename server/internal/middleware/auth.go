package middleware

import (
	"net/http"
	"os"
	"strings"
	"hypermenu-backend/internal/utils" // برای دسترسی به Claims
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "توکن احراز هویت الزامی است"})
			c.Abort()
			return
		}

		// جدا کردن "Bearer " از توکن
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "فرمت توکن صحیح نیست"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		// پارس و اعتبارسنجی توکن
		token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "توکن نامعتبر است"})
			c.Abort()
			return
		}

		// استخراج ادعاها و قرار دادن در کانتکست برای استفاده در هندلرها
		if claims, ok := token.Claims.(*utils.Claims); ok {
			c.Set("adminID", claims.AdminID)
			c.Set("username", claims.Username)
		}

		c.Next()
	}
}
