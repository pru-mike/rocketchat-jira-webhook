package app

import (
	"encoding/json"
	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/pru-mike/rocketchat-jira-webhook/jira"
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
	"github.com/pru-mike/rocketchat-jira-webhook/rocketchat"
	"net/http"
)

type App struct {
	errToRocket bool
	validate    *rocketchat.Validate
	jira        *jira.Jira
	out         *rocketchat.OutputBuilder
}

func New(cfg *config.Config) *App {
	return &App{
		errToRocket: cfg.App.ErrToRocket,
		validate:    rocketchat.SetupValidator(&cfg.Rocketchat),
		jira:        jira.NewClient(&cfg.Jira),
		out:         rocketchat.NewOutputBuilder(&cfg.Message),
	}
}

func (app *App) Jira(w http.ResponseWriter, req *http.Request) {

	var in rocketchat.Input
	err := json.NewDecoder(req.Body).Decode(&in)
	if err != nil {
		logger.Errorf("can't decode message: %v", err)
		return
	}
	logger.Debugf("input message: %+v", in)

	err = app.validate.Struct(in)
	if err != nil {
		logger.Debugf("validation failed %v", err)
		return
	}

	keys := app.validate.ValidateKeys(jira.ParseKeys(in.Text))
	if len(keys) == 0 {
		logger.Debug("jira keys not found")
		return
	}
	logger.Debugf("found jira keys '%+v'", keys)

	issues, err := app.jira.GetIssues(keys)
	if err != nil {
		logger.Error(err)
	}
	if len(issues) == 0 {
		if err != nil && app.errToRocket {
			writeResponse(w, app.out.NewMsg(err.Error()))
		}
	} else {
		writeResponse(w, app.out.New(issues))
	}
}

func (app *App) Health(w http.ResponseWriter, _ *http.Request) {
	me, err := app.jira.GetMyself()
	ok := true
	var errStr string
	if err != nil {
		ok = false
		errStr = err.Error()
		logger.Errorf("can't get myself: %v", err)
	}
	health := struct {
		Ok   bool `json:"ok"`
		Jira struct {
			Name  string `json:"name"`
			Error string `json:"error"`
		} `json:"jira"`
	}{
		Ok: ok,
		Jira: struct {
			Name  string `json:"name"`
			Error string `json:"error"`
		}{
			Name:  me.DisplayName,
			Error: errStr,
		},
	}
	if !ok {
		writeError(w, health)
	} else {
		writeResponse(w, health)
	}
}

func writeResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Error(err)
	}
}

func writeError(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	writeResponse(w, data)
}
