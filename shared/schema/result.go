package schema

type Result interface {
	Result() any
}

type InitializeResultSchema struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      Implementation     `json:"serverInfo"`
	Instructions    *string            `json:"instructions,omitempty"`
}

func (r *InitializeResultSchema) Result() any {
	return r
}

type EmptyResultSchema struct{}

func (r *EmptyResultSchema) Result() any {
	return r
}

type CreateMessageResultSchema[T ContentSchema] struct {
	Model      string  `json:"model"`
	StopReason *string `json:"stopReason,omitempty"` // endTurn or stopSequence or maxTokens
	Role       string  `json:"role"`                 // user or assistant
	Content    T       `json:"content"`
}

func (r *CreateMessageResultSchema[T]) Result() any {
	return r
}

type ListRootResultSchema struct {
	Roots []RootSchema `json:"roots"`
}

func (r *ListRootResultSchema) Result() any {
	return r
}

type ReadResourceResultSchema[T ResourceContentSchema] struct {
	Contents []T `json:"contents"`
}

func (r *ReadResourceResultSchema[T]) Result() any {
	return r
}

type ListResourcesResultSchema struct {
	Resources []ResourceSchema `json:"resources"`
}
