package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidFormat(t *testing.T) {
	assert.False(t, isValidFormat("toml"))
	assert.True(t, isValidFormat("json"))
}
