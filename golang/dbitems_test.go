package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("foo", "bar")
	assert.Nil(t, err)
	assert.Equal(t, "foo", user.Username)
	assert.NotEqual(t, "bar", user.Password)
}
