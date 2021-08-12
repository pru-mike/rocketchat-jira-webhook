package app

import (
	"net/http"
)

func (app *App) JiraConfluence(w http.ResponseWriter, req *http.Request) {

	in, err := app.readMessage(req.Body)
	if err != nil {
		return
	}

	if app.checkConnection(w, confluenceConn, app.confluenceErr) != nil {
		return
	}

	if app.checkConnection(w, jiraConn, app.jiraErr) != nil {
		return
	}

	issues, errJira := app.GetJiraIssues(in.Text)
	pages, errConfluence := app.GetConfluencePages(in.Text)

	if len(pages) == 0 && len(issues) == 0 {
		if app.errToRocket {
			var sendErr string
			if errJira != nil {
				sendErr += errJira.Error()
			}
			if errConfluence != nil {
				sendErr += errConfluence.Error()
			}
			ok(w, app.jiraErr.Output(sendErr))
		}
	} else {
		ok(w,
			app.multiplexOut.Output(
				app.jiraOut.New(issues),
				app.confluenceOut.New(pages),
			),
		)
	}

}
