package rocketchat

import (
	"github.com/pru-mike/rocketchat-jira-webhook/config"
	"github.com/pru-mike/rocketchat-jira-webhook/jira"
	"github.com/pru-mike/rocketchat-jira-webhook/utils"
	"sort"
	"strconv"
)

type JiraOutputBuilder struct {
	*OutputBuilder
	cfg                  *config.MessageJira
	priorityToColor      map[int]string
	priorityToPrecedence map[int]int
}

func NewJiraOutputBuilder(cfg *config.MessageJira) *JiraOutputBuilder {

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

	return &JiraOutputBuilder{
		NewOutputBuilder(&cfg.Message),
		cfg,
		priorityToColor,
		priorityToPrecedence,
	}
}

var _ OutputInterface = (*JiraOutput)(nil)

type JiraOutput struct {
	*JiraOutputBuilder
	issues []*jira.Issue
}

func (o *JiraOutput) GetText() string {
	return o.BuildMessage("Found %d issue", len(o.issues))
}

func (o *JiraOutput) AttachmentsLen() int {
	return len(o.issues)
}

func (o *JiraOutput) Attachments() []Attachment {
	nextIcon := o.NextIconGetter()
	attachments := make([]Attachment, len(o.issues))
	for i, issue := range o.issues {
		attachment := Attachment{
			Collapsed:  true,
			Title:      o.BuildTitle(issue),
			TitleLink:  issue.Link(),
			AuthorName: o.BuildAuthor(issue),
			AuthorIcon: nextIcon(issue),
			Text:       o.TrimMaxLen(o.StripTags(o.Unescape(issue.Description()))),
			Color:      o.Color(issue.Priority()),
		}
		o.AddFields(issue, &attachment)
		o.AddQuote(&attachment)
		attachments[i] = attachment
	}

	return attachments
}

func (o *JiraOutputBuilder) New(issues []*jira.Issue) *JiraOutput {
	if o.cfg.SortByPrecedence && len(issues) > 1 {
		sort.Slice(issues, func(i, j int) bool {
			pri1 := issues[j].Priority().GetID()
			pri2 := issues[i].Priority().GetID()
			pre1, ok1 := o.priorityToPrecedence[pri1]
			pre2, ok2 := o.priorityToPrecedence[pri2]
			if ok1 && ok2 {
				return pre1 > pre2
			}
			return pri1 > pri2
		})
	}
	return &JiraOutput{
		JiraOutputBuilder: o,
		issues:            issues,
	}
}

func (o *JiraOutputBuilder) Output(issues []*jira.Issue) *Output {
	return New(o.New(issues))
}

func (o *JiraOutputBuilder) NextIconGetter() func(issue *jira.Issue) string {
	if len(o.AuthorIcons) == 0 {
		return func(issue *jira.Issue) string {
			return ""
		}
	}

	authorIconsGetter := utils.NextIconGetter(o.AuthorIcons)

	if len(o.cfg.InactiveAuthorIcons) == 0 {
		return func(_ *jira.Issue) string {
			return authorIconsGetter()
		}
	}

	inactiveAuthorIconGetter := utils.NextIconGetter(o.cfg.InactiveAuthorIcons)
	inactiveAuthor := o.cfg.InactiveAuthor
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

func (o *JiraOutputBuilder) Color(priority *jira.Priority) (color string) {
	if !o.cfg.PriorityColors {
		return o.DefaultColor
	}
	if color, ok := o.priorityToColor[priority.GetID()]; ok {
		return color
	}
	return o.DefaultColor
}

func (o *JiraOutputBuilder) AddFields(issue *jira.Issue, attachment *Attachment) {
	for _, field := range o.cfg.Fields {
		switch field {
		case config.Status:
			attachment.AddShortField("Status", issue.Status())
		case config.Type:
			attachment.AddShortField("Type", issue.IssueType())
		case config.Priority:
			attachment.AddShortField("Priority", issue.Priority().Name)
		case config.Resolution:
			attachment.AddShortField("Resolution", issue.Resolution())
		case config.Assignee:
			attachment.AddShortField("Assignee", o.GetPersonName(issue.Assignee()))
		case config.Reporter:
			attachment.AddShortField("Reporter", o.GetPersonName(issue.Reporter()))
		case config.Creator:
			attachment.AddShortField("Creator", o.GetPersonName(issue.Creator()))
		case config.Created:
			attachment.AddShortField("Created", issue.Created().Format(o.DatetimeLayout))
		case config.Updated:
			attachment.AddShortField("Updated", issue.Updated().Format(o.DatetimeLayout))
		case config.Watches:
			attachment.AddShortField("Watches", strconv.Itoa(issue.Watches()))
		case config.Components:
			attachment.AddShortField("Components", issue.Components())
		case config.Labels:
			attachment.AddShortField("Labels", issue.Labels())
		}
	}
}
