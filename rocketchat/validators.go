package rocketchat

import (
	"github.com/go-playground/validator/v10"
	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/pru-mike/rocketchat-jira-webhook/jira"
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
)

type Validate struct {
	*validator.Validate
	blacklistedJiraProjectKeys []string
	whitelistedJiraProjectKeys []string
}

func SetupValidator(cfg *config.Rocketchat) *Validate {
	var err error
	validate := validator.New()

	if err = validate.RegisterValidation("token", ValidateToken(cfg)); err != nil {
		logger.Fatalf("ValidateToken registration failed: %v", err)
	}
	if err = validate.RegisterValidation("edits", ValidateEdits(cfg)); err != nil {
		logger.Fatalf("ValidateEdits registration failed: %v", err)
	}
	if err = validate.RegisterValidation("bots", ValidateBots(cfg)); err != nil {
		logger.Fatalf("ValidateBots registration failed: %v", err)
	}
	validate.RegisterStructValidation(ValidateInputMessageStructLevel(cfg), Input{})

	return &Validate{
		validate,
		cfg.BlacklistedJiraProjectKeys,
		cfg.WhitelistedJiraProjectKeys,
	}
}

func (v *Validate) ValidateJiraKeys(keys []string) []string {
	if len(v.blacklistedJiraProjectKeys) > 0 {
		var i int
	BlackListLoop:
		for _, key := range keys {
			prjKey := jira.StripKey(key)
			for _, forbiddenKey := range v.blacklistedJiraProjectKeys {
				if prjKey == forbiddenKey {
					logger.Debugf("forbidden key with blacklist key %s", key)
					continue BlackListLoop
				}
			}
			keys[i] = key
			i++
		}
		keys = keys[:i]
	}
	if len(v.whitelistedJiraProjectKeys) > 0 {
		var i int
		for _, key := range keys {
			prjKey := jira.StripKey(key)
			var allowed bool
			for _, allowedKey := range v.whitelistedJiraProjectKeys {
				if prjKey == allowedKey {
					allowed = true
					break
				}
			}
			if !allowed {
				logger.Debugf("forbidden key with whitelist key %s", key)
				continue
			}
			keys[i] = key
			i++
		}
		keys = keys[:i]
	}

	return keys
}

func AlwaysValid(validator.FieldLevel) bool {
	return true
}

func ValidateToken(config *config.Rocketchat) func(validator.FieldLevel) bool {
	if len(config.Tokens) == 0 {
		return AlwaysValid
	}
	return func(fl validator.FieldLevel) bool {
		reqToken := fl.Field().String()
		for _, confToken := range config.Tokens {
			if reqToken == confToken {
				return true
			}
		}
		return false
	}
}

func ValidateEdits(config *config.Rocketchat) func(validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		isEdited := fl.Field().Bool()
		return !isEdited || config.AllowEdits
	}
}

func ValidateBots(config *config.Rocketchat) func(validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		isBot := fl.Field().Bool()
		return !isBot || config.AllowBots
	}
}

func ValidateInputMessageStructLevel(config *config.Rocketchat) func(validator.StructLevel) {
	return func(sl validator.StructLevel) {
		ValidateUserStructLevel(config)(sl)
		ValidateChannelStructLevel(config)(sl)
	}
}

func ValidateUserStructLevel(config *config.Rocketchat) func(validator.StructLevel) {
	return validateWhiteListStructLevel(
		"UserID", "UserName",
		config.WhitelistedUsers, config.BlacklistedUsers,
	)
}

func ValidateChannelStructLevel(config *config.Rocketchat) func(validator.StructLevel) {
	return validateWhiteListStructLevel(
		"ChannelID", "ChannelName",
		config.WhitelistedChannels, config.BlacklistedChannels,
	)
}

func validateWhiteListStructLevel(id, name string, whitelist []string, blacklist []string) func(validator.StructLevel) {
	if len(whitelist) == 0 && len(blacklist) == 0 {
		return func(validator.StructLevel) {}
	}
	return func(sl validator.StructLevel) {
		ID := sl.Current().FieldByName(id).String()
		Name := sl.Current().FieldByName(name).String()

		for _, forbidden := range blacklist {
			if ID == forbidden {
				sl.ReportError(ID, id, id, "blacklisted", "")
				return
			}
			if Name == forbidden {
				sl.ReportError(Name, name, name, "blacklisted", "")
				return
			}
		}
		if len(whitelist) == 0 {
			return
		}
		for _, allowed := range whitelist {
			if ID == allowed || Name == allowed {
				return
			}
		}
		sl.ReportError(Name, name, name, "whitelisted", "")
	}
}
