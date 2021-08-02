package rocketchat

import (
	"github.com/pru-mike/rocketchat-jira-webhook/assets"
	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/pru-mike/rocketchat-jira-webhook/jira"
	"golang.org/x/text/message"
	"html"
	"math/rand"
	"strconv"
	"unicode/utf8"
)

type Output struct {
	config      config.Config
	Alias       string       `json:"alias"`
	Avatar      string       `json:"avatar"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	AuthorName string  `json:"author_name"`
	Collapsed  bool    `json:"collapsed"`
	Title      string  `json:"title"`
	TitleLink  string  `json:"title_link"`
	Color      string  `json:"color"`
	Text       string  `json:"text"`
	Fields     []Field `json:"fields"`
}

func (attachment *Attachment) addShortField(title string, value string) {
	attachment.addField(title, value, true)
}

func (attachment *Attachment) addField(title string, value string, short bool) {
	attachment.Fields = append(attachment.Fields, Field{
		Title: title,
		Value: value,
		Short: short,
	})
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type OutputBuilder struct {
	*config.Message
	priorityToColor map[int]string
	printer         *message.Printer
}

func NewOutputBuilder(cfg *config.Message) *OutputBuilder {

	priorityToColor := make(map[int]string, len(cfg.PriorityIDPrecedence))
	var i int
	for _, priority := range cfg.PriorityIDPrecedence {
		if i > len(cfg.ColorsByPriority[i]) {
			i = 0
		}
		priorityToColor[priority] = cfg.ColorsByPriority[i]
		i++
	}

	return &OutputBuilder{
		cfg,
		priorityToColor,
		message.NewPrinter(cfg.LangTag()),
	}
}

func (b *OutputBuilder) NewMsg(text string) *Output {
	return &Output{
		Alias:  b.Username,
		Avatar: b.IconURL,
		Text:   text,
	}
}

func (b *OutputBuilder) New(issues []*jira.Issue) *Output {

	msg := b.NewMsg(b.printer.Sprintf("Found %d issue", len(issues)))

	for _, issue := range issues {
		attachment := Attachment{
			Collapsed: true,
			Title:     b.unescapeHTML(issue.Title()),
			TitleLink: issue.Link(),
			Text:      b.unescapeHTML(b.trim(issue.Description())),
			Color:     b.color(issue.Priority()),
		}
		b.addFields(issue, &attachment)
		b.addQuote(&attachment)
		msg.Attachments = append(msg.Attachments, attachment)
	}

	return msg
}

func (b *OutputBuilder) addQuote(attachment *Attachment) {
	if len(assets.Quotes) > 0 {
		if b.QuoteProbability > rand.Float32() {
			q := assets.Quotes[rand.Intn(len(assets.Quotes))]
			attachment.addField(q.Author, q.Quote, false)
		}
	}
}

func (b *OutputBuilder) getPersonName(p *jira.Person) string {
	if b.UseRealNames {
		return p.RealName()
	}
	return p.JiraName()
}

func (b *OutputBuilder) addFields(issue *jira.Issue, attachment *Attachment) {
	for _, field := range b.Fields {
		switch field {
		case config.Status:
			attachment.addShortField("Status", issue.Status())
		case config.Type:
			attachment.addShortField("Type", issue.IssueType())
		case config.Priority:
			attachment.addShortField("Priority", issue.Priority().Name)
		case config.Resolution:
			attachment.addShortField("Resolution", issue.Resolution())
		case config.Assignee:
			attachment.addShortField("Assignee", b.getPersonName(issue.Assignee()))
		case config.Reporter:
			attachment.addShortField("Reporter", b.getPersonName(issue.Reporter()))
		case config.Creator:
			attachment.addShortField("Creator", b.getPersonName(issue.Creator()))
		case config.Created:
			attachment.addShortField("Created", issue.Created().Format(b.DatetimeLayout))
		case config.Updated:
			attachment.addShortField("Updated", issue.Updated().Format(b.DatetimeLayout))
		case config.Watches:
			attachment.addShortField("Watches", strconv.Itoa(issue.Watches()))
		case config.Components:
			attachment.addShortField("Components", issue.Components())
		case config.Labels:
			attachment.addShortField("Labels", issue.Labels())
		}
	}
}

func (b *OutputBuilder) trim(text string) string {
	if b.MaxTextLen > 0 && utf8.RuneCountInString(text) > b.MaxTextLen {
		i := 0
		for r := range text {
			if i == b.MaxTextLen {
				return text[:r] + "\u2026"
			}
			i++
		}
	}
	return text
}

func (b *OutputBuilder) color(priority *jira.Priority) (color string) {
	if !b.PriorityColors {
		return b.DefaultColor
	}
	if color, ok := b.priorityToColor[priority.GetID()]; ok {
		return color
	}
	return b.DefaultColor
}

func (b *OutputBuilder) unescapeHTML(text string) string {
	if b.UnescapeHTML {
		return html.UnescapeString(text)
	}
	return text
}
