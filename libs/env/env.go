package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	LoadEnv()
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// testing
	log.Println(os.Getenv("DB_USERNAME"))
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
