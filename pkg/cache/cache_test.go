package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKey(t *testing.T) {
	Init("test.json")

	k1 := NewKey[string]("k1")
	k2 := NewKey[int]("k2")
	k3 := NewKey[bool]("k3")

	k1.Set("v1")
	k2.Set(2)
	k3.Set(true)

	assert.Equal(t, "v1", k1.Get())
	assert.Equal(t, 2, k2.Get())
	assert.Equal(t, true, k3.Get())
}
