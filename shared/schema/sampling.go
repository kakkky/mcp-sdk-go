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
