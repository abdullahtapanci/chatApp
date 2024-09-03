package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const SECRET_KEY = "MEhWzCzIVB149XLDxgIpjMuZjI6hnE0VCd6oUib7fEM"

// Define a struct to represent the claims (payload) of the JWT token
type Claims struct {
	Id       string `json:"id"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

// Generate JWT token
func generateJWTToken(id string, userName string) (string, error) {
	// Create the claims
	claims := &Claims{
		Id:       id,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Sign the token with a secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// Validate the JWT token from the Authorization header
func ValidateToken(c *fiber.Ctx) (*jwt.Token, jwt.MapClaims, error) {
	// Get the token from the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		fmt.Println("No token provided")
		return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "No token provided")
	}

	// Parse the token
	tokenString := authHeader[len("Bearer "):] // Remove "Bearer " prefix
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		fmt.Println("Invalid token")
		return nil, nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	// Access claims and return them along with the parsed token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Failed to extract claims")
		return nil, nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to extract claims")
	}

	return token, claims, nil
}
