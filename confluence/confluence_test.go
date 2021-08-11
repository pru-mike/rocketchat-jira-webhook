package confluence

import (
	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"reflect"
	"strings"
	"testing"
)

func TestConfluence_FindPagesViewIDs(t *testing.T) {
	tests := []struct {
		url    string
		text   string
		result []string
	}{
		{
			url:    "https://my.confluence.com",
			text:   "",
			result: []string{},
		},
		{
			url:    "https://my.confluence.com",
			text:   "test123123",
			result: []string{},
		},
		{
			url:    "https://my.confluence.com",
			text:   "https://my.confluence.com?aaaaaa",
			result: []string{},
		},
		{
			url:    "https://my.confluence.com",
			text:   "https://my.confluence.com/pages/viewpage.action?pageId=123321987",
			result: []string{"123321987"},
		},
		{
			url:    "https://my.confluence.com",
			text:   "https://my.yyyyy.com/pages/viewpage.action?pageId=123321987",
			result: []string{},
		},
		{
			url:    "https://my.confluence.com",
			text:   "https://my.confluence.com/pages/viewpage.action?pageId=123321987 https://my.confluence.com/pages/viewpage.action?pageId=31244521",
			result: []string{"123321987", "31244521"},
		},
		{
			url: "https://my.confluence.com",
			text: `https://my.confluence.com/pages/viewpage.action?pageId=123321987 
https://my.confluence.com/pages/viewpage.action?pageId=31244521
https://my.confluence.com/pages/viewpage.action?pageId=123321987
https://my.confluence.com/pages/viewpage.action?pageId=1231231534`,
			result: []string{"123321987", "31244521", "1231231534"},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.text, func(t *testing.T) {
			c := NewClient(&config.Confluence{URL: test.url})
			result := c.FindPagesViewIDs(test.text)
			if strings.Join(result, ",") != strings.Join(test.result, ",") {
				t.Errorf("got %q, whant %q", result, test.result)
			}
		})
	}

}

func TestConfluence_FindPagesSpaceTitles(t *testing.T) {
	tests := []struct {
		url    string
		text   string
		result []SpaceTitle
	}{
		{
			url:    "https://my.confluence.com",
			text:   "",
			result: []SpaceTitle{},
		},
		{
			url:    "https://my.confluence.com",
			text:   "test123123",
			result: []SpaceTitle{},
		},
		{
			url:    "https://my.confluence.com",
			text:   "https://my.confluence.com?aaaaaa",
			result: []SpaceTitle{},
		},
		{
			url:  "https://my.confluence.com",
			text: "https://my.confluence.com/display/TST/myPage%20Title",
			result: []SpaceTitle{
				{
					Space: "TST",
					Title: "myPage%20Title",
				},
			},
		},
		{
			url:  "https://my.confluence.com",
			text: "https://my.confluence.com/display/TST/myPage%20Title https://my.confluence.com/display/TST/myPage%20Title",
			result: []SpaceTitle{
				{
					Space: "TST",
					Title: "myPage%20Title",
				},
			},
		},
		{
			url: "https://my.confluence.com",
			text: `
https://my.confluence.com/display/TST/myPage%20Title
https://my.confluence.com/display/ZZZ/page123
https://my.confluence.com/display/TST/myPage%20Title
https://my.confluence.com/display/ZZZ/page123
https://my.confluence.com/display/XXXX/aaaaa

`,
			result: []SpaceTitle{
				{
					Space: "TST",
					Title: "myPage%20Title",
				},
				{
					Space: "ZZZ",
					Title: "page123",
				}, {
					Space: "XXXX",
					Title: "aaaaa",
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.text, func(t *testing.T) {
			c := NewClient(&config.Confluence{URL: test.url})
			result := c.FindPagesSpaceTitles(test.text)
			if !reflect.DeepEqual(result, test.result) {
				t.Errorf("got %q, whant %q", result, test.result)
			}
		})
	}

}
