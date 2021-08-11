package jira

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Time time.Time

func (t *Time) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	pt, err := time.Parse("2006-01-02T15:04:05.999-0700", str)
	if err != nil {
		return err
	}
	*t = Time(pt)
	return nil
}

type FieldsID int

func (id *FieldsID) UnmarshalJSON(b []byte) error {
	var strID string
	err := json.Unmarshal(b, &strID)
	if err != nil {
		return err
	}
	intID, err := strconv.ParseInt(strID, 10, 32)
	if err != nil {
		return err
	}
	*id = FieldsID(int(intID))
	return nil
}

type Priority struct {
	ID      FieldsID `json:"id"`
	Name    string   `json:"name"`
	Self    string   `json:"self"`
	IconURL string   `json:"IconURL"`
}

func (p *Priority) GetID() int {
	return int(p.ID)
}

type Person struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
}

func (p *Person) JiraName() string {
	return p.Name
}

func (p *Person) RealName() string {
	return p.DisplayName
}

type Issue struct {
	browseURL string
	ID        FieldsID `json:"id"`
	Self      string   `json:"self"`
	Key       string   `json:"key"`
	Fields    struct {
		Summary     string   `json:"summary"`
		Description string   `json:"description"`
		Created     Time     `json:"created"`
		Updated     Time     `json:"updated"`
		Priority    Priority `json:"priority"`
		Status      struct {
			ID   FieldsID `json:"id"`
			Name string   `json:"name"`
		} `json:"status"`
		IssueType struct {
			ID   FieldsID `json:"id"`
			Name string   `json:"name"`
		} `json:"issuetype"`
		Resolution struct {
			ID   FieldsID `json:"id"`
			Name string   `json:"name"`
		} `json:"resolution"`
		Assignee Person `json:"assignee"`
		Reporter Person `json:"reporter"`
		Creator  Person `json:"creator"`
		Watches  struct {
			WatchCount int `json:"watchCount"`
		} `json:"watches"`
		Components []struct {
			ID   FieldsID `json:"id"`
			Name string   `json:"name"`
		} `json:"components"`
		Labels []string `json:"labels"`
	} `json:"fields"`
}

func (i *Issue) Link() string {
	return i.browseURL + "/" + i.Key
}

func (i *Issue) DefaultTitle() string {
	return i.Key + " " + i.Fields.Summary
}

func (i *Issue) GetKey() string {
	return i.Key
}

func (i *Issue) GetSummary() string {
	return i.Fields.Summary
}

func (i *Issue) Description() string {
	return i.Fields.Description
}

func (i *Issue) Priority() *Priority {
	return &i.Fields.Priority
}

func (i *Issue) Status() string {
	return i.Fields.Status.Name
}

func (i *Issue) IssueType() string {
	return i.Fields.IssueType.Name
}

func (i *Issue) Resolution() string {
	return i.Fields.Resolution.Name
}

func (i *Issue) Assignee() *Person {
	return &i.Fields.Assignee
}

func (i *Issue) Reporter() *Person {
	return &i.Fields.Reporter
}

func (i *Issue) Creator() *Person {
	return &i.Fields.Creator
}

func (i *Issue) Created() time.Time {
	return time.Time(i.Fields.Created)
}

func (i *Issue) Updated() time.Time {
	return time.Time(i.Fields.Updated)
}

func (i *Issue) Watches() int {
	return i.Fields.Watches.WatchCount
}

func (i *Issue) Components() string {
	components := make([]string, len(i.Fields.Components))
	for i, c := range i.Fields.Components {
		components[i] = c.Name
	}
	return strings.Join(components, ",")
}

func (i *Issue) Labels() string {
	return strings.Join(i.Fields.Labels, ",")
}
