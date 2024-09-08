package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func NullableInt64(n int) *int64 {
	result := int64(n)

	return &result
}

func TestValidateJWT(t *testing.T) {
	secretKey := []byte("secret")

	t.Run("when bearerToken is not in the correct format", func(t *testing.T) {
		_, err := ValidateJWT("", secretKey)

		assert.ErrorContains(t, err, "failed to parse bearerToken")
	})

	t.Run("when bearerToken is empty", func(t *testing.T) {
		_, err := ValidateJWT("", secretKey)

		assert.ErrorContains(t, err, "failed to parse bearerToken")
	})

	t.Run("when bearerToken is valid", func(t *testing.T) {
		claim := NewUserClaim(
			int64(5),
			"username",
			true,
		)

		accessToken, _ := CreateAndSignToken(claim, secretKey)

		result, err := ValidateJWT(accessToken, secretKey)

		if assert.Nil(t, err) {
			assert.Equal(t, "username", result.Username)
			assert.Equal(t, "5", result.ID)
		}
	})

	t.Run("when bearerToken was signed with a different secret key", func(t *testing.T) {
		claim := NewUserClaim(
			int64(5),
			"username",
			false,
		)

		accessToken, _ := CreateAndSignToken(claim, []byte("notSoSecretKey"))

		_, err := ValidateJWT(accessToken, secretKey)

		assert.ErrorContains(t, err, "failed to parse bearerToken")
	})
}
