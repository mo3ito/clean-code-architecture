package models

import (
	"time"
	"gorm.io/gorm"
)

// Business مدل کسب‌وکار (رستوران / کافه)
type Business struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	Name                  string         `gorm:"type:varchar(150);not null" json:"name"`
	Slug                  string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"` // مثل: cafe-ahoura
	Phone                 string         `gorm:"type:varchar(20)" json:"phone"`
	Address               string         `gorm:"type:text" json:"address"`
	LogoURL               string         `gorm:"type:varchar(255)" json:"logo_url,omitempty"`
	Theme                 string         `gorm:"type:json" json:"theme,omitempty"` // برای ذخیره رنگ‌بندی سفارشی منو
	Plan                  string         `gorm:"type:varchar(50);default:'trial'" json:"plan"` // trial, bronze, silver, gold
	SubscriptionExpiresAt *time.Time     `json:"subscription_expires_at"` // تاریخ انقضای پکیج
	IsActive              bool           `gorm:"default:true" json:"is_active"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"` // حذف نرم (Soft Delete)
}

// متد کمکی برای بررسی وضعیت فعال بودن اشتراک
func (b *Business) HasActiveSubscription() bool {
	if !b.IsActive {
		return false
	}
	if b.SubscriptionExpiresAt == nil {
		return false
	}
	return b.SubscriptionExpiresAt.After(time.Now())
}

// BusinessAdmin مدل ادمین / صاحب یا پرسنل کسب‌وکار
type BusinessAdmin struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	BusinessID   uint           `gorm:"not null;index" json:"business_id"`
	Business     Business       `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE" json:"business,omitempty"`
	Username     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	PasswordHash string         `gorm:"not null" json:"-"`
	FullName     string         `gorm:"type:varchar(100)" json:"full_name"`
	Phone        string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	Role         string         `gorm:"type:varchar(20);default:'owner'" json:"role"` // owner, manager, waiter
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
