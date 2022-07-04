package model

import (
	"time"

	"gorm.io/gorm"
)

type Common struct {
	ID        uint            `json:"id" gorm:"primarykey;autoIncrement"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
