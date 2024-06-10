package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthKey string

const TokenKey AuthKey = "Token"
const UserClaimKey AuthKey = "UserClaim"

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

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("cannot generate hash from string")
	}

	return string(hash), nil
}

func ComparePassword(password string, hashedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hashedPassword)); err != nil {
		log.Println(err)
		return false
	}

	return true
}
