package app

import (
	"github.com/pru-mike/rocketchat-jira-webhook/confluence"
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
	"net/http"
)

func (app *App) GetConfluencePages(text string) ([]*confluence.Page, error) {
	if app.confluence != nil {
		ids := app.confluence.FindPagesIDs(text)
		if len(ids) == 0 {
			logger.Debug("confluence ids not found")
			return nil, nil
		}
		logger.Debugf("found confluence ids '%+v'", ids)

		pages, err := app.confluence.GetPages(ids)
		if err != nil {
			logger.Error(err)
		}
		return pages, err
	}
	return nil, nil
}

func (app *App) Confluence(w http.ResponseWriter, req *http.Request) {

	in, err := app.readMessage(req.Body)
	if err != nil {
		return
	}

	if app.checkConnection(w, confluenceConn, app.confluenceErr) != nil {
		return
	}

	pages, err := app.GetConfluencePages(in.Text)

	if len(pages) == 0 {
		if err != nil && app.errToRocket {
			ok(w, app.confluenceErr.Output(err.Error()))
		}
	} else {
		ok(w, app.confluenceOut.Output(pages))
	}
}
