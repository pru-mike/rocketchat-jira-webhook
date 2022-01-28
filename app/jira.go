package app

import (
	"net/http"

	"github.com/pru-mike/rocketchat-jira-webhook/jira"
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
)

func (app *App) GetJiraIssues(text string) ([]*jira.Issue, error) {
	if app.jira == nil {
		return nil, nil
	}

	keys := app.validate.ValidateJiraKeys(app.jira.ParseKeys(text))
	if len(keys) == 0 {
		logger.Debug("jira keys not found")
		return nil, nil
	}
	logger.Debugf("found jira keys '%+v'", keys)

	issues, err := app.jira.GetIssues(keys)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return issues, nil
}

func (app *App) Jira(w http.ResponseWriter, req *http.Request) {
	in, err := app.readMessage(req.Body)
	if err != nil {
		return
	}

	if err := app.checkConnection(w, jiraConn, app.jiraErr); err != nil {
		return
	}

	issues, err := app.GetJiraIssues(in.TextWithoutReply())

	if len(issues) == 0 {
		if err != nil && app.errToRocket {
			ok(w, app.jiraErr.Output(err.Error()))
		}
		return
	}

	ok(w, app.jiraOut.Output(issues))
}
