package rocketchat

import (
	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/pru-mike/rocketchat-jira-webhook/confluence"
	"github.com/pru-mike/rocketchat-jira-webhook/utils"
)

type ConfluenceOutputBuilder struct {
	*OutputBuilder
	cfg *config.MessageConfluence
}

func NewConfluenceOutputBuilder(cfg *config.MessageConfluence) *ConfluenceOutputBuilder {
	return &ConfluenceOutputBuilder{
		NewOutputBuilder(&cfg.Message),
		cfg,
	}
}

var _ OutputInterface = (*ConfluenceOutput)(nil)

type ConfluenceOutput struct {
	*ConfluenceOutputBuilder
	pages []*confluence.Page
}

func (o *ConfluenceOutput) GetText() string {
	return o.BuildMessage("Found %d doc", len(o.pages))
}

func (o *ConfluenceOutput) AttachmentsLen() int {
	return len(o.pages)
}

func (o *ConfluenceOutput) Attachments() []Attachment {
	nextIcon := o.NextIconGetter()
	attachments := make([]Attachment, len(o.pages))
	for i, page := range o.pages {
		attachment := Attachment{
			Collapsed:  true,
			Title:      o.BuildTitle(page),
			TitleLink:  page.Link(),
			AuthorName: o.BuildAuthor(page),
			AuthorIcon: nextIcon(),
			Text:       o.TrimMaxLen(o.StripTags(o.Unescape(page.Description()))),
			Color:      o.Color(),
		}
		o.AddFields(page, &attachment)
		o.AddQuote(&attachment)
		attachments[i] = attachment
	}

	return attachments
}

func (o *ConfluenceOutputBuilder) New(pages []*confluence.Page) *ConfluenceOutput {
	return &ConfluenceOutput{
		ConfluenceOutputBuilder: o,
		pages:                   pages,
	}
}

func (o *ConfluenceOutputBuilder) Output(pages []*confluence.Page) *Output {
	return New(o.New(pages))
}

func (o *ConfluenceOutput) NextIconGetter() func() string {
	if len(o.AuthorIcons) == 0 {
		return func() string {
			return ""
		}
	}

	authorIconsGetter := utils.NextIconGetter(o.AuthorIcons)
	return func() string {
		return authorIconsGetter()
	}
}

func (o *ConfluenceOutput) Color() string {
	return o.DefaultColor
}

func (o *ConfluenceOutputBuilder) AddFields(page *confluence.Page, attachment *Attachment) {
	for _, field := range o.cfg.Fields {
		switch field {
		case config.SpaceName:
			attachment.AddField("Space", page.SpaceName(), false)
		case config.CreatedBy:
			attachment.AddShortField("CreatedBy", o.GetPersonName(page.CreatedBy()))
		case config.CreatedDate:
			attachment.AddShortField("CreatedDate", page.CreatedDate().Format(o.DatetimeLayout))
		case config.UpdatedBy:
			attachment.AddShortField("UpdatedBy", o.GetPersonName(page.UpdatedBy()))
		case config.UpdatedDate:
			attachment.AddShortField("UpdatedDate", page.UpdatedDate().Format(o.DatetimeLayout))
		case config.LastVersionNumber:
			attachment.AddShortField("Last Version", page.LastVersionNumber())
		case config.ConfluenceStatus:
			attachment.AddShortField("Status", page.GetStatus())
		case config.Latest:
			attachment.AddShortField("Latest", page.Latest())
		}
	}
}
