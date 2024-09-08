package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/carlqt/ezsplit/graph/model"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthKey string

const TokenKey AuthKey = "Token"
const UserClaimKey AuthKey = "UserClaim"

type UserClaim struct {
	jwt.RegisteredClaims
	ID       string          `json:"id"`
	Username string          `json:"username"`
	State    model.UserState `json:"state"`
}

func NewUserClaim(id int64, username string, isVerified bool) UserClaim {
	userID := strconv.Itoa(int(id))
	state := model.UserStateGuest

	if isVerified {
		state = model.UserStateVerified
	}

	return UserClaim{
		ID:       userID,
		Username: username,
		State:    state,
	}
}

func CreateAndSignToken(userClaim UserClaim, secret []byte) (string, error) {
	userClaim.IssuedAt = jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(bearerToken string, secret []byte) (UserClaim, error) {
	token, err := jwt.ParseWithClaims(bearerToken, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return UserClaim{}, fmt.Errorf("failed to parse bearerToken: %w", err)
	}

	if claims, ok := token.Claims.(*UserClaim); ok && token.Valid {
		return *claims, nil
	} else {
		return UserClaim{}, fmt.Errorf("unknown claims type")
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
		return false
	}

	return true
}
