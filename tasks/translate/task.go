package translate

import (
	"context"
	"fmt"

	"github.com/dtylman/aitasks/chat"
	"github.com/dtylman/aitasks/prompts"
	"github.com/zendev-sh/goai"
	"github.com/zendev-sh/goai/provider"
)

// Task orchestrates translation workflows.
type Task struct {
	// provider is the language model provider used for translation and related tasks.
	provider provider.LanguageModel
	//AutoProofread if true, automatically runs a proofreading step after translation to improve quality.
	AutoProofread bool
	//MaxRetries specifies how many times to retry the translation step on failure before giving up.
	MaxRetries int
}

// New creates a new translation Task with the given provider and options.
func New(provider provider.LanguageModel) *Task {
	t := &Task{
		provider:      provider,
		AutoProofread: true,
		MaxRetries:    3,
	}
	return t
}

// Translate translates a single paragraph.
func (t *Task) Translate(ctx context.Context, req *Request) (*Result, error) {
	result, err := t.doTranslate(ctx, req)
	if err != nil {
		return nil, err
	}

	if t.AutoProofread {
		result, err = t.doProofread(ctx, req, result.Translation)
		if err != nil {
			return nil, fmt.Errorf("proofread: %w", err)
		}
	}

	return result, nil
}

// Proofread proofreads an existing translation.
func (t *Task) Proofread(ctx context.Context, req *Request, translation string) (*Result, error) {
	return t.doProofread(ctx, req, translation)
}

// Fix re-translates a paragraph that was flagged as poor quality.
func (t *Task) Fix(ctx context.Context, req *Request, badTranslation string) (*Result, error) {
	return t.doFix(ctx, req, badTranslation)
}

// PopulateProject uses the LLM to fill in project context details (genre, synopsis,
// writing style, glossary, characters) given a ProjectContext with Title and Author set.
func (t *Task) PopulateProject(ctx context.Context, project *ProjectContext) (*ProjectContext, error) {
	systemPrompt, err := prompts.Render("translate", "default", provider.RoleSystem, "populate_project", project)
	if err != nil {
		return nil, err
	}
	userPrompt, err := prompts.Render("translate", "default", provider.RoleUser, "populate_project", project)
	if err != nil {
		return nil, err
	}

	var chatReq chat.Request
	chatReq.AddMessage(provider.RoleSystem, systemPrompt)
	chatReq.AddMessage(provider.RoleUser, userPrompt)

	resp, err := chat.ChatInto[ProjectContext](ctx, t.provider, &chatReq)
	if err != nil {
		return nil, fmt.Errorf("populate project: %w. resp: %v", err, resp)
	}

	// Preserve the original title and author
	resp.Object.Title = project.Title
	resp.Object.Author = project.Author

	return &resp.Object, nil
}

func (t *Task) doTranslate(ctx context.Context, req *Request) (*Result, error) {
	if req.Text == "" {
		return &Result{Translation: ""}, nil
	}
	if req.Style == "" {
		req.Style = "default"
	}
	systemPrompt, err := prompts.Render("translate", req.Style, provider.RoleSystem, "translate", req)
	if err != nil {
		return nil, err
	}
	userPrompt, err := prompts.Render("translate", req.Style, provider.RoleUser, "translate", req)
	if err != nil {
		return nil, err
	}

	chatReq := &chat.Request{
		Messages: []provider.Message{
			{Role: provider.RoleSystem, Content: []provider.Part{{Text: systemPrompt}}},
			{Role: provider.RoleUser, Content: []provider.Part{{Text: userPrompt}}},
		},
	}
	chatReq.AddOption(goai.WithMaxRetries(t.MaxRetries))

	resp, err := chat.ChatInto[Result](ctx, t.provider, chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to translate: %w, resp: %v:", err, resp)
	}
	return &resp.Object, nil
}

func (t *Task) doProofread(ctx context.Context, tr *Request, translation string) (*Result, error) {
	req := &ProofreadRequest{
		TranslationReq: tr,
		DraftText:      translation,
	}

	systemPrompt, err := prompts.Render("translate", "default", provider.RoleSystem, "proofread", req)
	if err != nil {
		return nil, err
	}

	userPrompt, err := prompts.Render("translate", "default", provider.RoleUser, "proofread", req)
	if err != nil {
		return nil, err
	}

	chatReq := &chat.Request{
		Messages: []provider.Message{
			{Role: provider.RoleSystem, Content: []provider.Part{{Text: systemPrompt}}},
			{Role: provider.RoleUser, Content: []provider.Part{{Text: userPrompt}}},
		},
	}
	chatReq.AddOption(goai.WithMaxRetries(t.MaxRetries))

	resp, err := chat.ChatInto[Result](ctx, t.provider, chatReq)
	if err != nil {
		return nil, fmt.Errorf("proofread: %w, resp: %v", err, resp)
	}

	return &resp.Object, nil
}

func (t *Task) doFix(ctx context.Context, req *Request, badTranslation string) (*Result, error) {
	fixReq := &FixRequest{
		TranslationReq: req,
		DraftText:      badTranslation,
	}

	systemPrompt, err := prompts.Render("translate", "default", provider.RoleSystem, "fix", fixReq)
	if err != nil {
		return nil, err
	}
	userPrompt, err := prompts.Render("translate", "default", provider.RoleUser, "fix", fixReq)
	if err != nil {
		return nil, err
	}

	chatReq := &chat.Request{
		Messages: []provider.Message{
			{Role: provider.RoleSystem, Content: []provider.Part{{Text: systemPrompt}}},
			{Role: provider.RoleUser, Content: []provider.Part{{Text: userPrompt}}},
		},
	}
	chatReq.AddOption(goai.WithMaxRetries(t.MaxRetries))

	resp, err := chat.ChatInto[Result](ctx, t.provider, chatReq)
	if err != nil {
		return nil, fmt.Errorf("fix: %w, resp: %v", err, resp)
	}

	return &resp.Object, nil
}
