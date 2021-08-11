package app

import (
	"net/http"
	"reflect"
)

type connHealth struct {
	Name string `json:"name"`
	Err  string `json:"error"`
}

type health struct {
	Ok         bool       `json:"ok"`
	Jira       connHealth `json:"jira"`
	Confluence connHealth `json:"confluence"`
}

func newHealthResponse() *health {
	return &health{
		Ok: true,
		Jira: connHealth{
			Name: "",
			Err:  "jira not configured",
		},
		Confluence: connHealth{
			Name: "",
			Err:  "confluence not configured",
		},
	}
}

type healthChecker interface {
	GetCurrentUser() (string, error)
}

func (h *health) checkStatus(check healthChecker, conn *connHealth) {
	if check != nil && !reflect.ValueOf(check).IsNil() {
		var err error
		conn.Name, err = check.GetCurrentUser()
		if err != nil {
			h.Ok = false
			conn.Err = err.Error()
		} else {
			conn.Err = ""
		}
	}
}

func (app *App) Health(w http.ResponseWriter, _ *http.Request) {
	response := newHealthResponse()
	response.checkStatus(app.jira, &response.Jira)
	response.checkStatus(app.confluence, &response.Confluence)

	if !response.Ok {
		err(w, response)
	} else {
		ok(w, response)
	}
}
