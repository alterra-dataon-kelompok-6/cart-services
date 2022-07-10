package carts

import (
	"errors"
	"fmt"
	"log"

	"product-services/internal/dto"
	"product-services/internal/factory"
	model "product-services/internal/models"
	"product-services/internal/repository"
	"product-services/libs/api"
)

type Service interface {
	// admin roles
	GetAll() (*[]model.Cart, error)
	GetById(payload dto.CartRequestParams) (*dto.CartResponseGetById, error)
	GetCartItemDetails(payload dto.CartItemDetailRequestParams) (*dto.CartResponseCartItemDetails, error)
	Update(id uint, payload dto.CartRequestBodyUpdate) (*model.CartItem, error)
	Delete(payload dto.CartRequestParams) (interface{}, error)
	// customer_roles and admin_roles
	Create(payload dto.CartRequestBodyCreate) (*model.Cart, *model.CartItem, error)
	GetCustomerCart(customer_id uint) (*dto.CartResponseGetById, error)
	GetCustomerCartItemDetails(CustomerID uint, payload dto.CartItemDetailRequestParams) (*dto.CartResponseCartItemDetails, error)
	UpdateCustomerCart(CustomerID uint, payload dto.CartRequestBodyUpdate) (*model.CartItem, error)
	DeleteCustomerCart(CustomerID uint) (interface{}, error)
}

type service struct {
	CartRepository repository.CartRepository
}

func NewService(f *factory.Factory) Service {
	return &service{
		CartRepository: repository.NewCartRepo(f.DB),
	}
}

var CartItemRepo = repository.NewCartItemRepo(factory.NewFactory().DB)

func (s service) Create(payload dto.CartRequestBodyCreate) (*model.Cart, *model.CartItem, error) {
	var newCart = model.Cart{
		CustomerID: payload.CustomerID,
	}
	// check is Customer ID not 0
	if newCart.CustomerID == 0 {
		return nil, nil, errors.New("customer")
	}
	// check cart is already exist ?
	currentCart, err := s.CartRepository.GetCart(0, payload.CustomerID)
	if err != nil {
		return nil, nil, err
	}

	// log.Println(currentCart, "debug ketemu 1")

	var cart *model.Cart

	// if current cart not found, create new cart
	if currentCart.CustomerID == 0 {
		cart, err = s.CartRepository.Create(newCart)

		if err != nil {
			return nil, nil, err
		}
		// log.Println("debug insert doble", cart)

	} else {
		cart = currentCart
	}

	// check cart item is already exist ?
	currentCartItem, err := CartItemRepo.GetCartItem(0, cart.ID, payload.ProductID)
	if err != nil {
		return nil, nil, err
	}

	// check product stock
	product := api.GetProduct(payload.ProductID)
	stock := product.Data.Product.Stock

	log.Println("api get product", product)

	var newCartItem = model.CartItem{
		CartID:    cart.ID, // diambil dari data diatas
		ProductID: payload.ProductID,
		Qty:       payload.Qty,
	}

	// if product id == 0 / product not found return error product not found
	if product.Data.Product.ID == 0 {
		return nil, nil, fmt.Errorf("product with id %d not found", payload.ProductID)
	}

	var cartItem *model.CartItem
	if currentCartItem.ID == 0 {
		// check stock
		if stock < newCartItem.Qty {
			return nil, nil, fmt.Errorf("quantity exceeds available stock %d", stock)
		}
		cartItem, err = CartItemRepo.Create(newCartItem)
	} else {
		var updatedData map[string]interface{} = make(map[string]interface{})
		qty := currentCartItem.Qty + payload.Qty
		log.Println(qty, currentCartItem.Qty, payload.Qty, "qty", currentCartItem)
		updatedData["qty"] = qty

		// check stock
		if stock < qty {
			return nil, nil, fmt.Errorf("quantity exceeds available stock %d", stock)
		}

		cartItem, err = CartItemRepo.Update(currentCartItem.ID, updatedData)
		if err != nil {
			return nil, nil, err
		}
	}

	if err != nil {
		return nil, nil, err
	}
	return cart, cartItem, nil
}

func (s service) GetAll() (*[]model.Cart, error) {
	carts, err := s.CartRepository.GetAll()
	if err != nil || len(*carts) <= 0 {
		return nil, errors.New("data is empty")
	}
	return carts, nil
}

func (s service) GetById(payload dto.CartRequestParams) (*dto.CartResponseGetById, error) {
	var result = new(dto.CartResponseGetById)

	cart, err := s.CartRepository.GetCart(payload.ID, 0)
	if err != nil {
		return nil, err
	}

	result.Cart = *cart

	cart_items, err := CartItemRepo.GetAll(cart.ID)
	if err != nil {
		return nil, err
	}

	result.CartItems = *cart_items

	log.Println("service", cart, err)
	if err != nil || cart.ID == 0 {
		return nil, errors.New("data not found")
	}
	return result, nil
}

func (s service) GetCustomerCart(customer_id uint) (*dto.CartResponseGetById, error) {
	var result = new(dto.CartResponseGetById)

	// get customer cart
	cart, err := s.CartRepository.GetCart(0, customer_id)
	if err != nil {
		return nil, err
	}

	result.Cart = *cart

	// get cart item
	cart_items, err := CartItemRepo.GetAll(cart.ID)
	if err != nil {
		return nil, err
	}

	result.CartItems = *cart_items

	log.Println("service", cart, err)
	if err != nil || cart.ID == 0 {
		return nil, errors.New("data not found")
	}
	return result, nil
}

func (s service) GetCartItemDetails(payload dto.CartItemDetailRequestParams) (*dto.CartResponseCartItemDetails, error) {
	log.Println("params", payload)

	var result dto.CartResponseCartItemDetails

	// // payload cart
	// var payloadCart dto.CartRequestParams
	// payloadCart.ID = payload.ID

	// get data cart
	cart, err := s.CartRepository.GetCart(payload.ID, 0)
	if err != nil {
		return nil, err
	}

	result.Cart = *cart

	// get cart item
	cart_item, err := CartItemRepo.GetCartItem(payload.CartItemID, 0, 0)
	if err != nil {
		return nil, err
	}

	if cart.ID != cart_item.CartID {
		return nil, fmt.Errorf("cart item id %d with cart id %d not found", payload.CartItemID, cart.ID)
	}

	result.CartItem = *cart_item

	// get product details
	product := api.GetProduct(cart_item.ProductID)

	// if product id == 0 / product not found return error product not found
	if product.Data.Product.ID == 0 {
		return nil, fmt.Errorf("product with id %d not found", cart_item.ProductID)
	}

	result.Product.ProductID = product.Data.ID
	result.Product.ProductName = product.Data.Name

	return &result, nil
}

func (s service) GetCustomerCartItemDetails(CustomerID uint, payload dto.CartItemDetailRequestParams) (*dto.CartResponseCartItemDetails, error) {
	log.Println("params", payload)

	var result dto.CartResponseCartItemDetails

	// find cart id with customer id
	cart, err := s.CartRepository.GetCart(0, CustomerID)
	// get data cart
	// cart, err := s.CartRepository.GetCart(payload.ID, 0)
	if err != nil {
		return nil, err
	}

	result.Cart = *cart

	// get cart item
	cart_item, err := CartItemRepo.GetCartItem(payload.CartItemID, 0, 0)
	if err != nil {
		return nil, err
	}

	if cart.ID != cart_item.CartID {
		return nil, fmt.Errorf("cart item id %d with cart id %d not found", payload.CartItemID, cart.ID)
	}

	result.CartItem = *cart_item

	// get product details
	product := api.GetProduct(cart_item.ProductID)

	// if product id == 0 / product not found return error product not found
	if product.Data.Product.ID == 0 {
		return nil, fmt.Errorf("product with id %d not found", cart_item.ProductID)
	}

	result.Product.ProductID = product.Data.ID
	result.Product.ProductName = product.Data.Name

	return &result, nil
}

func (s service) Update(id uint, payload dto.CartRequestBodyUpdate) (*model.CartItem, error) {
	// get cartItemsId
	cartId := id
	cartItemId, err := CartItemRepo.GetCartItem(0, cartId, payload.ProductID)
	if err != nil {
		return nil, err
	}

	updatedCartItem := make(map[string]interface{})
	updatedCartItem["cart_id"] = cartId
	updatedCartItem["product_id"] = payload.ProductID
	updatedCartItem["qty"] = payload.Qty

	// check product stock
	product := api.GetProduct(payload.ProductID)
	stock := product.Data.Product.Stock

	log.Println("api get product", product)

	// if product id == 0 / product not found return error product not found
	if product.Data.Product.ID == 0 {
		return nil, fmt.Errorf("product with id %d not found", payload.ProductID)
	}

	if payload.Qty > stock {
		return nil, fmt.Errorf("quantity exceeds available stock %d", stock)
	}

	cartItem, err := CartItemRepo.Update(cartItemId.ID, updatedCartItem)
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (s service) UpdateCustomerCart(CustomerID uint, payload dto.CartRequestBodyUpdate) (*model.CartItem, error) {
	// get cart_id
	cart, err := s.CartRepository.GetCart(0, CustomerID)
	if err != nil {
		return nil, err
	}
	// get cart_item_id
	cartItemId, err := CartItemRepo.GetCartItem(0, cart.ID, payload.ProductID)
	if err != nil {
		return nil, err
	}

	updatedCartItem := make(map[string]interface{})
	updatedCartItem["cart_id"] = cart.ID
	updatedCartItem["product_id"] = payload.ProductID
	updatedCartItem["qty"] = payload.Qty

	// check product stock
	product := api.GetProduct(payload.ProductID)
	stock := product.Data.Product.Stock

	log.Println("api get product", product)

	// if product id == 0 / product not found return error product not found
	if product.Data.Product.ID == 0 {
		return nil, fmt.Errorf("product with id %d not found", payload.ProductID)
	}

	if payload.Qty > stock {
		return nil, fmt.Errorf("quantity exceeds available stock %d", stock)
	}

	cartItem, err := CartItemRepo.Update(cartItemId.ID, updatedCartItem)
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (s service) Delete(payload dto.CartRequestParams) (interface{}, error) {
	deleted, err := s.CartRepository.Delete(payload.ID)
	if err != nil {
		return nil, err
	}
	log.Println(deleted, "deleted")
	return deleted, err
}

func (s service) DeleteCustomerCart(CustomerID uint) (interface{}, error) {
	// get cart_id
	cart, err := s.CartRepository.GetCart(0, CustomerID)
	if err != nil {
		return nil, err
	}
	log.Println("cart id", cart.ID, cart)
	deleted, err := s.CartRepository.Delete(cart.ID)
	if err != nil {
		return nil, err
	}
	// log.Println(deleted, "deleted")
	return deleted, err
}
