package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStripSlash(t *testing.T) {
	testStr := ""
	stripSlash(&testStr)
	assert.Equal(t, testStr, "")

	testStr = "https://google.com"
	stripSlash(&testStr)
	assert.Equal(t, testStr, "https://google.com")

	testStr = "https://google.com/"
	stripSlash(&testStr)
	assert.Equal(t, testStr, "https://google.com")
}

func TestContains(t *testing.T) {
	assert.Equal(t, contains("xxx", []string{}), false)
	assert.Equal(t, contains("xxx", []string{"aaa","bbb","ccc"}), false)
	assert.Equal(t, contains("xxx", []string{"aaa","bbb","ccc", "xxx"}), true)
}