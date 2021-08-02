package jira

type Myself struct {
	Self         string
	Key          string
	Name         string
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool
	Deleted      bool
	TimeZone     string `json:"timeZone"`
	Locale       string
}
