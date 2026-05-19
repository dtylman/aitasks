package deepseek

// Config holds DeepSeek provider configuration.
type Config struct {
	// APIKey is the DeepSeek API key.
	APIKey string
	// Model is the default model to use (e.g., "deepseek-v4-pro", "deepseek-v4-flash").
	Model string
	// Thinking enables thinking mode (chain-of-thought reasoning). Default is false (disabled).
	Thinking bool
}
