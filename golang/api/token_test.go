package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestULIDGen(t *testing.T) {
	_, err := NewULIDNow()
	assert.Nil(t, err)
}
