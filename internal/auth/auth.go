package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/carlqt/ezsplit/graph/model"
	jwt "github.com/golang-jwt/jwt/v5"
)

const TokenKey = "Token"

type UserClaim struct {
	jwt.RegisteredClaims
	ID       string `json:"id"`
	Username string `json:"username"`
}

// getBearerToken extracts the bearer token from the Authorization header.
// An error is returned if Authorization header is missing or the format is invalid.
func getBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header format must be Bearer <token>")
	}

	return parts[1], nil
}

// TODO: Move all auth/jwt related functions to a separate package
func ValidateBearerToken(bearerToken string, secret []byte) (model.User, error) {
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		log.Println(err)
		return model.User{}, err
	}

	if claims, ok := token.Claims.(UserClaim); ok && token.Valid {
		// TODO: Research if this is enough or should we validate the user in the database
		return model.User{
			ID:       claims.ID,
			Username: claims.Username,
		}, nil
	}

	errMsg := fmt.Sprint("invalid token: ", bearerToken)
	return model.User{}, errors.New(errMsg)
}

// BearerTokenMiddleware extracts the bearer token from the Authorization header
// and stores it in the context.
// The error is ignored because the token is optional and the resolver will handle the error.
func BearerTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken, _ := getBearerToken(r)

		ctx := context.WithValue(r.Context(), TokenKey, bearerToken)
		newReq := r.WithContext(ctx)

		next.ServeHTTP(w, newReq)
	})
}
