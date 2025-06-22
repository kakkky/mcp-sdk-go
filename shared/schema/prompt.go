package schema

type PromptSchema struct {
	Name        string                `json:"name"`
	Description string                `json:"description,omitempty"`
	Auguments   []PromptAugmentSchema `json:"arguments,omitempty"`
}

type PromptAugmentSchema struct {
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	Required         bool   `json:"required,omitempty"`
	CompletionValues []string
}

type PromptMessageSchema struct {
	Role    string              `json:"role"` // "user" or "assistant"
	Content PromptContentSchema `json:"content,omitempty"`
}

// TextContentSchema | ImageContentSchema | AudioContentSchema | EmbeddedResourceSchema
type PromptContentSchema interface {
	Content() any // Returns the schema itContent, used for type assertion
}
