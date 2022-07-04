package carts

import (
	"errors"
	"log"

	model "product-services/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	Create(cart model.Cart) (*model.Cart, error)
	CreateCartItem(cartItem model.CartItem) (*model.CartItem, error)
	GetAll() (*[]model.Cart, error)
	GetById(id uint) (*model.Cart, error)
	GetCartItemById(id uint) (*model.CartItem, error)
	GetCartItemId(cart_id, product_id uint) (*model.CartItem, error)
	Update(id uint, cart map[string]interface{}) (*model.Cart, error)
	UpdateCartItem(id uint, updatedData map[string]interface{}) (*model.CartItem, error)
	Delete(id uint) (bool, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}

func (r *repository) Create(cart model.Cart) (*model.Cart, error) {
	var currentCart model.Cart
	if err := r.DB.Debug().Where("customer_id = ?", cart.CustomerID).Find(&currentCart).Error; err != nil {
		log.Println(err)
	}
	log.Println("customer_id ", currentCart.CustomerID, "cart", currentCart)

	if currentCart.CustomerID == cart.CustomerID {
		return &currentCart, nil
	}

	if err := r.DB.Save(&cart).Error; err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *repository) CreateCartItem(cartItem model.CartItem) (*model.CartItem, error) {
	var currentCartItem *model.CartItem
	var err error
	currentCartItem, err = r.GetCartItemId(cartItem.CartID, cartItem.ProductID)

	if err != nil {
		log.Println(err)
	}
	log.Println("cart_id ", currentCartItem.CartID, "product_id", currentCartItem.ProductID, "cart items", currentCartItem)

	var qty uint

	if currentCartItem.CartID == cartItem.CartID {
		var updatedData = make(map[string]interface{})
		updatedData["cart_id"] = currentCartItem.CartID
		updatedData["product_id"] = currentCartItem.ProductID
		qty = currentCartItem.Qty + cartItem.Qty
		updatedData["qty"] = qty

		// update qty
		updates, err := r.UpdateCartItem(currentCartItem.ID, updatedData)
		log.Println("Create Cart Items", updates, err)
		if err != nil {
			return nil, err
		}
		return updates, nil
	}

	if err := r.DB.Save(&cartItem).Error; err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *repository) GetAll() (*[]model.Cart, error) {
	var categories []model.Cart
	if err := r.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return &categories, nil
}

func (r *repository) GetById(id uint) (*model.Cart, error) {
	var cart model.Cart
	if err := r.DB.Where("id = ?", id).Find(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *repository) GetCartItemById(id uint) (*model.CartItem, error) {
	var cartItem model.CartItem
	if err := r.DB.Where("id = ?", id).Find(&cartItem).Error; err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *repository) GetCartItemId(cart_id, product_id uint) (*model.CartItem, error) {
	var currentCartItem *model.CartItem

	if err := r.DB.Debug().Where("cart_id = ? AND product_id = ?", cart_id, product_id).Find(&currentCartItem).Error; err != nil {
		return nil, err
	}
	return currentCartItem, nil
}

func (r *repository) Update(id uint, updatedData map[string]interface{}) (*model.Cart, error) {
	if cart, _ := r.GetById(id); cart.ID <= 0 {
		return nil, errors.New("data not found")
	}
	var newCart model.Cart
	if err := r.DB.Model(model.Cart{}).Where("id = ?", id).Updates(updatedData).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Where("id = ?", id).Find(&newCart).Error; err != nil {
		return nil, err
	}

	return &newCart, nil
}

func (r *repository) UpdateCartItem(id uint, updatedData map[string]interface{}) (*model.CartItem, error) {
	log.Println("udpate cart items id, updatedData", id, updatedData)
	if cartItems, _ := r.GetCartItemById(id); cartItems.ID <= 0 {
		return nil, errors.New("data not found")
	}
	var newCartItem model.CartItem
	if err := r.DB.Model(model.CartItem{}).Where("id = ?", id).Updates(updatedData).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Where("id = ?", id).Find(&newCartItem).Error; err != nil {
		return nil, err
	}

	return &newCartItem, nil
}

func (r *repository) Delete(id uint) (bool, error) {
	if cart, _ := r.GetById(id); cart.ID <= 0 {
		return false, errors.New("data not found")
	}
	if err := r.DB.Debug().Where("id = ?", id).Delete(&model.Cart{}).Error; err != nil {
		return false, errors.New("failed to delete data")
	}

	return true, nil
}
