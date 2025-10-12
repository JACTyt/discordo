package http_test

import (
	"testing"

	"github.com/ayn2op/discordo/internal/http"
	"github.com/stretchr/testify/assert"
)

func TestHeaders(t *testing.T) {
	headers := http.Headers()

	// Check mandatory fields set
	assert.Equal(t, "*/*", headers.Get("Accept"))
	assert.Equal(t, "gzip, deflate, br, zstd", headers.Get("Accept-Encoding"))
	assert.Equal(t, "en-US,en;q=0.9", headers.Get("Accept-Language"))
	assert.Equal(t, "https://discord.com", headers.Get("Origin"))
	assert.Equal(t, "https://discord.com/channels/@me", headers.Get("Referer"))

	// Check optional fields set
	assert.NotEmpty(t, headers.Get("X-Discord-Locale"))
	assert.NotEmpty(t, headers.Get("X-Super-Properties"))

	// Check other fields
	assert.Equal(t, "u=0, i", headers.Get("Priority"))
	assert.Equal(t, "empty", headers.Get("Sec-Fetch-Dest"))
	assert.Equal(t, "cors", headers.Get("Sec-Fetch-Mode"))
	assert.Equal(t, "same-origin", headers.Get("Sec-Fetch-Site"))
	assert.Equal(t, "bugReporterEnabled", headers.Get("X-Debug-Options"))
}
