package schema

type Request interface {
	Method() string
	Params() any
}

type InitializeRequestSchema struct {
	MethodName string                  `json:"method"`
	ParamsData InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      Implementation     `json:"clientInfo"`
}

func (r *InitializeRequestSchema) Method() string {
	return r.MethodName
}

func (r *InitializeRequestSchema) Params() any {
	return r.ParamsData
}

type PingRequestSchema struct {
	MethodName string `json:"method"`
}

func (r *PingRequestSchema) Method() string {
	return r.MethodName
}

func (r *PingRequestSchema) Params() any {
	return nil
}

type CreateMessageRequestSchema[T ContentSchema] struct {
	MethodName string                        `json:"method"`
	ParamsData CreateMessageRequestParams[T] `json:"params"`
}

type CreateMessageRequestParams[T ContentSchema] struct {
	Messages        []SamplingMessageSchema[T] `json:"messages"`
	SystemPrompt    *string                    `json:"systemPrompt,omitempty"`
	IncludeContext  *string                    `json:"includeContext,omitempty"` // none or  thisServer or  allServer
	Temperature     *int                       `json:"temperature,omitempty"`
	MaxTokens       int                        `json:"maxTokens"`
	StopSequences   []*string                  `json:"stopSequences,omitempty"`
	Metadata        any                        `json:"metadata,omitempty"`
	ModelPreference *ModelPreferencesSchema    `json:"modelPreference,omitempty"`
}

func (r *CreateMessageRequestSchema[T]) Method() string {
	return r.MethodName
}

func (r *CreateMessageRequestSchema[T]) Params() any {
	return r.ParamsData
}

type ListRootsRequestSchema struct {
	MethodName string `json:"method"`
}

func (r *ListRootsRequestSchema) Method() string {
	return r.MethodName
}

func (r *ListRootsRequestSchema) Params() any {
	return nil
}
