package chat

import (
	"context"
	"fmt"

	"github.com/zendev-sh/goai"
	"github.com/zendev-sh/goai/provider"
)

// ChatInto generates a JSON schema from the target type, attaches it to
// the request, calls Chat, and decodes the response into target.
func ChatInto[T any](ctx context.Context, provider provider.LanguageModel, req *Request, target T) (*goai.ObjectResult[T], error) {
	schema, err := NewJSONSchema(target)
	if err != nil {
		return nil, fmt.Errorf("schema generation: %w", err)
	}

	opts := req.Options
	opts = append(opts, goai.WithExplicitSchema(schema.RawMessage()))
	opts = append(opts, goai.WithMessages(req.Messages...))
	return goai.GenerateObject[T](ctx, provider, opts...)
}
