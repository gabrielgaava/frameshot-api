package utils

import (
	"example/web-service-gin/src/core/entity"
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"strings"
	"time"
)

type JwtService struct {
	jwksURL string
	jwks    *keyfunc.JWKS
}

// NewJwtService creates a new instance of JwtService with url for public keys
func NewJwtService(jwksURL string) (*JwtService, error) {

	// Updates each 1 hour
	options := keyfunc.Options{
		RefreshInterval: 1 * time.Hour,
	}

	jwks, err := keyfunc.Get(jwksURL, options)

	if err != nil {
		return nil, fmt.Errorf("error to retrive JWKS data: %w", err)
	}

	return &JwtService{jwksURL: jwksURL, jwks: jwks}, nil
}

// GetUser Validates the token and get user data from JWT claims
func (j *JwtService) GetUser(jwtToken string) (user *entity.User, err error) {

	cleanToken := removeBearerPrefix(jwtToken)
	token, err := jwt.Parse(cleanToken, j.jwks.Keyfunc)

	if err != nil {
		return nil, fmt.Errorf("error on JWT Token validation: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("JWT Token is invalid")
	}

	// Acesse as claims do payload.
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		user = &entity.User{
			Id:            claims["sub"].(string),
			Email:         claims["email"].(string),
			EmailVerified: claims["email_verified"].(bool),
		}

		return user, nil
	}

	log.Fatalf("Failed to parse JWT Token")
	return nil, nil
}

// removeBearerPrefix removes the "Bearer" prefix of token it exists
func removeBearerPrefix(token string) string {
	const prefix = "Bearer "
	if strings.HasPrefix(token, prefix) {
		return strings.TrimPrefix(token, prefix)
	}
	return token
}
