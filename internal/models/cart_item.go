package model

type CartItem struct {
	Common
	CartID    uint `json:"cart_id" gorm:"index:idx_cart,unique"`
	ProductID uint `json:"product_id" gorm:"index:idx_cart,unique"`
	Qty       uint `json:"qty"`
}
