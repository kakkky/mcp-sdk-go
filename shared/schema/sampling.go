package schema

type SamplingMessageSchema[T ContentSchema] struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content T      `json:"content"`
}

type ContentSchema interface {
	TextContentSchema | ImageContentSchema | AudioContentSchema
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
