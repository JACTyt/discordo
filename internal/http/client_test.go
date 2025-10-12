package http_test

import (
	"testing"

	"github.com/ayn2op/discordo/internal/http"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	token := "dummy"
	client := http.NewClient(token)

	assert.NotNil(t, client, "NewClient should return a non-nil client")
	// If tokens are equal, then token was set correctly
	assert.Equal(t, client.Token, token, "Token from api should be equal with input token")
}
