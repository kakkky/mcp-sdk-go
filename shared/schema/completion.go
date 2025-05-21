package schema

type ReferenceSchema interface {
	Type() string
	UriOrName() string
}

type ResourceReferenceSchema struct {
	TypeData string `json:"type"`
	UriData  string `json:"uri"`
}

func (r *ResourceReferenceSchema) Type() string {
	return r.TypeData
}
func (r *ResourceReferenceSchema) UriOrName() string {
	return r.UriData
}

type PromptReferenceSchema struct {
	TypeData string `json:"type"`
	NameData string `json:"name"`
}

func (r *PromptReferenceSchema) Type() string {
	return r.TypeData
}

func (r *PromptReferenceSchema) UriOrName() string {
	return r.NameData
}

type CompleteRequestParamsArg struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CompletionSchema struct {
	Values  []string `json:"values"`
	Total   int      `json:"total,omitempty"`
	HasMore *bool    `json:"hasMore,omitempty"`
}
