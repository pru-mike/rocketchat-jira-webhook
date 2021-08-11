package rocketchat

type Output struct {
	Alias       string       `json:"alias"`
	Avatar      string       `json:"avatar"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments,omitempty"`
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

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func (attachment *Attachment) AddShortField(title string, value string) {
	attachment.AddField(title, value, true)
}

func (attachment *Attachment) AddField(title string, value string, short bool) {
	attachment.Fields = append(attachment.Fields, Field{
		Title: title,
		Value: value,
		Short: short,
	})
}

type OutputInterface interface {
	Username() string
	IconURL() string
	GetText() string
	Attachments() []Attachment
	AttachmentsLen() int
}

func New(out OutputInterface) *Output {
	return &Output{
		Alias:       out.Username(),
		Avatar:      out.IconURL(),
		Text:        out.GetText(),
		Attachments: out.Attachments(),
	}
}
