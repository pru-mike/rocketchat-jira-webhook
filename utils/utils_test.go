package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStripSlash(t *testing.T) {
	testStr := ""
	StripSlash(&testStr)
	assert.Equal(t, testStr, "")

	testStr = "https://google.com"
	StripSlash(&testStr)
	assert.Equal(t, testStr, "https://google.com")

	testStr = "https://google.com/"
	StripSlash(&testStr)
	assert.Equal(t, testStr, "https://google.com")
}

func TestContains(t *testing.T) {
	assert.Equal(t, Contains("xxx", []string{}), false)
	assert.Equal(t, Contains("xxx", []string{"aaa", "bbb", "ccc"}), false)
	assert.Equal(t, Contains("xxx", []string{"aaa", "bbb", "ccc", "xxx"}), true)
}
