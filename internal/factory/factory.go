package factory

import (
	"product-services/database"

	"gorm.io/gorm"
)

type Factory struct {
	DB *gorm.DB
}

func NewFactory() *Factory {
	return &Factory{
		DB: database.GetConnection(),
	}
}
