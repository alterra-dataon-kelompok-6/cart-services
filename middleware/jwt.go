package middleware

import (
	"errors"
	"log"
	"net/http"
	"product-services/libs/env"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
/* type jwtCustomClaims struct {
	Authorized bool `json:"authorized"`
	UserID     uint `json:"userId"`
	jwt.StandardClaims
} */

var secret string = env.GetEnv("JWT_SECRET")

/* func CreateToken(name string) (string, error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		name,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, nil
} */

// func middleware to validate jwt from request headers
func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authToken := c.Request().Header.Get("Authorization")
		// log.Println("01 - authToken", authToken)
		if authToken == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"status":  false,
				"message": "user unauthorized",
			})
		}

		tokenString := strings.Split(authToken, " ")[1]
		// log.Println("02 - tokenString", tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("error token")
			}
			return []byte(secret), nil
		})
		// log.Println("03 - token", token)
		if !token.Valid || err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"status":  false,
				"message": "user unauthorized",
			})
		}

		return next(c)
	}
}

// func middlewa to get userId from token jwt
func GetUserIdFromToken(e echo.Context) uint {
	authToken := e.Request().Header.Get("Authorization")
	// log.Println("01 - authToken", authToken)
	if authToken == "" {
		return 0
	}

	tokenString := strings.Split(authToken, " ")[1]
	// log.Println("02 - tokenString", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error token")
		}
		return []byte(secret), nil
	})
	// log.Println("03 - token", token)
	if !token.Valid || err != nil {
		log.Println("tidak valid ternyata :)")
		return 0
	}
	userId := token.Claims.(jwt.MapClaims)["userId"].(float64)
	// log.Println("04", userId)
	if userId != 0 {
		return uint(userId)
	}
	return 0
}
