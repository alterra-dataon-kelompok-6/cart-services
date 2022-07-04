package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func CreateToken(name string) (string, error) {
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
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return t, nil
}

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authToken := c.Request().Header.Get("Authorization")
		log.Println("01 - authToken", authToken)
		if authToken == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"status":  false,
				"message": "unauthorized",
			})
		}

		tokenString := strings.Split(authToken, " ")[1]
		log.Println("02 - tokenString", tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("error token")
			}
			return []byte("secret"), nil
		})
		log.Println("03 - token", token)
		if !token.Valid || err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"status":  false,
				"message": "unauthorized",
			})
		}
		return next(c)
	}
}
