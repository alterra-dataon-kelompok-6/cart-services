package carts

import (
	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("", h.GetAll)
	g.GET("/:id", h.GetById)
	g.POST("", h.Create /*, middleware.ValidateToken*/)
	g.PUT("/:id", h.Update /*, middleware.ValidateToken*/)
	g.DELETE("/:id", h.Delete /*, middleware.ValidateToken*/)
}
