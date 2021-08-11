package rocketchat

import "strings"

type MultiplexOutputBuilder struct{}

func NewMultiplexOutputBuilder() *MultiplexOutputBuilder {
	return &MultiplexOutputBuilder{}
}

var _ OutputInterface = (*MultiplexOutput)(nil)

type MultiplexOutput struct {
	output1, output2 OutputInterface
}

func (o *MultiplexOutput) Username() string {
	return o.output1.Username()
}

func (o *MultiplexOutput) IconURL() string {
	return o.output1.IconURL()
}

func (o *MultiplexOutput) AttachmentsLen() int {
	return o.output1.AttachmentsLen() + o.output2.AttachmentsLen()
}

func (o *MultiplexOutput) GetText() string {
	var texts []string
	if o.output1.AttachmentsLen() != 0 {
		texts = append(texts, o.output1.GetText())
	}
	if o.output2.AttachmentsLen() != 0 {
		texts = append(texts, o.output2.GetText())
	}
	return strings.Join(texts, ", ")
}

func (o *MultiplexOutput) Attachments() []Attachment {
	var attachments []Attachment
	attachments = append(attachments, o.output1.Attachments()...)
	attachments = append(attachments, o.output2.Attachments()...)
	return attachments
}

func (o *MultiplexOutputBuilder) Output(o1, o2 OutputInterface) *Output {
	return New(&MultiplexOutput{
		o1, o2,
	})
}
