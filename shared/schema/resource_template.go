package schema

type ResourceTemplateSchema struct {
	UriTemplate string `json:"uriTemplate"`
	Name        string `json:"name"`
	*ResourceMetadata
}
