package rocketchat

import (
	"bytes"
	"github.com/pru-mike/rocketchat-jira-webhook/assets"
	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/pru-mike/rocketchat-jira-webhook/jira"
	"github.com/pru-mike/rocketchat-jira-webhook/logger"
	"golang.org/x/text/message"
	"html"
	"math/rand"
	"sort"
	"strconv"
	"text/template"
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
	AuthorName  string  `json:"author_name"`
	AuthorLink  string  `json:"author_link,omitempty"`
	AuthorIcon  string  `json:"author_icon,omitempty"`
	Collapsed   bool    `json:"collapsed"`
	Title       string  `json:"title"`
	TitleLink   string  `json:"title_link,omitempty"`
	MessageLink string  `json:"message_link,omitempty"`
	ImageURL    string  `json:"image_url,omitempty"`
	ThumbURL    string  `json:"thumb_url,omitempty"`
	Color       string  `json:"color"`
	Text        string  `json:"text"`
	Fields      []Field `json:"fields"`
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
	priorityToColor      map[int]string
	priorityToPrecedence map[int]int
	printer              *message.Printer
	titleTmpl            *template.Template
	authorTmpl           *template.Template
}

func NewOutputBuilder(cfg *config.Message) *OutputBuilder {

	priorityToColor := make(map[int]string, len(cfg.PriorityIDPrecedence))
	priorityToPrecedence := make(map[int]int, len(cfg.PriorityIDPrecedence))
	var i int
	for p, priority := range cfg.PriorityIDPrecedence {
		priorityToPrecedence[priority] = p
		if i > len(cfg.ColorsByPriority) {
			i = 0
		}
		if len(cfg.ColorsByPriority) > 0 {
			priorityToColor[priority] = cfg.ColorsByPriority[i]
		}
		i++
	}

	return &OutputBuilder{
		cfg,
		priorityToColor,
		priorityToPrecedence,
		message.NewPrinter(cfg.LangTag()),
		template.Must(template.New("titleTmpl").Parse(cfg.TitleTemplate)),
		template.Must(template.New("authorTmpl").Parse(cfg.AuthorTemplate)),
	}
}

func (b *OutputBuilder) NewMsg(text string) *Output {
	return &Output{
		Alias:  b.Username,
		Avatar: b.IconURL,
		Text:   text,
	}
}

func (b *OutputBuilder) getTitle(issue *jira.Issue) string {
	var title bytes.Buffer
	err := b.titleTmpl.Execute(&title, issue)
	if err != nil {
		logger.Errorf("can't execute title_template %v", err)
		return b.unescapeHTML(issue.DefaultTitle())
	}
	return b.unescapeHTML(title.String())
}

func (b *OutputBuilder) getAuthor(issue *jira.Issue) string {
	var author bytes.Buffer
	err := b.authorTmpl.Execute(&author, issue)
	if err != nil {
		logger.Errorf("can't execute author_template %v", err)
		return ""
	}
	return author.String()
}

func getNextElem(src []string, n uint) string {
	if len(src) > 0 {
		return src[n%uint(len(src))]
	}
	return ""
}

func getNextLogo(logos []string, n int) string {
	logo, _ := assets.GetLogo(getNextElem(logos, uint(n)))
	return logo
}

func (b *OutputBuilder) makeNextAuthorIconGetter() func(issue *jira.Issue) string {
	if len(b.AuthorIcons) == 0 {
		return func(issue *jira.Issue) string {
			return ""
		}
	}
	authorIconsGetter := b.makeNextIconGetter(b.AuthorIcons)

	if len(b.InactiveAuthorIcons) == 0 {
		return func(_ *jira.Issue) string {
			return authorIconsGetter()
		}
	}

	inactiveAuthorIconGetter := b.makeNextIconGetter(b.InactiveAuthorIcons)
	inactiveAuthor := b.InactiveAuthor
	return func(issue *jira.Issue) string {
		var isActive bool
		switch inactiveAuthor {
		case config.Assignee:
			isActive = issue.Assignee().Active
		case config.Creator:
			isActive = issue.Creator().Active
		default:
			isActive = issue.Reporter().Active
		}
		if isActive {
			return authorIconsGetter()
		}
		return inactiveAuthorIconGetter()
	}
}

func (b *OutputBuilder) makeNextIconGetter(icons []string) func() string {
	if len(icons) == 1 {
		return func() string {
			return getNextLogo(icons, 0)
		}
	}
	iconsCopy := make([]string, len(icons))
	copy(iconsCopy, icons)
	rand.Shuffle(len(iconsCopy), func(i, j int) {
		iconsCopy[i], iconsCopy[j] = iconsCopy[j], iconsCopy[i]
	})
	i := rand.Intn(len(iconsCopy))
	return func() string {
		i++
		return getNextLogo(iconsCopy, i)
	}
}

func (b *OutputBuilder) New(issues []*jira.Issue) *Output {

	msg := b.NewMsg(b.printer.Sprintf("Found %d issue", len(issues)))

	if b.SortByPrecedence {
		sort.Slice(issues, func(i, j int) bool {
			pri1 := issues[j].Priority().GetID()
			pri2 := issues[i].Priority().GetID()
			pre1, ok1 := b.priorityToPrecedence[pri1]
			pre2, ok2 := b.priorityToPrecedence[pri2]
			if ok1 && ok2 {
				return pre1 > pre2
			}
			return pri1 > pri2
		})
	}

	getNextAuthorIcon := b.makeNextAuthorIconGetter()
	for _, issue := range issues {
		attachment := Attachment{
			Collapsed:  true,
			Title:      b.getTitle(issue),
			TitleLink:  issue.Link(),
			AuthorName: b.getAuthor(issue),
			AuthorIcon: getNextAuthorIcon(issue),
			Text:       b.unescapeHTML(b.trim(issue.Description())),
			Color:      b.color(issue.Priority()),
		}
		b.addFields(issue, &attachment)
		b.addQuote(&attachment)
		msg.Attachments = append(msg.Attachments, attachment)
	}

	return msg
}

func (b *OutputBuilder) addQuote(attachment *Attachment) {
	if ok, q := assets.GetQuoteWithProb(b.QuoteProbability); ok {
		attachment.addField(q.Author, q.Quote, false)
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
