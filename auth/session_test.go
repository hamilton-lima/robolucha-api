package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsEmpty(t *testing.T) {
	data := make([]string, 0)
	c := contains(data, "foo")
	assert.Equal(t, false, c)
}

func TestContainsTrue(t *testing.T) {
	data := []string{"one", "two", "three"}
	c := contains(data, "two")
	assert.Equal(t, true, c)
}
