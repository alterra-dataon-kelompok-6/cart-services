package carts

import (
	"errors"
	"log"

	"product-services/internal/dto"
	"product-services/internal/factory"
	model "product-services/internal/models"
)

type Service interface {
	Create(payload dto.CartRequestBodyCreate) (*model.Cart, *model.CartItem, error)
	GetAll() (*[]model.Cart, error)
	GetById(payload dto.CartRequestParams) (*model.Cart, error)
	Update(id uint, payload dto.CartRequestBodyUpdate) (*model.CartItem, error)
	Delete(payload dto.CartRequestParams) (interface{}, error)
}

type service struct {
	CartRepository Repository
}

func NewService(f *factory.Factory) Service {
	return &service{
		CartRepository: NewRepo(f.DB),
	}
}

func (s service) Create(payload dto.CartRequestBodyCreate) (*model.Cart, *model.CartItem, error) {
	var newCart = model.Cart{
		CustomerID: payload.CustomerID,
	}

	cart, err := s.CartRepository.Create(newCart)
	if err != nil {
		return nil, nil, err
	}
	var newCartItem = model.CartItem{
		CartID:    cart.ID, // diambil dari data diatas
		ProductID: payload.ProductID,
		Qty:       payload.Qty,
	}
	cartItem, err := s.CartRepository.CreateCartItem(newCartItem)

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

func (s service) GetById(payload dto.CartRequestParams) (*model.Cart, error) {
	cart, err := s.CartRepository.GetById(payload.ID)
	log.Println("service", cart, err)
	if err != nil || cart.ID == 0 {
		return nil, errors.New("data not found")
	}
	return cart, nil
}

func (s service) Update(id uint, payload dto.CartRequestBodyUpdate) (*model.CartItem, error) {
	// get cartItemsId
	cartId := id
	cartItemId, err := s.CartRepository.GetCartItemId(cartId, payload.ProductID)
	if err != nil {
		return nil, err
	}

	updatedCartItem := make(map[string]interface{})
	updatedCartItem["cart_id"] = cartId
	updatedCartItem["product_id"] = payload.ProductID
	updatedCartItem["qty"] = payload.Qty

	cartItem, err := s.CartRepository.UpdateCartItem(cartItemId.ID, updatedCartItem)
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
