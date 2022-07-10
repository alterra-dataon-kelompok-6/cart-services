package carts

import (
	"product-services/middleware"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("/all", h.GetAll)                                     // admin
	g.GET("/:id/cart_item/:cart_item_id", h.GetCartItemDetails) // admin
	g.GET("/:id", h.GetById)                                    // admin
	g.PUT("/:id", h.Update /*, middleware.ValidateToken*/)      // admin
	g.DELETE("/:id", h.Delete /*, middleware.ValidateToken*/)   // admin

	g.GET("", h.GetCustomerCart, middleware.ValidateToken)                                    // customer
	g.GET("/cart_item/:cart_item_id", h.GetCustomerCartItemDetails, middleware.ValidateToken) // customer
	g.POST("", h.Create, middleware.ValidateToken)                                            // customer
	g.PUT("", h.UpdateCustomerCart, middleware.ValidateToken)                                 // customer
	g.DELETE("", h.DeleteCustomerCart, middleware.ValidateToken)                              // customer
}
