package cart_items

import (
	"errors"
	"log"

	model "product-services/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	Create(cartItem model.CartItem) (*model.CartItem, error)
	GetAll(cart_id uint) (*[]model.CartItem, error)
	GetCartItem(id, cart_id, product_id uint) (*model.CartItem, error)
	Update(id uint, updatedData map[string]interface{}) (*model.CartItem, error)
	Delete(id uint) (bool, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepo(DB *gorm.DB) Repository {
	return &repository{DB: DB}
}

func (r *repository) Create(cartItem model.CartItem) (*model.CartItem, error) {
	if err := r.DB.Save(&cartItem).Error; err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *repository) GetAll(cart_id uint) (*[]model.CartItem, error) {
	var cart_items []model.CartItem
	if cart_id != 0 {
		if err := r.DB.Debug().Where("cart_id = ?", cart_id).Find(&cart_items).Error; err != nil {
			return nil, err
		}
		return &cart_items, nil
	}
	log.Println(cart_items, "cart_items get all")

	if err := r.DB.Find(&cart_items).Error; err != nil {
		return nil, err
	}
	return &cart_items, nil
}

func (r *repository) GetCartItem(id, cart_id, product_id uint) (*model.CartItem, error) {
	var cartItem model.CartItem
	if cart_id != 0 && product_id != 0 {
		if err := r.DB.Where("cart_id = ? AND product_id = ?", cart_id, product_id).Find(&cartItem).Error; err != nil {
			return nil, err
		}
		return &cartItem, nil
	}
	if err := r.DB.Where("id = ?", id).Find(&cartItem).Error; err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *repository) Update(id uint, updatedData map[string]interface{}) (*model.CartItem, error) {
	log.Println("udpate cart items id, updatedData", id, updatedData)
	if cartItems, _ := r.GetCartItem(id, 0, 0); cartItems.ID <= 0 {
		return nil, errors.New("data not found")
	}
	var newCartItem model.CartItem
	if err := r.DB.Debug().Model(model.CartItem{}).Where("id = ?", id).Updates(updatedData).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Where("id = ?", id).Find(&newCartItem).Error; err != nil {
		return nil, err
	}

	return &newCartItem, nil
}

func (r *repository) Delete(id uint) (bool, error) {
	if cartItem, _ := r.GetCartItem(id, 0, 0); cartItem.ID <= 0 {
		return false, errors.New("data not found")
	}
	if err := r.DB.Debug().Where("id = ?", id).Delete(&model.CartItem{}).Error; err != nil {
		return false, errors.New("failed to delete data")
	}

	return true, nil
}
