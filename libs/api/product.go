package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"product-services/libs/env"
)

type Category struct {
	Category   string `json:"category"`
	CategoryID uint   `json:"category_id"`
}

type Product struct {
	ID          uint   `json:"id"`
	CategoryID  int    `json:"category_id"`
	Name        string `json:"name"`
	Stock       uint   `json:"stock"`
	Price       uint   `json:"price"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

type Data struct {
	Product
	Category Category `json:"category"`
}

type ProductResponseApi struct {
	Data   Data `json:"data"`
	Status bool `json:"status"`
}

var productBaseUrl = env.GetEnv("URL_PRODUCT_SERVICES")

func GetProduct(id uint) *ProductResponseApi {
	url := fmt.Sprintf("%v/products/%d", productBaseUrl, id)
	log.Println(url)

	// resp, err := http.NewRequest(http.MethodGet, url, nil)
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
	}

	var data ProductResponseApi
	json.Unmarshal(body, &data)

	log.Println("resp", resp, "data", data)
	return &data
}
