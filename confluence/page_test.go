package confluence

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	jt := struct {
		DateTime Time `json:"time"`
	}{}
	dt := []byte(`{"time":"2020-06-01T16:07:48.328+03:00"}`)
	err := json.Unmarshal(dt, &jt)
	assert.NoError(t, err)
	assert.Equal(t, "2020-06-01 16:07:48.328 +0300 MSK", time.Time(jt.DateTime).String())
}
