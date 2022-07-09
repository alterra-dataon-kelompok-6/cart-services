package dto

import model "product-services/internal/models"

type CartRequestParams struct {
	ID uint `json:"id" param:"id" query:"id" form:"id" xml:"id"`
}

type CartRequestBodyCreate struct {
	CustomerID uint `json:"customer_id" validate:"required"`
	ProductID  uint `json:"product_id" validate:"required"`
	Qty        uint `json:"qty" validate:"required"`
}

type CartRequestBodyUpdate struct {
	ProductID uint `json:"product_id" validate:"required"`
	Qty       uint `json:"qty" validate:"required"`
}

type Customer struct {
	CustomerID   uint   `json:"customer_id" validate:"required"`
	CustomerName string `json:"customer_name" validate:"required"`
}

type CartResponseGetById struct {
	model.Cart
	CartItems []model.CartItem `json:"cart_items"`
	// Customer  Customer         `json:"customer"` // next development
}
