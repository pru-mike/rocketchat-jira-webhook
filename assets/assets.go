package assets

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
	"math/rand"
	"strings"
)

//go:embed data.json
var quotesRaw []byte

//go:embed *.png
var logo embed.FS

type Quote struct {
	Author string
	Quote  string
}

var quotes []Quote

func GetQuoteWithProb(prob float32) (bool, *Quote) {
	if len(quotes) > 0 {
		if prob > rand.Float32() {
			q := quotes[rand.Intn(len(quotes))]
			return true, &q
		}
	}
	return false, nil
}

type logoContainer map[string]string

func (l *logoContainer) add(name string, logo []byte) {
	if len(logo) > 0 {
		b := strings.Builder{}
		b.WriteString("data:image/png;base64,")
		b.WriteString(base64.StdEncoding.EncodeToString(logo))
		(*l)[name] = b.String()
	}
}

var logoTypes logoContainer

func addLogo(name string, logo []byte) {
	if logoTypes == nil {
		logoTypes = make(map[string]string)
	}
	logoTypes.add(name, logo)
}

func GetLogo(name string) (string, bool) {
	logo, ok := logoTypes[name]
	return logo, ok
}

func loadQuotes() {
	err := json.Unmarshal(quotesRaw, &quotes)
	if err != nil {
		logger.Error(err)
	}
}

func loadLogoTypes() {
	dir, err := logo.ReadDir(".")
	if err != nil {
		logger.Error(err)
	}
	for _, dir := range dir {
		name := dir.Name()
		if len(name) > 4 && name[len(name)-4:] == ".png" {
			data, err := logo.ReadFile(dir.Name())
			if err != nil {
				logger.Error(err)
				continue
			}
			addLogo(name[:len(name)-4], data)
		}
	}
}

func init() {
	loadQuotes()
	loadLogoTypes()
}
