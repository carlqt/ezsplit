package auth

import (
	"testing"

	"github.com/carlqt/ezsplit/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestValidateBearerToken(t *testing.T) {
	secretKey := []byte("secret")

	t.Run("when bearerToken is not in the correct format", func(t *testing.T) {
		_, err := ValidateBearerToken("", secretKey)

		assert.ErrorContains(t, err, "failed to parse bearerToken")
	})

	t.Run("when bearerToken is empty", func(t *testing.T) {
		_, err := ValidateBearerToken("", secretKey)

		assert.ErrorContains(t, err, "failed to parse bearerToken")
	})

	t.Run("when bearerToken is valid", func(t *testing.T) {
		claim := NewUserClaim(
			int32(5),
			"username",
			repository.Verified,
		)

		accessToken, _ := CreateAndSignToken(claim, secretKey)

		result, err := ValidateBearerToken(accessToken, secretKey)

		if assert.Nil(t, err) {
			assert.Equal(t, "username", result.Username)
			assert.Equal(t, "5", result.ID)
		}
	})

	t.Run("when bearerToken was signed with a different secret key", func(t *testing.T) {
		claim := NewUserClaim(
			int32(5),
			"username",
			repository.Verified,
		)

		accessToken, _ := CreateAndSignToken(claim, []byte("notSoSecretKey"))

		_, err := ValidateBearerToken(accessToken, secretKey)

		assert.ErrorContains(t, err, "failed to parse bearerToken")
	})
}
