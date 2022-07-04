package model

type Cart struct {
	Common
	CustomerID uint `json:"customer_id" gorm:"uniqueIndex"`
}
