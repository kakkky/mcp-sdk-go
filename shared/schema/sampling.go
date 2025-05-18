package schema

type SamplingMessageSchema[T ContentSchema] struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content T      `json:"content"`
}

type ContentSchema interface {
	TextContentSchema | ImageContentSchema | AudioContentSchema
}

type TextContentSchema struct {
	Type string `json:"type"` // text
	Text string `json:"text"`
}

type ImageContentSchema struct {
	Type     string `json:"type"` // image
	Data     string `json:"data"` // base64 encoded image data
	MimeType string `json:"mimeType"`
}

type AudioContentSchema struct {
	Type     string `json:"type"` // audio
	Data     string `json:"data"` // base64 encoded audio data
	MimeType string `json:"mimeType"`
}

type ModelPreferencesSchema struct {
	Hints                []*ModelHintsSchema `json:"hints,omitempty"`
	CostPriority         *int                `json:"costPriority,omitempty"`         // 0 <= costPriority <= 1
	SpeedPriority        *int                `json:"speedPriority,omitempty"`        // 0 <= speedPriority <= 1
	IntelligencePriority *int                `json:"intelligencePriority,omitempty"` // 0 <= intelligencePriority <= 1
}
type ModelHintsSchema struct {
	Name string `json:"name,omitempty"`
}
