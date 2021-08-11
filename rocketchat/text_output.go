package rocketchat

import (
	"github.com/pru-mike/rocketchat-jira-webhook/config"
)

type TextOutputBuilder struct {
	*OutputBuilder
}

func NewTextOutputBuilder(cfg *config.Message) *TextOutputBuilder {
	return &TextOutputBuilder{
		NewOutputBuilder(cfg),
	}
}

var _ OutputInterface = (*TextOutput)(nil)

type TextOutput struct {
	*TextOutputBuilder
	text string
}

func (o *TextOutput) GetText() string {
	return o.text
}

func (o *TextOutput) AttachmentsLen() int {
	return 0
}

func (o *TextOutput) Attachments() []Attachment {
	return nil
}

func (o *TextOutputBuilder) Output(text string) *Output {
	return New(&TextOutput{
		TextOutputBuilder: o,
		text:              text,
	})

}
