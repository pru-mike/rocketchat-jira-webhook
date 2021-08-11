package rocketchat

import (
	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/pru-mike/rocketchat-jira-webhook/jira"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJiraOutputBuilder_Color(t *testing.T) {
	b := JiraOutputBuilder{
		priorityToColor: map[int]string{
			1: "Red",
			2: "Green",
			3: "Blue",
		},
		cfg: &config.MessageJira{
			PriorityColors: false,
		},
		OutputBuilder: &OutputBuilder{
			Message: &config.Message{
				DefaultColor: "Black",
			},
		},
	}
	assert.Equal(t, "Black", b.Color(&jira.Priority{ID: 1}))

	b.cfg.PriorityColors = true

	assert.Equal(t, "Red", b.Color(&jira.Priority{ID: 1}))
	assert.Equal(t, "Red", b.Color(&jira.Priority{ID: 1}))
	assert.Equal(t, "Black", b.Color(&jira.Priority{ID: 123}))

}
