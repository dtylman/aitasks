package ocr

import (
	"context"
	"fmt"

	"github.com/dtylman/aitasks/chat"
	"github.com/dtylman/aitasks/prompts"
)

// Task orchestrates OCR text cleanup.
type Task struct {
	client chat.Client
}

// New creates a new OCR cleanup Task with the given client and options.
func New(client chat.Client) *Task {
	t := &Task{
		client: client,
	}
	return t
}

// Clean processes raw OCR segments and returns structured, cleaned text.
func (t *Task) Clean(ctx context.Context, req *Request) (*Response, error) {
	systemPrompt, err := prompts.Render("ocr", "default", chat.RoleSystem, "clean", req)
	if err != nil {
		return nil, fmt.Errorf("render system prompt: %w", err)
	}
	userPrompt, err := prompts.Render("ocr", "default", chat.RoleUser, "clean", req)
	if err != nil {
		return nil, fmt.Errorf("render user prompt: %w", err)
	}

	chatReq := &chat.Request{
		Messages: []chat.Message{
			{Role: chat.RoleSystem, Content: systemPrompt},
			{Role: chat.RoleUser, Content: userPrompt},
		},
	}

	var result Response
	resp, err := chat.ChatInto(ctx, t.client, chatReq, &result)
	if err != nil {
		content := ""
		if resp != nil {
			content = resp.Content
		}
		return nil, fmt.Errorf("chat failed: %v, response: '%v'", err, content)
	}

	return &result, nil
}
