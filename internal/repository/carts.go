package repository

import (
	"errors"

	model "product-services/internal/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	Create(cart model.Cart) (*model.Cart, error)
	GetAll() (*[]model.Cart, error)
	GetCart(id, customer_id uint) (*model.Cart, error)
	Update(id uint, cart map[string]interface{}) (*model.Cart, error)
	Delete(id uint) (bool, error)
}

type cartRepository struct {
	DB *gorm.DB
}

func NewCartRepo(DB *gorm.DB) CartRepository {
	return &cartRepository{DB: DB}
}

func (r *cartRepository) Create(cart model.Cart) (*model.Cart, error) {
	if err := r.DB.Save(&cart).Error; err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) GetAll() (*[]model.Cart, error) {
	var carts []model.Cart
	if err := r.DB.Find(&carts).Error; err != nil {
		return nil, err
	}
	return &carts, nil
}

func (r *cartRepository) GetCart(id, customer_id uint) (*model.Cart, error) {
	var cart model.Cart

	if customer_id != 0 {
		if err := r.DB.Debug().Where("customer_id = ?", customer_id).Find(&cart).Error; err != nil {
			return nil, err
		}
		return &cart, nil
	}

	if err := r.DB.Where("id = ?", id).Find(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) Update(id uint, updatedData map[string]interface{}) (*model.Cart, error) {
	if cart, _ := r.GetCart(id, 0); cart.ID <= 0 {
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

func (r *cartRepository) Delete(id uint) (bool, error) {
	if cart, _ := r.GetCart(id, 0); cart.ID <= 0 {
		return false, errors.New("data not found")
	}
	if err := r.DB.Debug().Where("id = ?", id).Delete(&model.Cart{}).Error; err != nil {
		return false, errors.New("failed to delete data")
	}

	return true, nil
}
