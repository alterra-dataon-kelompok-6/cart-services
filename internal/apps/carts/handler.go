package carts

import (
	"log"
	"net/http"

	"product-services/internal/dto"
	"product-services/internal/factory"
	"product-services/middleware"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service Service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

// get all cart
func (h handler) GetAll(e echo.Context) error {
	categories, err := h.service.GetAll()

	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "data not found",
		})
	}
	return e.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"data":   categories,
	})

}

// get cart by customer id
func (h handler) GetCustomerCart(e echo.Context) error {
	customer_id := middleware.GetUserIdFromToken(e)

	categories, err := h.service.GetCustomerCart(customer_id)

	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "data not found",
		})
	}
	return e.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"data":   categories,
	})

}

// get cart by cart_id
func (h handler) GetById(e echo.Context) error {
	// id, _ := strconv.Atoi(e.Param("id"))
	var payload dto.CartRequestParams

	if err := (&echo.DefaultBinder{}).BindPathParams(e, &payload); err != nil {
		log.Println(err)
	}

	cart, err := h.service.GetById(payload)
	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "data not found",
		})
	}
	return e.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"data":   cart,
	})
}

// create new cart
func (h handler) Create(e echo.Context) error {
	userId := middleware.GetUserIdFromToken(e)
	log.Println(userId, "userId")
	var payload dto.CartRequestBodyCreate
	if err := e.Bind(&payload); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "invalid data",
		})
	}
	cart, cartItem, err := h.service.Create(payload)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "failed to create data",
		})
	}
	return e.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"data": map[string]interface{}{
			"cart":       cart,
			"cart_items": cartItem,
		},
	})
}

// update cart data
func (h handler) Update(e echo.Context) error {
	// id, _ := strconv.Atoi(e.Param("id"))
	var id dto.CartRequestParams
	var payload dto.CartRequestBodyUpdate

	if err := (&echo.DefaultBinder{}).BindPathParams(e, &id); err != nil {
		log.Println(err)
	}

	if err := e.Bind(&payload); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "invalid data",
		})
	}

	cart, err := h.service.Update(id.ID, payload)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "failed to update data",
		})
	}
	return e.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"data":   cart,
	})

}

// delete cart data
func (h handler) Delete(e echo.Context) error {
	// id, _ := strconv.Atoi(e.Param("id"))
	var payload dto.CartRequestParams

	if err := (&echo.DefaultBinder{}).BindPathParams(e, &payload); err != nil {
		log.Println(err)
	}
	_, err := h.service.Delete(payload)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "failed to delete data",
		})
	}
	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "data has been deleted",
	})
}
