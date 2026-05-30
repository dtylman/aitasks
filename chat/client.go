package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// Client is the interface that provider adapters implement.
// When req.Schema is non-nil, the provider should request JSON output
// from the model using whatever native mechanism it supports.
type Client interface {
	Chat(ctx context.Context, req *Request) (*Response, error)
}

// ChatInto generates a JSON schema from the target type, attaches it to
// the request, calls Chat, and decodes the response into target.
func ChatInto(ctx context.Context, c Client, req *Request, target any) (*Response, error) {
	schema, err := NewJSONSchema(target)
	if err != nil {
		return nil, fmt.Errorf("schema generation: %w", err)
	}
	req.Schema = schema

	maxRetries := req.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}
	for attempt := range maxRetries {
		resp, err := Chat(ctx, c, req)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(resp.Content), target)
		if err != nil {
			if attempt == maxRetries-1 {
				return resp, fmt.Errorf("decode response: %w, %v", err, resp.Content)
			}
			log.Printf("ChatInto: decode error (attempt %d/%d): %v, retrying...", attempt+1, maxRetries, err)
			// Append the failed exchange and ask the model to fix the JSON
			req.Messages = append(req.Messages,
				Message{Role: RoleAssistant, Content: resp.Content},
				Message{Role: RoleUser, Content: fmt.Sprintf("Your response was not valid JSON: %v. Please return only valid JSON matching the schema.", err)},
			)
			continue
		}
		return resp, nil
	}
	// unreachable
	return nil, fmt.Errorf("ChatInto: exhausted retries")
}

func Chat(ctx context.Context, c Client, req *Request) (*Response, error) {
	schemaType := "none"
	if req.Schema != nil {
		schemaType = req.Schema.Type
	}
	message := fmt.Sprintf("invoking %T:%v chat with %v:%v", c, req.Model, schemaType, len(req.Messages))
	log.Println(message)
	logChat(message, req)
	resp, err := c.Chat(ctx, req)
	if err != nil {
		return nil, err
	}
	logChat("response", resp)
	logChat("==============================================================", nil)
	return resp, nil
}
