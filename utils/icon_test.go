package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNextElem(t *testing.T) {
	assert.Equal(t, "", getNextElem(nil, 100))
	assert.Equal(t, "", getNextElem([]string{}, 0))
	assert.Equal(t, "", getNextElem([]string{}, 10))
	assert.Equal(t, "aaa", getNextElem([]string{"aaa"}, 0))
	assert.Equal(t, "aaa", getNextElem([]string{"aaa"}, 10))
	assert.Equal(t, "aaa", getNextElem([]string{"aaa", "bbb"}, 0))
	assert.Equal(t, "bbb", getNextElem([]string{"aaa", "bbb"}, 1))
	assert.Equal(t, "aaa", getNextElem([]string{"aaa", "bbb"}, 2))
	assert.Equal(t, "bbb", getNextElem([]string{"aaa", "bbb"}, 3))
}
