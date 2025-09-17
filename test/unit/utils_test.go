package unit

import (
	"testing"

	"github.com/AtlasOpx/devprep/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hashedPassword, err := utils.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)
	assert.True(t, len(hashedPassword) > 50)
}

func TestCheckPasswordHash_Valid(t *testing.T) {
	password := "testpassword123"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	isValid := utils.CheckPasswordHash(password, hashedPassword)

	assert.True(t, isValid)
}

func TestCheckPasswordHash_Invalid(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	isValid := utils.CheckPasswordHash(wrongPassword, hashedPassword)

	assert.False(t, isValid)
}

func TestGenerateSessionToken(t *testing.T) {
	token1, _ := utils.GenerateSessionToken()
	token2, _ := utils.GenerateSessionToken()

	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)
	assert.NotEqual(t, token1, token2)
	assert.True(t, len(token1) > 20)
	assert.True(t, len(token2) > 20)
}

func TestGenerateSessionToken_Uniqueness(t *testing.T) {
	tokens := make(map[string]bool)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		token, _ := utils.GenerateSessionToken()
		assert.False(t, tokens[token], "Generated duplicate token: %s", token)
		tokens[token] = true
	}

	assert.Equal(t, iterations, len(tokens))
}
