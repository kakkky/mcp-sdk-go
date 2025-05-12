package schema

// Server
type ServerCapabilities struct {
	Experimental any `json:"experimental,omitempty"`
	*Logging     `json:"logging"`
	*Prompts     `json:"prompts,omitempty"`
	*Resources   `json:"resources,omitempty"`
	*Tools       `json:"tools,omitempty"`
}

type Logging struct{}

type Prompts struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type Resources struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}
type Tools struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// Client
type ClientCapabilities struct {
	Experimental any `json:"experimental,omitempty"`
	*Sampling    `json:"sampling"`
	*Roots       `json:"roots,omitempty"`
}

type Sampling struct{}

type Roots struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// リクエストや通知で使用される型

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

type RootSchema struct {
	Uri  string  `json:"uri"` // starting with file://
	Name *string `json:"name,omitempty"`
}
