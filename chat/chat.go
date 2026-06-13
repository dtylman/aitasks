package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zendev-sh/goai"
	"github.com/zendev-sh/goai/provider"
)

// ChatInto generates a JSON schema from the target type, attaches it to
// the request, calls Chat, and decodes the response into target.
func ChatInto[T any](ctx context.Context, llm provider.LanguageModel, req *Request) (*goai.ObjectResult[T], error) {
	schema, err := NewJSONSchema(new(T))
	if err != nil {
		return nil, fmt.Errorf("schema generation: %w", err)
	}

	req.AddMessage(provider.RoleUser, "Respond with a JSON object matching this schema: "+string(schema.RawMessage()))

	opts := req.Options
	opts = append(opts, goai.WithExplicitSchema(schema.RawMessage()))
	opts = append(opts, goai.WithMessages(req.Messages...))

	maxRetries := 3

	for maxRetries > 0 {
		tr, err := goai.GenerateText(ctx, llm, opts...)
		if err != nil {
			return nil, err
		}

		objectResult := &goai.ObjectResult[T]{
			Usage:            tr.TotalUsage,
			FinishReason:     tr.FinishReason,
			Response:         tr.Response,
			ProviderMetadata: tr.ProviderMetadata,
			ResponseMessages: tr.ResponseMessages,
			Steps:            tr.Steps,
		}

		err = json.Unmarshal([]byte(tr.Text), &objectResult.Object)
		if err == nil {
			return objectResult, nil
		}
		maxRetries--
		req.Messages = append(req.Messages, tr.ResponseMessages...)
		req.AddMessage(provider.RoleUser, "The response was not valid JSON. Please respond with a valid JSON object matching the schema. The error was: "+err.Error())
	}

	return nil, fmt.Errorf("failed to get a valid response after %d retries", maxRetries)
}
