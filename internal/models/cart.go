package model

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID         uint            `json:"id" gorm:"primarykey;autoIncrement"`
	CustomerID uint            `json:"customer_id" gorm:"index:idx_cart,unique"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  *gorm.DeletedAt `json:"deleted_at" gorm:"index:idx_cart,unique"`
}
