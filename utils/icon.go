package utils

import (
	"github.com/pru-mike/rocketchat-jira-webhook/assets"
	"math/rand"
)

func getNextElem(src []string, n uint) string {
	if len(src) > 0 {
		return src[n%uint(len(src))]
	}
	return ""
}

func getNextLogo(logos []string, n int) string {
	logo, _ := assets.GetLogo(getNextElem(logos, uint(n)))
	return logo
}

func NextIconGetter(icons []string) func() string {
	if len(icons) == 1 {
		return func() string {
			return getNextLogo(icons, 0)
		}
	}
	iconsCopy := make([]string, len(icons))
	copy(iconsCopy, icons)
	rand.Shuffle(len(iconsCopy), func(i, j int) {
		iconsCopy[i], iconsCopy[j] = iconsCopy[j], iconsCopy[i]
	})
	i := rand.Intn(len(iconsCopy))
	return func() string {
		i++
		return getNextLogo(iconsCopy, i)
	}
}
