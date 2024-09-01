package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("my_secret_key")

type JWTClaim struct {
	Username string
	jwt.RegisteredClaims
}
