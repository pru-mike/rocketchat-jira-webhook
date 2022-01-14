package app

import (
	"net/http"

	"github.com/pru-mike/rocketchat-jira-webhook/jira"
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
)

func (app *App) GetJiraIssues(text string) ([]*jira.Issue, error) {
	if app.jira != nil {
		keys := app.validate.ValidateJiraKeys(app.jira.ParseKeys(text))
		if len(keys) == 0 {
			logger.Debug("jira keys not found")
			return nil, nil
		}
		logger.Debugf("found jira keys '%+v'", keys)

		issues, err := app.jira.GetIssues(keys)
		if err != nil {
			logger.Error(err)
		}
		return issues, err
	}
	return nil, nil
}

func (app *App) Jira(w http.ResponseWriter, req *http.Request) {
	in, err := app.readMessage(req.Body)
	if err != nil {
		return
	}

	if app.checkConnection(w, jiraConn, app.jiraErr) != nil {
		return
	}

	issues, err := app.GetJiraIssues(in.Text)

	if len(issues) == 0 {
		if err != nil && app.errToRocket {
			ok(w, app.jiraErr.Output(err.Error()))
		}
	} else {
		ok(w, app.jiraOut.Output(issues))
	}
}
