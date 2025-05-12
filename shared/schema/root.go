package schema

type RootSchema struct {
	Uri  string  `json:"uri"` // starting with file://
	Name *string `json:"name,omitempty"`
}
