package rocketchat

import (
	"encoding/json"
	"testing"

	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/stretchr/testify/assert"
)

func TestValidateToken(t *testing.T) {
	v := SetupValidator(&config.Rocketchat{
		Tokens: []string{"testoken"},
	})
	assert.NoError(t, v.Struct(Input{
		Text:  "text",
		Token: "testoken",
	}))
	assert.Error(t, v.Struct(Input{
		Text:  "text",
		Token: "aaaa",
	}))
}

func TestValidateEdits(t *testing.T) {
	v := SetupValidator(&config.Rocketchat{
		AllowEdits: true,
	})
	assert.NoError(t, v.Struct(Input{
		Text:     "text",
		IsEdited: true,
	}))
	assert.NoError(t, v.Struct(Input{
		Text:     "text",
		IsEdited: false,
	}))

	v = SetupValidator(&config.Rocketchat{
		AllowEdits: false,
	})
	assert.Error(t, v.Struct(Input{
		Text:     "text",
		IsEdited: true,
	}))
	assert.NoError(t, v.Struct(Input{
		Text:     "text",
		IsEdited: false,
	}))
}

func TestValidateBots(t *testing.T) {
	var human, bot Input
	var err error

	err = json.Unmarshal([]byte(`{
		"text":"text",
		"bot":{"i":"HhXyzyXbFFGcEyGFM"}
	}`), &bot)
	assert.NoError(t, err)

	err = json.Unmarshal([]byte(`{
		"text":"text",
		"bot":false
	}`), &human)
	assert.NoError(t, err)

	v := SetupValidator(&config.Rocketchat{
		AllowBots: true,
	})
	assert.NoError(t, v.Struct(human))
	assert.NoError(t, v.Struct(bot))

	v = SetupValidator(&config.Rocketchat{
		AllowBots: false,
	})
	assert.NoError(t, v.Struct(human))
	assert.Error(t, v.Struct(bot))
}

func TestValidateUserStructLevel(t *testing.T) {
	v := SetupValidator(&config.Rocketchat{
		BlacklistedUsers: []string{"NAME", "ID"},
	})
	assert.NoError(t, v.Struct(Input{
		Text:     "text",
		UserName: "name",
		UserID:   "id",
	}))
	assert.Error(t, v.Struct(Input{
		Text:     "text",
		UserName: "NAME",
		UserID:   "id",
	}))
	assert.Error(t, v.Struct(Input{
		Text:     "text",
		UserName: "name",
		UserID:   "ID",
	}))

	v = SetupValidator(&config.Rocketchat{
		BlacklistedUsers: []string{"NAME", "ID"},
		WhitelistedUsers: []string{"NAME", "ID"},
	})
	assert.Error(t, v.Struct(Input{
		Text:     "text",
		UserName: "NAME",
		UserID:   "id",
	}))
	assert.Error(t, v.Struct(Input{
		Text:     "text",
		UserName: "name",
		UserID:   "ID",
	}))

	v = SetupValidator(&config.Rocketchat{
		BlacklistedUsers: []string{"NAME", "ID"},
		WhitelistedUsers: []string{"NAME2"},
	})
	assert.Error(t, v.Struct(Input{
		Text:     "text",
		UserName: "testName",
		UserID:   "testID",
	}))
	assert.Error(t, v.Struct(Input{
		Text:     "text",
		UserName: "name",
		UserID:   "ID",
	}))
	assert.NoError(t, v.Struct(Input{
		Text:     "text",
		UserName: "NAME2",
		UserID:   "xxxx",
	}))
	assert.Error(t, v.Struct(Input{
		Text:     "text",
		UserName: "NAME2",
		UserID:   "ID",
	}))
}

func TestValidateKeys(t *testing.T) {
	v := SetupValidator(&config.Rocketchat{
		WhitelistedJiraProjectKeys: []string{},
		BlacklistedJiraProjectKeys: []string{},
	})
	assert.Equal(t, []string{}, v.ValidateJiraKeys([]string{}))
	assert.Equal(t, []string{"PRJ-123", "TST-111", "ZZZ-31"}, v.ValidateJiraKeys([]string{"PRJ-123", "TST-111", "ZZZ-31"}))

	v = SetupValidator(&config.Rocketchat{
		WhitelistedJiraProjectKeys: []string{},
		BlacklistedJiraProjectKeys: []string{"ZZZ"},
	})
	assert.Equal(t, []string{"PRJ-123", "TST-111"}, v.ValidateJiraKeys([]string{"PRJ-123", "TST-111", "ZZZ-31"}))

	v = SetupValidator(&config.Rocketchat{
		WhitelistedJiraProjectKeys: []string{"ZZZ"},
		BlacklistedJiraProjectKeys: []string{},
	})
	assert.Equal(t, []string{"ZZZ-31"}, v.ValidateJiraKeys([]string{"PRJ-123", "TST-111", "ZZZ-31"}))

	v = SetupValidator(&config.Rocketchat{
		WhitelistedJiraProjectKeys: []string{"ZZZ", "PRJ"},
		BlacklistedJiraProjectKeys: []string{"ZZZ"},
	})
	assert.Equal(t, []string{"PRJ-123", "PRJ-333"}, v.ValidateJiraKeys([]string{"PRJ-123", "TST-111", "PRJ-333", "ZZZ-31"}))
}
