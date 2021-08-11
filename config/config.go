package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/pru-mike/rocketchat-jira-webhook/assets"
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

type Config struct {
	App               App
	Jira              *Jira
	Confluence        *Confluence
	MessageJira       MessageJira       `mapstructure:"message_jira"`
	MessageConfluence MessageConfluence `mapstructure:"message_confluence"`
	Rocketchat        Rocketchat
}

func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.App.Host, c.App.Port)
}

type App struct {
	Host        string
	Port        int
	LogLevel    string `mapstructure:"log_level"`
	ErrToRocket bool   `mapstructure:"err_to_rocket"`
}

type Rocketchat struct {
	Tokens                     []string
	WhitelistedUsers           []string `mapstructure:"whitelisted_users"`
	BlacklistedUsers           []string `mapstructure:"blacklisted_users"`
	WhitelistedChannels        []string `mapstructure:"whitelisted_channels"`
	BlacklistedChannels        []string `mapstructure:"blacklisted_channels"`
	WhitelistedJiraProjectKeys []string `mapstructure:"whitelisted_jira_keys"`
	BlacklistedJiraProjectKeys []string `mapstructure:"blacklisted_jira_keys"`
	AllowEdits                 bool     `mapstructure:"allow_edits"`
	AllowBots                  bool     `mapstructure:"allow_bots"`
}

type Message struct {
	Username         string   `mapstructure:"username"`
	IconURL          string   `mapstructure:"icon_url"`
	MaxTextLen       int      `mapstructure:"max_text_length"`
	DefaultColor     string   `mapstructure:"default_color"`
	UseRealNames     bool     `mapstructure:"use_real_names"`
	DatetimeLayout   string   `mapstructure:"datetime_layout"`
	MsgLang          string   `mapstructure:"msg_lang"`
	QuoteProbability float32  `mapstructure:"quote_prob"`
	UnescapeHTML     bool     `mapstructure:"unescape_html"`
	StripTags        bool     `mapstructure:"strip_tags"`
	TitleTemplate    string   `mapstructure:"title_template"`
	ShowAuthor       bool     `mapstructure:"show_author"`
	AuthorTemplate   string   `mapstructure:"author_template"`
	AuthorIcons      []string `mapstructure:"author_icons"`
	Fields           []string `mapstructure:"fields"`
}

func (m *Message) LangTag() language.Tag {
	t, err := language.Parse(m.MsgLang)
	if err != nil {
		logger.Errorf("can't parse lang tag '%s': %v", m.MsgLang, err)
		return language.English
	}
	return t
}

func setDefaults() {
	viper.SetDefault("app.host", "0.0.0.0")
	viper.SetDefault("app.port", "4567")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.err_to_rocket", true)
}

func setDefaultsMessage(message string) {
	viper.SetDefault(message+".use_real_names", true)
	viper.SetDefault(message+".datetime_layout", "02/01/2006 15:04")
	viper.SetDefault(message+".msg_lang", "en")
	viper.SetDefault(message+".quote_prob", 0.009)
	viper.SetDefault(message+".unescape_html", true)
	viper.SetDefault(message+".show_author", true)
	viper.SetDefault(message+".author_icons", []string{
		"stickman-apple", "stickman-bike", "stickman-excercise", "stickman-excercise2",
		"stickman-excercise3", "stickman-heart", "stickman-heart2", "stickman-jump", "stickman-mail",
		"stickman-massage", "stickman-massage2", "stickman-meditation", "stickman-relax", "stickman-run",
		"stickman-sauna", "stickman-shower", "stickman-spa", "stickman-sport", "stickman-sport2",
		"stickman-study", "stickman-swimmer", "stickman-treadmil", "stickman-walker", "stickman-weightlifting",
		"stickman-yoga", "stickman-yoga2", "stickman-yoga3", "stickman-yoga4", "stickman-yoga5", "stickman",
		"stickman2",
	})
}

func loadLogo(m *Message) {
	if logo, ok := assets.GetLogo(m.IconURL); ok {
		m.IconURL = logo
	}
}

func Load(configFile string) (*Config, error) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath("/etc/rocketchat-jira-webhook")
		viper.AddConfigPath(".")
	}
	setDefaults()
	setDefaultsJira()
	setDefaultsConfluence()
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

	if config.Jira == nil && config.Confluence == nil {
		return nil, fmt.Errorf("define one of [confluence] or [jira] section is must")
	}

	if config.Jira != nil {
		if err := preProcJira(&config); err != nil {
			return nil, err
		}
	}

	if config.Confluence != nil {
		if err := preProcConfluence(&config); err != nil {
			return nil, err
		}
	}

	loadLogo(&config.MessageJira.Message)
	loadLogo(&config.MessageConfluence.Message)

	return &config, nil
}
