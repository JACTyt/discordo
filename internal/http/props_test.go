package http

import (
	"encoding/base64"
	"encoding/json"

	"github.com/diamondburned/arikawa/v3/gateway"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifyProperties(t *testing.T) {
	identifyProperties := IdentifyProperties()
	assert.Equal(t, "Windows", identifyProperties["os"])
	assert.Equal(t, "Chrome", identifyProperties["browser"])
	assert.NotEmpty(t, identifyProperties["client_launch_id"])
}

func Test_superProps(t *testing.T) {
	s, err := superProps()
	assert.NoError(t, err)

	data, err := base64.StdEncoding.DecodeString(s)
	assert.NoError(t, err)

	var props gateway.IdentifyProperties
	err = json.Unmarshal(data, &props)
	assert.NoError(t, err)

	assert.NotContains(t, props, "is_fast_connect")
	assert.NotContains(t, props, "gateway_connect_reasons")
}
