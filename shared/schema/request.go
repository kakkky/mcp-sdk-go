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
	SystemPrompt    string                     `json:"systemPrompt,omitempty"`
	IncludeContext  string                     `json:"includeContext,omitempty"` // none or  thisServer or  allServer
	Temperature     int                        `json:"temperature,omitempty"`
	MaxTokens       int                        `json:"maxTokens"`
	StopSequences   []string                   `json:"stopSequences,omitempty"`
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

type ListResourceRequestSchema struct {
	MethodName string `json:"method"`
}

func (r *ListResourceRequestSchema) Method() string {
	return r.MethodName
}

func (r *ListResourceRequestSchema) Params() any {
	return nil
}

type ListResourceTemplatesRequestSchema struct {
	MethodName string `json:"method"`
}

func (r *ListResourceTemplatesRequestSchema) Method() string {
	return r.MethodName
}
func (r *ListResourceTemplatesRequestSchema) Params() any {
	return nil
}

type ReadResourceRequestSchema struct {
	MethodName string                    `json:"method"`
	ParamsData ReadResourceRequestParams `json:"params"`
}

type ReadResourceRequestParams struct {
	Uri string `json:"uri"`
}

func (r *ReadResourceRequestSchema) Method() string {
	return r.MethodName
}
func (r *ReadResourceRequestSchema) Params() any {
	return r.ParamsData
}

type CompleteRequestSchema struct {
	MethodName string                `json:"method"`
	ParamsData CompleteRequestParams `json:"params"`
}

type CompleteRequestParams struct {
	Ref      ReferenceSchema          `json:"ref"`
	Argument CompleteRequestParamsArg `json:"argument"`
}

func (r *CompleteRequestSchema) Method() string {
	return r.MethodName
}
func (r *CompleteRequestSchema) Params() any {
	return r.ParamsData
}

type SetLevelRequestSchema struct {
	MethodName string                       `json:"method"`
	ParamsData SetLoggingLevelRequestParams `json:"params"`
}

func (r *SetLevelRequestSchema) Method() string {
	return r.MethodName
}
func (r *SetLevelRequestSchema) Params() any {
	return r.ParamsData
}

type SetLoggingLevelRequestParams struct {
	Level LoggingLevelSchema `json:"level"`
}

type ListPromptsRequestSchema struct {
	MethodName string `json:"method"`
}

func (r *ListPromptsRequestSchema) Method() string {
	return r.MethodName
}
func (r *ListPromptsRequestSchema) Params() any {
	return nil
}

type GetPromptRequestSchema struct {
	MethodName string                 `json:"method"`
	ParamsData GetPromptRequestParams `json:"params"`
}
type GetPromptRequestParams struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments,omitempty"` // 変数名と値のマップ
}

func (r *GetPromptRequestSchema) Method() string {
	return r.MethodName
}
func (r *GetPromptRequestSchema) Params() any {
	return r.ParamsData
}
