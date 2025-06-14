package schema

type TextContentSchema struct {
	Type string `json:"type"` // text
	Text string `json:"text"`
}

func (t *TextContentSchema) Content() any {
	return t
}

type ImageContentSchema struct {
	Type     string `json:"type"` // image
	Data     string `json:"data"` // base64 encoded image data
	MimeType string `json:"mimeType"`
}

func (i *ImageContentSchema) Content() any {
	return i
}

type AudioContentSchema struct {
	Type     string `json:"type"` // audio
	Data     string `json:"data"` // base64 encoded audio data
	MimeType string `json:"mimeType"`
}

func (a *AudioContentSchema) Content() any {
	return a
}

// The contents of a resource, embedded into a prompt or tool call result.
type EmbeddedResourceSchema struct {
	Type     string                `json:"type"` // "resource"
	Resource ResourceContentSchema `json:"resource"`
}

func (e *EmbeddedResourceSchema) Content() any {
	return e
}
