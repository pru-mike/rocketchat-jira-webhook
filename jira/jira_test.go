package jira

import (
	"strings"
	"testing"
)

func TestParseKeys(t *testing.T) {
	tests := []struct {
		request string
		result  []string
	}{
		{
			"no jira keys",
			[]string{},
		},
		{
			"one jira TEST-123 key",
			[]string{"TEST-123"},
		},
		{
			"dublicate keys TEST-123 TEST-123",
			[]string{"TEST-123"},
		},
		{
			"TEST-123 one jira key",
			[]string{"TEST-123"},
		},
		{
			"one jira key TEST-123",
			[]string{"TEST-123"},
		},
		{
			"one jira TEST-123! key",
			[]string{"TEST-123"},
		},
		{
			"one jira <TEST-123> key",
			[]string{"TEST-123"},
		},
		{
			"several jira TEST-123 key PROJ-18,XCD-21",
			[]string{"TEST-123", "PROJ-18", "XCD-21"},
		},
		{
			"https://mycompany.jira.com/browse/TEST-123",
			[]string{"TEST-123"},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.request, func(t *testing.T) {
			result := ParseKeys(test.request)
			if strings.Join(result, ",") != strings.Join(test.result, ",") {
				t.Errorf("got %q, whant %q", result, test.result)
			}
		})
	}
}

func TestStripKey(t *testing.T) {
	tests := []struct {
		request string
		result  string
	}{
		{
			"DEVELOP-123",
			"DEVELOP",
		},
		{
			"PRJ-0",
			"PRJ",
		},
		{
			"TEST",
			"TEST",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.request, func(t *testing.T) {
			result := StripKey(test.request)
			if result != test.result {
				t.Errorf("got %q, whant %q", result, test.result)
			}
		})
	}
}
