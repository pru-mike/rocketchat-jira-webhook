package rocketchat

import (
	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/pru-mike/rocketchat-jira-webhook/jira"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrim(t *testing.T) {
	b := OutputBuilder{
		Message: &config.Message{
			MaxTextLen: 5,
		},
	}
	assert.Equal(t, "test", b.trim("test"))
	assert.Equal(t, "test5", b.trim("test5"))
	assert.Equal(t, "test5…", b.trim("test51"))
	assert.Equal(t, "ТЕСТ", b.trim("ТЕСТ"))
	assert.Equal(t, "ТЕСТТ…", b.trim("ТЕСТТЕСТТЕСТ"))

	b = OutputBuilder{
		Message: &config.Message{
			MaxTextLen: 2,
		},
	}

	assert.Equal(t, "世界", b.trim("世界"))
	assert.Equal(t, "世界…", b.trim("世界世界"))
}

func TestColor(t *testing.T) {
	b := OutputBuilder{
		priorityToColor: map[int]string{
			1: "Red",
			2: "Green",
			3: "Blue",
		},
		Message: &config.Message{
			PriorityColors: false,
			DefaultColor:   "Black",
		},
	}
	assert.Equal(t, "Black", b.color(&jira.Priority{ID: 1}))

	b.Message.PriorityColors = true

	assert.Equal(t, "Red", b.color(&jira.Priority{ID: 1}))
	assert.Equal(t, "Red", b.color(&jira.Priority{ID: 1}))
	assert.Equal(t, "Black", b.color(&jira.Priority{ID: 123}))

}
