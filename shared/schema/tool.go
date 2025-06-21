package schema

type ToolSchema struct {
	Name        string                `json:"name"`                  // Name of the tool
	Description string                `json:"description,omitempty"` // Description of the tool
	InputSchema InputSchema           `json:"inputSchema"`           // Schema for the input to the tool
	Annotations *ToolAnotationsSchema `json:"annotations,omitempty"` // Annotations for the tool
}

type InputSchema struct {
	Type       string         `json:"type"`       // "object"
	Properties PropertySchema `json:"properties"` // Properties of the input object
}

type PropertySchema map[string]PropertyInfoSchema
type PropertyInfoSchema struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ToolAnotationsSchema struct {
	Title           string `json:"title,omitempty"`           // Title of the tool
	ReadOnlyHint    bool   `json:"readOnlyHint,omitempty"`    // If true, the tool is read-only
	DestructiveHint bool   `json:"destructiveHint,omitempty"` // If true, the tool is destructive
	IdempotentHint  bool   `json:"idempotentHint,omitempty"`  // If true, the tool is idempotent
	OpenWorldHint   bool   `json:"openWorldHint,omitempty"`   // If true, the tool is open world
}
