package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/pru-mike/rocketchat-jira-webhook/assets"
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
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

type Config struct {
	App        App
	Jira       Jira
	Rocketchat Rocketchat
	Message    Message
}

type App struct {
	Host        string
	Port        int
	LogLevel    string `mapstructure:"log_level"`
	ErrToRocket bool   `mapstructure:"err_to_rocket"`
}

type Jira struct {
	URL           string `validate:"required"`
	Username      string `validate:"required"`
	Password      string `validate:"required"`
	Timeout       time.Duration
	requestFields []string
}

func (j *Jira) RequestFields() []string {
	return j.requestFields
}

type Rocketchat struct {
	Tokens                 []string
	WhitelistedUsers       []string `mapstructure:"whitelisted_users"`
	BlacklistedUsers       []string `mapstructure:"blacklisted_users"`
	WhitelistedChannels    []string `mapstructure:"whitelisted_channels"`
	BlacklistedChannels    []string `mapstructure:"blacklisted_channels"`
	WhitelistedProjectKeys []string `mapstructure:"whitelisted_project_keys"`
	BlacklistedProjectKeys []string `mapstructure:"blacklisted_project_keys"`
	AllowEdits             bool     `mapstructure:"allow_edits"`
	AllowBots              bool     `mapstructure:"allow_bots"`
}

type Message struct {
	Username             string   `mapstructure:"username"`
	IconURL              string   `mapstructure:"icon_url"`
	MaxTextLen           int      `mapstructure:"max_text_length"`
	PriorityColors       bool     `mapstructure:"priority_colors"`
	DefaultColor         string   `mapstructure:"default_color"`
	Fields               []string `mapstructure:"fields"`
	UseRealNames         bool     `mapstructure:"use_real_names"`
	DatetimeLayout       string   `mapstructure:"datetime_layout"`
	PriorityIDPrecedence []int    `mapstructure:"priority_id_precedence"`
	ColorsByPriority     []string `mapstructure:"colors_by_priority"`
	MsgLang              string   `mapstructure:"msg_lang"`
	QuoteProbability     float32  `mapstructure:"quote_prob"`
	UnescapeHTML         bool     `mapstructure:"unescape_html"`
}

func (m *Message) LangTag() language.Tag {
	t, err := language.Parse(m.MsgLang)
	if err != nil {
		logger.Errorf("can't parse lang tag '%s': %v", m.MsgLang, err)
		return language.English
	}
	return t
}

func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.App.Host, c.App.Port)
}

func Load(configFile string) (*Config, error) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath("/etc/github.com/pru-mike/rocketchat-jira-webhook")
		viper.AddConfigPath(".")
	}
	viper.SetDefault("app.host", "0.0.0.0")
	viper.SetDefault("app.port", "4567")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.err_to_rocket", true)
	viper.SetDefault("jira.timeout", "30s")
	viper.SetDefault("message.max_text_len", 600)
	viper.SetDefault("message.priority_colors", true)
	viper.SetDefault("message.default_color", "#205081")
	viper.SetDefault("message.fields", JiraDefaultFields[:])
	viper.SetDefault("message.use_real_names", true)
	viper.SetDefault("message.datetime_layout", "02/01/2006 15:04")
	viper.SetDefault("message.priority_id_precedence", []int{1, 2, 3, 4, 5})
	viper.SetDefault("message.colors_by_priority", []string{"#FF4437", "#D04437", "#E3833C", "#F6C342", "#707070"})
	viper.SetDefault("message.msg_lang", "en")
	viper.SetDefault("message.quote_prob", 0.009)
	viper.SetDefault("message.unescape_html", true)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("can't read configuration file: %w", err)
	}

	var config Config

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	validate := validator.New()

	err = validate.Struct(config)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration file: %w", err)
	}

	stripSlash(&config.Jira.URL)

	for i := range config.Message.Fields {
		config.Message.Fields[i] = strings.ToLower(config.Message.Fields[i])
		field := config.Message.Fields[i]
		if !contains(field, JiraAllFields[:]) {
			return nil, fmt.Errorf("invalid configuration field: %s", field)
		}
		config.Jira.requestFields = append(config.Jira.requestFields, strings.ToLower(field))
	}
	if config.Message.PriorityColors && !contains(Priority, config.Message.Fields) {
		config.Jira.requestFields = append(config.Jira.requestFields, Priority)
	}

	if logo, ok := assets.GetLogo(config.Message.IconURL); ok {
		config.Message.IconURL = logo
	}

	return &config, nil
}

func stripSlash(str *string) {
	if len(*str) > 0 && (*str)[len(*str)-1] == '/' {
		*str = (*str)[:len(*str)-1]
	}
}

func contains(target string, list []string) bool {
	for _, val := range list {
		if target == val {
			return true
		}
	}
	return false
}
