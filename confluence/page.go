package confluence

import (
	"encoding/json"
	"strconv"
	"time"
)

type Time time.Time

func (t *Time) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	pt, err := time.Parse("2006-01-02T15:04:05.999-07:00", str)
	if err != nil {
		return err
	}
	*t = Time(pt)
	return nil
}

type Person struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

func (p *Person) JiraName() string {
	return p.Name
}

func (p *Person) RealName() string {
	return p.DisplayName
}

type BodyValue struct {
	Value string `json:"value"`
}

type Page struct {
	ID     string `json:"id"`
	Type   string `json:"types"`
	Status string `json:"status"`
	Title  string `json:"title"`
	Space  struct {
		ID   int    `json:"id"`
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"space"`
	Links struct {
		WebUI  string `json:"webui"`
		TinyUI string `json:"tinyui"`
		Base   string `json:"base"`
	} `json:"_links"`
	History struct {
		Latest      bool `json:"latest"`
		CreatedDate Time `json:"createdDate"`
		CreatedBy   struct {
			Username    string `json:"username"`
			DisplayName string `json:"displayName"`
		} `json:"createdBy"`
		LastUpdated struct {
			By struct {
				Username    string `json:"username"`
				DisplayName string `json:"displayName"`
			} `json:"by"`
			When   Time `json:"when"`
			Number int  `json:"number"`
		} `json:"lastUpdated"`
	} `json:"history"`
	Body struct {
		Editor              *BodyValue `json:"editor,omitempty"`
		View                *BodyValue `json:"view,omitempty"`
		ExportView          *BodyValue `json:"export_view,omitempty"`
		Storage             *BodyValue `json:"storage,omitempty"`
		AnonymousExportView *BodyValue `json:"anonymous_export_view,omitempty"`
		StyledView          *BodyValue `json:"styled_view,omitempty"`
	} `json:"body"`
}

func (p *Page) Description() string {
	if p.Body.Editor != nil {
		return p.Body.Editor.Value
	} else if p.Body.View != nil {
		return p.Body.View.Value
	} else if p.Body.ExportView != nil {
		return p.Body.ExportView.Value
	} else if p.Body.Storage != nil {
		return p.Body.Storage.Value
	} else if p.Body.AnonymousExportView != nil {
		return p.Body.AnonymousExportView.Value
	} else if p.Body.StyledView != nil {
		return p.Body.StyledView.Value
	}
	return ""
}

func (p *Page) SpaceName() string {
	return p.Space.Name
}

func (p *Page) CreatedBy() *Person {
	return &Person{
		Name:        p.History.CreatedBy.Username,
		DisplayName: p.History.CreatedBy.DisplayName,
	}
}

func (p *Page) CreatedDate() time.Time {
	return time.Time(p.History.CreatedDate)
}

func (p *Page) UpdatedBy() *Person {
	return &Person{
		Name:        p.History.LastUpdated.By.Username,
		DisplayName: p.History.LastUpdated.By.DisplayName,
	}
}

func (p *Page) UpdatedDate() time.Time {
	return time.Time(p.History.LastUpdated.When)
}

func (p *Page) Link() string {
	return p.Links.Base + p.Links.TinyUI
}

func (p *Page) GetStatus() string {
	return p.Status
}

func (p *Page) Latest() string {
	return strconv.FormatBool(p.History.Latest)
}

func (p *Page) LastVersionNumber() string {
	return strconv.Itoa(p.History.LastUpdated.Number)
}
