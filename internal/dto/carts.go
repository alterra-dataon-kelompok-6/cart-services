package dto

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
