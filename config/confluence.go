package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/pru-mike/rocketchat-jira-webhook/assets"
	"github.com/pru-mike/rocketchat-jira-webhook/utils"
	"github.com/spf13/viper"
)

const (
	CreatedBy         string = "created_by"
	CreatedDate       string = "created_date"
	UpdatedBy         string = "updated_by"
	UpdatedDate       string = "updated_date"
	Latest            string = "latest"
	ConfluenceStatus  string = "status"
	LastVersionNumber string = "last_version"
	SpaceName         string = "space_name"
)

var ConfluenceAllFields = [...]string{
	SpaceName, CreatedBy, CreatedDate, UpdatedBy, UpdatedDate, LastVersionNumber, ConfluenceStatus, Latest,
}
var ConfluenceDefaultFields = [...]string{
	SpaceName, CreatedBy, CreatedDate, UpdatedBy, UpdatedDate, LastVersionNumber, ConfluenceStatus, Latest,
}

const (
	StyledView          string = "styled_view"
	AnonymousExportView string = "anonymous_export_view"
	Storage             string = "storage"
	ExportView          string = "export_view"
	View                string = "view"
	Editor              string = "editor"
)

var ConfluenceBodyExpand = [...]string{
	StyledView, AnonymousExportView, Storage, ExportView, View, Editor,
}

type Confluence struct {
	URL                   string `validate:"required"`
	Username              string `validate:"required"`
	Password              string `validate:"required"`
	Timeout               time.Duration
	FindPagesByViewID     string `mapstructure:"find_pages_viewid"`
	FindPagesBySpaceTitle string `mapstructure:"find_pages_spacetitle"`
	BodyExpand            string `mapstructure:"body_expand"`
}

type MessageConfluence struct {
	Message `mapstructure:",squash"`
}

func setDefaultsConfluence() {
	setDefaultsMessage("message_confluence")
	viper.SetDefault("message_confluence.max_text_length", 1800)
	viper.SetDefault("message_confluence.fields", ConfluenceDefaultFields[:])
	viper.SetDefault("message_confluence.title_template", "{{ .Space.Key }} | {{ .Title }}")
	viper.SetDefault("message_confluence.author_template", "{{ .CreatedBy.DisplayName }}")
}

func preProcConfluence(config *Config) error {

	utils.StripSlash(&config.Confluence.URL)

	if logo, ok := assets.GetLogo(config.MessageConfluence.IconURL); ok {
		config.MessageConfluence.IconURL = logo
	}

	for i := range config.MessageConfluence.Fields {
		config.MessageConfluence.Fields[i] = strings.ToLower(config.MessageConfluence.Fields[i])
		field := config.MessageConfluence.Fields[i]
		if !utils.Contains(field, ConfluenceAllFields[:]) {
			return fmt.Errorf("invalid configuration field: %s", field)
		}
	}

	if !config.MessageConfluence.ShowAuthor {
		config.MessageConfluence.AuthorTemplate = ""
		config.MessageConfluence.AuthorIcons = []string{}
	}

	if config.Confluence.Timeout == 0 {
		config.Confluence.Timeout, _ = time.ParseDuration("30s")
	}

	if config.Confluence.BodyExpand == "" {
		config.Confluence.BodyExpand = AnonymousExportView
	}

	if !utils.Contains(config.Confluence.BodyExpand, ConfluenceBodyExpand[:]) {
		return fmt.Errorf("invalid configuration field: %s value: %s", "body_expand", config.Confluence.BodyExpand)
	}

	return nil
}
