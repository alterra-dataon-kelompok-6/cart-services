package dto

import model "product-services/internal/models"

type CartRequestParams struct {
	ID uint `json:"id" param:"id" query:"id" form:"id" xml:"id"`
}

type CartItemDetailRequestParams struct {
	ID         uint `json:"id" param:"id" query:"id" form:"id" xml:"id"`
	CartItemID uint `json:"cart_item_id" param:"cart_item_id" query:"cart_item_id" form:"cart_item_id" xml:"cart_item_id"`
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

type Product struct {
	ProductID   uint   `json:"product_id" validate:"required"`
	ProductName string `json:"product_name" validate:"required"`
}

type CartResponseGetById struct {
	model.Cart
	CartItems []model.CartItem `json:"cart_items"`
	// Customer  Customer         `json:"customer"` // next development
}

type CartResponseCartItemDetails struct {
	model.Cart
	CartItem model.CartItem `json:"cart_item"`
	Product  Product        `json:"product"`
	// Customer  Customer       `json:"customer"` // next development
}
