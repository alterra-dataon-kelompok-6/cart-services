package model

type CartItem struct {
	Common
	CartID    uint `json:"cart_id" gorm:"index:idx_cart_item,unique"`
	ProductID uint `json:"product_id" gorm:"index:idx_cart_item,unique"`
	Qty       uint `json:"qty"`
}
