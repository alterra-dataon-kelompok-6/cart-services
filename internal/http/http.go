package http

import (
	"product-services/internal/apps/carts"
	"product-services/internal/factory"

	"github.com/labstack/echo/v4"
)

func NewHttp(e *echo.Echo, f *factory.Factory) {
	carts.NewHandler(f).Route(e.Group("/carts"))
}
