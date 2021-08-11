package config

import (
	"fmt"
	"github.com/pru-mike/rocketchat-jira-webhook/utils"
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (
	Assignee   string = "assignee"
	Status     string = "status"
	Reporter   string = "reporter"
	Creator    string = "creator"
	Priority   string = "priority"
	Resolution string = "resolution"
	Type       string = "type"
	Created    string = "created"
	Updated    string = "updated"
	Watches    string = "watches"
	Components string = "components"
	Labels     string = "labels"
)

var JiraAllFields = [...]string{
	Assignee, Status, Reporter, Creator, Priority, Resolution, Type, Created, Updated, Watches, Components, Labels,
}
var JiraDefaultFields = [...]string{
	Priority, Type, Status, Resolution, Assignee, Reporter, Created, Updated,
}

type Jira struct {
	URL           string `validate:"required"`
	Username      string `validate:"required"`
	Password      string `validate:"required"`
	Timeout       time.Duration
	requestFields []string
	FindKeyRegexp string `mapstructure:"find_keys_regexp"`
}

func (j *Jira) RequestFields() []string {
	return j.requestFields
}

type MessageJira struct {
	Message              `mapstructure:",squash"`
	PriorityColors       bool     `mapstructure:"priority_colors"`
	PriorityIDPrecedence []int    `mapstructure:"priority_id_precedence"`
	SortByPrecedence     bool     `mapstructure:"sort_by_precedence"`
	ColorsByPriority     []string `mapstructure:"colors_by_priority"`
	InactiveAuthor       string   `mapstructure:"inactive_author"`
	InactiveAuthorIcons  []string `mapstructure:"inactive_author_icons"`
}

func setDefaultsJira() {
	setDefaultsMessage("message_jira")
	viper.SetDefault("message_jira.max_text_length", 600)
	viper.SetDefault("message_jira.priority_colors", true)
	viper.SetDefault("message_jira.default_color", "#205081")
	viper.SetDefault("message_jira.fields", JiraDefaultFields[:])
	viper.SetDefault("message_jira.priority_id_precedence", []int{1, 2, 3, 4, 5})
	viper.SetDefault("message_jira.colors_by_priority", []string{"#000000", "#ff5500", "#F6C342", "#00ff66", "#0095ff"})
	viper.SetDefault("message_jira.title_template", "{{.GetKey}} {{.GetSummary}}")
	viper.SetDefault("message_jira.author_template", "{{ .Reporter.DisplayName }}")
	viper.SetDefault("message_jira.inactive_author", "reporter")
	viper.SetDefault("message_jira.inactive_author_icons", []string{"candle"})
	viper.SetDefault("message_jira.sort_by_precedence", true)
}

func preProcJira(config *Config) error {

	utils.StripSlash(&config.Jira.URL)

	for i := range config.MessageJira.Fields {
		config.MessageJira.Fields[i] = strings.ToLower(config.MessageJira.Fields[i])
		field := config.MessageJira.Fields[i]
		if !utils.Contains(field, JiraAllFields[:]) {
			return fmt.Errorf("invalid configuration field: %s", field)
		}
		config.Jira.requestFields = append(config.Jira.requestFields, strings.ToLower(field))
	}
	if config.MessageJira.PriorityColors && !utils.Contains(Priority, config.MessageJira.Fields) {
		config.Jira.requestFields = append(config.Jira.requestFields, Priority)
	}

	if !config.MessageJira.ShowAuthor {
		config.MessageJira.AuthorTemplate = ""
		config.MessageJira.AuthorIcons = []string{}
	} else if !utils.Contains(Reporter, config.MessageJira.Fields) {
		config.Jira.requestFields = append(config.Jira.requestFields, Reporter)
	}

	if config.Jira.Timeout == 0 {
		config.Jira.Timeout, _ = time.ParseDuration("30s")
	}

	return nil
}
