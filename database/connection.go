package database

import (
	"fmt"
	"log"

	model "product-services/internal/models"
	"product-services/libs/env"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

type Tables []interface{}

var tables Tables = Tables{&model.Cart{}, &model.CartItem{}}

func init() {
	dbConn = ConnectDB()
	MigrateDB(dbConn)
}

func connectMysql() *gorm.DB {
	config := dbConfig{
		DB_Username: env.GetEnv("DB_USERNAME"),
		DB_Password: env.GetEnv("DB_PASSWORD"),
		DB_Port:     env.GetEnv("DB_PORT"),
		DB_Host:     env.GetEnv("DB_HOST"),
		DB_Name:     env.GetEnv("DB_NAME"),
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", config.DB_Username, config.DB_Password, config.DB_Host, config.DB_Port, config.DB_Name)
	var err error
	dbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil
	}
	return dbConn
}

func ConnectDB() *gorm.DB {
	// use mysql
	return connectMysql()
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(tables...)
}

func GetConnection() *gorm.DB {
	if dbConn == nil {
		ConnectDB()
	}
	log.Println("Connected to database", dbConn)
	return dbConn
}
