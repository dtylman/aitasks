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

func (r *Request) AddOption(opt goai.Option) {
	if r.Options == nil {
		r.Options = []goai.Option{}
	}
	r.Options = append(r.Options, opt)
}
