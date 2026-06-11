package chat

import (
	"github.com/zendev-sh/goai"
	"github.com/zendev-sh/goai/provider"
)

// Request represents a chat completion request.
type Request struct {
	// Messages is the conversation history.
	Messages []provider.Message
	// Opts more options
	Options []goai.Option
}

// Response represents a chat completion response.
func (r *Request) AddOption(opt goai.Option) {
	if r.Options == nil {
		r.Options = []goai.Option{}
	}
	r.Options = append(r.Options, opt)
}

// AddMessage adds a message to the conversation history.
func (r *Request) AddMessage(role provider.Role, content string) {
	if r.Messages == nil {
		r.Messages = []provider.Message{}
	}
	r.Messages = append(r.Messages, provider.Message{
		Role: role,
		Content: []provider.Part{
			{Text: content,
				Type: provider.PartText,
			},
		},
	})
}
