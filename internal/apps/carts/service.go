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
	Update(id uint, payload dto.CartRequestBodyUpdate) (*model.CartItem, error)
	Delete(payload dto.CartRequestParams) (interface{}, error)
	// customer_roles and admin_roles
	Create(payload dto.CartRequestBodyCreate) (*model.Cart, *model.CartItem, error)
	GetCustomerCart(customer_id uint) (*dto.CartResponseGetById, error)
	UpdateCustomerCart(CustomerID uint, payload dto.CartRequestBodyUpdate) (*model.CartItem, error)
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

	log.Println(currentCart, "debug ketemu 1")

	var cart *model.Cart

	if currentCart.CustomerID == 0 {
		cart, err = s.CartRepository.Create(newCart)

		if err != nil {
			return nil, nil, err
		}
		log.Println("debug insert doble", cart)

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
	categories, err := s.CartRepository.GetAll()
	if err != nil || len(*categories) <= 0 {
		return nil, errors.New("data is empty")
	}
	return categories, nil
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

	cart, err := s.CartRepository.GetCart(0, customer_id)
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
