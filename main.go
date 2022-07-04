package main

import (
	"product-services/internal/factory"
	"product-services/internal/http"
	m "product-services/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	f := factory.NewFactory()
	http.NewHttp(app, f)

	m.LogMiddleware(app)

	app.Start(":8088")
}
