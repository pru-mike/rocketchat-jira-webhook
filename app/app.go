package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/pru-mike/rocketchat-jira-webhook/confluence"
	"github.com/pru-mike/rocketchat-jira-webhook/jira"
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
	"github.com/pru-mike/rocketchat-jira-webhook/rocketchat"
)

type App struct {
	errToRocket   bool
	validate      *rocketchat.Validate
	jira          *jira.Jira
	jiraOut       *rocketchat.JiraOutputBuilder
	jiraErr       *rocketchat.TextOutputBuilder
	confluence    *confluence.Confluence
	confluenceOut *rocketchat.ConfluenceOutputBuilder
	confluenceErr *rocketchat.TextOutputBuilder
	multiplexOut  *rocketchat.MultiplexOutputBuilder
}

func New(cfg *config.Config) *App {
	return &App{
		errToRocket:   cfg.App.ErrToRocket,
		validate:      rocketchat.SetupValidator(&cfg.Rocketchat),
		jira:          jiraClient(cfg.Jira),
		jiraOut:       rocketchat.NewJiraOutputBuilder(&cfg.MessageJira),
		jiraErr:       rocketchat.NewTextOutputBuilder(&cfg.MessageJira.Message),
		confluence:    confluenceClient(cfg.Confluence),
		confluenceOut: rocketchat.NewConfluenceOutputBuilder(&cfg.MessageConfluence),
		confluenceErr: rocketchat.NewTextOutputBuilder(&cfg.MessageConfluence.Message),
		multiplexOut:  rocketchat.NewMultiplexOutputBuilder(),
	}
}

func jiraClient(cfg *config.Jira) *jira.Jira {
	if cfg != nil {
		return jira.NewClient(cfg)
	}
	return nil
}

func confluenceClient(cfg *config.Confluence) *confluence.Confluence {
	if cfg != nil {
		return confluence.NewClient(cfg)
	}
	return nil
}

func (app *App) readMessage(r io.Reader) (*rocketchat.Input, error) {
	var in rocketchat.Input
	err := json.NewDecoder(r).Decode(&in)
	if err != nil {
		logger.Errorf("can't decode message: %v", err)
		return nil, err
	}

	logger.Debugf("input message: %+v", in)
	err = app.validate.Struct(in)
	if err != nil {
		logger.Debugf("validation failed %v", err)
		return nil, err
	}
	return &in, nil
}

func ok(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error(err)
	}
}

func err(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error(err)
	}
}

type conn string

const (
	jiraConn       conn = "jira"
	confluenceConn conn = "confluence"
)

func (app *App) checkConnection(w http.ResponseWriter, c conn, out *rocketchat.TextOutputBuilder) error {
	var isSet bool
	switch c {
	case jiraConn:
		isSet = app.jira != nil
	case confluenceConn:
		isSet = app.confluence != nil
	default:
		panic(fmt.Sprintf("wrond connection type: ;%s'", c))
	}
	if !isSet {
		err := fmt.Errorf("rocketchat-jira-webhook configuration error: '%s' connections is turned off", c)
		logger.Error(err.Error())
		if app.errToRocket {
			ok(w, out.Output(err.Error()))
		}
		return err
	}
	return nil
}
