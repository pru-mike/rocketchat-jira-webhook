package rocketchat

import (
	"encoding/json"
	"time"
)

type Bot bool

type Input struct {
	Token       string    `json:"token" validate:"token"`
	ChannelID   string    `json:"channel_id"`
	ChannelName string    `json:"channel_name"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	Text        string    `json:"text" validate:"required"`
	MessageID   string    `json:"message_id"`
	SiteURL     string    `json:"siteUrl"`
	Timestamp   time.Time `json:"timestamp"`
	IsEdited    bool      `json:"IsEdited" validate:"edits"`
	Bot         Bot       `json:"bot" validate:"bots"`
	Alias       string    `json:"alias"`
	TriggerWord string    `json:"trigger_word"`
}

func (bot *Bot) UnmarshalJSON(b []byte) error {
	var inputVal bool
	if err := json.Unmarshal(b, &inputVal); err != nil {
		inputVal = true
	}
	*bot = Bot(inputVal)
	return nil
}
