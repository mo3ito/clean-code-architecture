package models

import "time"

type Admin struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	// اضافه کردن size:255 مشکل را به طور کامل حل می‌کند
	Username  string    `gorm:"size:255;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
