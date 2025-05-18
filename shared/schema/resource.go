package schema

type ResourceSchema struct {
	Uri  string `json:"uri"`
	Name string `json:"name"`
	*ResourceMetadata
}

type ResourceMetadata struct {
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

type ResourceContentSchema interface {
	Uri() string
	MimeType() string
	Content() string
}

type TextResourceContentsSchema struct {
	UriData      string `json:"uri"`
	MimeTypeData string `json:"mimeType"`
	ContentData  string `json:"text"`
}

func (r *TextResourceContentsSchema) Uri() string {
	return r.UriData
}
func (r *TextResourceContentsSchema) MimeType() string {
	return r.MimeTypeData
}
func (r *TextResourceContentsSchema) Content() string {
	return r.ContentData
}

type BlobResourceContentsSchema struct {
	UriData      string `json:"uri"`
	MimeTypeData string `json:"mimeType"`
	ContentData  string `json:"blob"`
}

func (r *BlobResourceContentsSchema) Uri() string {
	return r.UriData
}
func (r *BlobResourceContentsSchema) MimeType() string {
	return r.MimeTypeData
}
func (r *BlobResourceContentsSchema) Content() string {
	return r.ContentData
}
