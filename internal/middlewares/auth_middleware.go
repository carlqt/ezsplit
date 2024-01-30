package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal"
	"github.com/golang-jwt/jwt"
)

const userKey = "User"

func getBearerToken(r *http.Request) string {
	authHeaderString := r.Header.Get("Authorization")
	authHeader := strings.Split(authHeaderString, " ")

	return authHeader[1]
}

func validateBearerToken(bearerToken string, secret []byte) (*model.User, error) {
	// Notes: Find a way to return a concrete type instead of an empty interface
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// TODO: Research if this is enough or should we validate the user in the database
		return &model.User{
			ID:       claims["id"].(string),
			Username: claims["username"].(string),
		}, nil
	}

	return nil, err
}

func AuthMiddleware(next http.Handler, conf internal.EnvConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey := []byte(conf.JWTSecret)
		bearerToken := getBearerToken(r)

		user, err := validateBearerToken(bearerToken, secretKey)

		// TODO: Figure out what to do with the error here
		if err != nil {
			log.Println(err)
		}

		// TODO: Fix warning
		ctx := context.WithValue(r.Context(), userKey, user)
		newReq := r.WithContext(ctx)

		log.Println("GET Middleware")

		next.ServeHTTP(w, newReq)
	})
}
