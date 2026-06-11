package ocr

import (
	"context"
	"fmt"

	"github.com/dtylman/aitasks/chat"
	"github.com/dtylman/aitasks/prompts"
	"github.com/zendev-sh/goai/provider"
)

// Task orchestrates OCR text cleanup.
type Task struct {
	provider provider.LanguageModel
}

// New creates a new OCR cleanup Task with the given provider.
func New(provider provider.LanguageModel) *Task {
	t := &Task{
		provider: provider,
	}
	return t
}

// Clean processes raw OCR segments and returns structured, cleaned text.
func (t *Task) Clean(ctx context.Context, req *Request) (*Response, error) {
	systemPrompt, err := prompts.Render("ocr", "default", provider.RoleSystem, "clean", req)
	if err != nil {
		return nil, fmt.Errorf("render system prompt: %w", err)
	}
	userPrompt, err := prompts.Render("ocr", "default", provider.RoleUser, "clean", req)
	if err != nil {
		return nil, fmt.Errorf("render user prompt: %w", err)
	}

	chatReq := &chat.Request{
		Messages: []provider.Message{
			{Role: provider.RoleSystem, Content: []provider.Part{{Text: systemPrompt}}},
			{Role: provider.RoleUser, Content: []provider.Part{{Text: userPrompt}}},
		},
	}

	resp, err := chat.ChatInto[Response](ctx, t.provider, chatReq)
	if err != nil {
		return nil, fmt.Errorf("chat failed: %w, resp: %v", err, resp)
	}

	return &resp.Object, nil
}
