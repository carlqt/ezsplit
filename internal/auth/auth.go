package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const TokenKey = "Token"

type UserClaim struct {
	jwt.RegisteredClaims
	ID       string `json:"id"`
	Username string `json:"username"`
}

func CreateAndSignToken(userClaim UserClaim, secret []byte) (string, error) {
	userClaim.IssuedAt = jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return signedToken, nil
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

func ValidateBearerToken(bearerToken string, secret []byte) (UserClaim, error) {
	token, err := jwt.ParseWithClaims(bearerToken, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		log.Println(err)
		return UserClaim{}, err
	} else if claims, ok := token.Claims.(*UserClaim); ok && token.Valid {
		return *claims, nil
	} else {
		log.Println("unknown claims type")
		return *claims, errors.New("unknown claims type")
	}
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
