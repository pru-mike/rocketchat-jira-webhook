package rocketchat

import (
	"bytes"
	"html"
	"text/template"
	"unicode/utf8"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/pru-mike/rocketchat-jira-webhook/assets"
	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
	"golang.org/x/text/message"
)

type OutputBuilder struct {
	*config.Message
	printer    *message.Printer
	titleTmpl  *template.Template
	authorTmpl *template.Template
}

func NewOutputBuilder(cfg *config.Message) *OutputBuilder {
	return &OutputBuilder{
		cfg,
		message.NewPrinter(cfg.LangTag()),
		template.Must(template.New("titleTmpl").Parse(cfg.TitleTemplate)),
		template.Must(template.New("authorTmpl").Parse(cfg.AuthorTemplate)),
	}
}

func (o *OutputBuilder) Username() string {
	return o.Message.Username
}

func (o *OutputBuilder) IconURL() string {
	return o.Message.IconURL
}

func (o *OutputBuilder) Unescape(text string) string {
	if o.Message.UnescapeHTML {
		return html.UnescapeString(text)
	}
	return text
}

func (o *OutputBuilder) StripTags(text string) string {
	if o.Message.StripTags {
		return strip.StripTags(text)
	}
	return text
}

func (o *OutputBuilder) BuildMessage(key string, data ...interface{}) string {
	return o.printer.Sprintf(key, data...)
}

func (o *OutputBuilder) BuildTitle(data interface{}) string {
	var title bytes.Buffer
	err := o.titleTmpl.Execute(&title, data)
	if err != nil {
		logger.Errorf("can't execute title_template %v", err)
		return ""
	}
	return o.Unescape(title.String())
}

func (o *OutputBuilder) BuildAuthor(data interface{}) string {
	var author bytes.Buffer
	err := o.authorTmpl.Execute(&author, data)
	if err != nil {
		logger.Errorf("can't execute author_template %v", err)
		return ""
	}
	return author.String()
}

func (o *OutputBuilder) TrimMaxLen(text string) string {
	if o.MaxTextLen > 0 && utf8.RuneCountInString(text) > o.MaxTextLen {
		i := 0
		for r := range text {
			if i == o.MaxTextLen {
				return text[:r] + "\u2026"
			}
			i++
		}
	}
	return text
}

func (o *OutputBuilder) AddQuote(attachment *Attachment) {
	if ok, q := assets.GetQuoteWithProb(o.QuoteProbability); ok {
		attachment.AddField(q.Author, q.Quote, false)
	}
}

type Person interface {
	RealName() string
	JiraName() string
}

func (o *OutputBuilder) GetPersonName(p Person) string {
	if o.UseRealNames {
		return p.RealName()
	}
	return p.JiraName()
}
