package rocketchat

import (
	"testing"

	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/stretchr/testify/assert"
)

func TestOutputBuilder_TrimMaxLen(t *testing.T) {
	b := OutputBuilder{
		Message: &config.Message{
			MaxTextLen: 5,
		},
	}
	assert.Equal(t, "test", b.TrimMaxLen("test"))
	assert.Equal(t, "test5", b.TrimMaxLen("test5"))
	assert.Equal(t, "test5…", b.TrimMaxLen("test51"))
	assert.Equal(t, "ТЕСТ", b.TrimMaxLen("ТЕСТ"))
	assert.Equal(t, "ТЕСТТ…", b.TrimMaxLen("ТЕСТТЕСТТЕСТ"))

	b = OutputBuilder{
		Message: &config.Message{
			MaxTextLen: 2,
		},
	}
	assert.Equal(t, "世界", b.TrimMaxLen("世界"))
	assert.Equal(t, "世界…", b.TrimMaxLen("世界世界"))

	b = OutputBuilder{
		Message: &config.Message{
			MaxTextLen: 1,
		},
	}
	assert.Equal(t, "世…", b.TrimMaxLen("世界世界"))

}
