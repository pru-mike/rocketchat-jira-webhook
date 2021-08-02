package assets

import (
	_ "embed"
	"encoding/json"
)

//go:embed data.json
var quotesJSON []byte

var Quotes []struct {
	Author string
	Quote  string
}

func init() {
	err := json.Unmarshal(quotesJSON, &Quotes)
	if err != nil {
		return
	}
}
