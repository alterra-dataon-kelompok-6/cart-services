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
	carts, err := h.service.GetAll()

	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "data not found",
		})
	}
	return e.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"data":   carts,
	})

}

// get cart by customer id
func (h handler) GetCustomerCart(e echo.Context) error {
	CustomerID := middleware.GetUserIdFromToken(e)
	log.Println(CustomerID, "get customer cart")
	carts, err := h.service.GetCustomerCart(CustomerID)

	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "data not found",
		})
	}
	return e.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"data":   carts,
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

func (h handler) GetCartItemDetails(e echo.Context) error {
	// bind params to payload
	var payload dto.CartItemDetailRequestParams
	if err := (&echo.DefaultBinder{}).BindPathParams(e, &payload); err != nil {
		log.Println(err)
	}

	result, err := h.service.GetCartItemDetails(payload)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": err.Error(),
		})
	}

	// return errors.New("error")
	return e.JSON(http.StatusOK, map[string]interface{}{
		"payload": result,
	})
}

func (h handler) GetCustomerCartItemDetails(e echo.Context) error {
	// get customer id form jwt token
	CustomerID := middleware.GetUserIdFromToken(e)

	// bind params to payload
	var payload dto.CartItemDetailRequestParams
	if err := (&echo.DefaultBinder{}).BindPathParams(e, &payload); err != nil {
		log.Println(err)
	}

	result, err := h.service.GetCustomerCartItemDetails(CustomerID, payload)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": err.Error(),
		})
	}

	// return errors.New("error")
	return e.JSON(http.StatusOK, map[string]interface{}{
		"payload": result,
	})
}

// create new cart
func (h handler) Create(e echo.Context) error {
	CustomerID := middleware.GetUserIdFromToken(e)
	log.Println(CustomerID, "CustomerID")
	var payload dto.CartRequestBodyCreate
	if err := e.Bind(&payload); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "invalid data",
		})
	}
	payload.CustomerID = CustomerID
	cart, cartItem, err := h.service.Create(payload)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": err.Error(),
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
			"message": err.Error(),
		})
	}
	return e.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"data":   cart,
	})

}

// update customer cart data
func (h handler) UpdateCustomerCart(e echo.Context) error {
	// id, _ := strconv.Atoi(e.Param("id"))
	CustomerID := middleware.GetUserIdFromToken(e)
	log.Println(CustomerID, "CustomerID")
	var payload dto.CartRequestBodyUpdate

	if err := e.Bind(&payload); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "invalid data",
		})
	}

	cart, err := h.service.UpdateCustomerCart(CustomerID, payload)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": err.Error(),
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

// delete customer cart data
func (h handler) DeleteCustomerCart(e echo.Context) error {
	// id, _ := strconv.Atoi(e.Param("id"))
	CustomerID := middleware.GetUserIdFromToken(e)
	log.Println(CustomerID, "CustomerID")

	_, err := h.service.DeleteCustomerCart(CustomerID)
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
