package schema

type Request interface {
	Method() string
	Params() any
}

// initialize
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

// ping
type PingRequestSchema struct {
	MethodName string `json:"method"`
}

func (r *PingRequestSchema) Method() string {
	return r.MethodName
}

func (r *PingRequestSchema) Params() any {
	return nil
}

// sampling/createMessage
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

// roots/list
type ListRootsRequestSchema struct {
	MethodName string `json:"method"`
}

func (r *ListRootsRequestSchema) Method() string {
	return r.MethodName
}

func (r *ListRootsRequestSchema) Params() any {
	return nil
}

// resources/list
type ListResourceRequestSchema struct {
	MethodName string `json:"method"`
}

func (r *ListResourceRequestSchema) Method() string {
	return r.MethodName
}

func (r *ListResourceRequestSchema) Params() any {
	return nil
}

// resources/templates/list
type ListResourceTemplatesRequestSchema struct {
	MethodName string `json:"method"`
}

func (r *ListResourceTemplatesRequestSchema) Method() string {
	return r.MethodName
}
func (r *ListResourceTemplatesRequestSchema) Params() any {
	return nil
}

// resources/read
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

// resources/subscribe
type SubscribeRequestSchema struct {
	MethodName string                 `json:"method"`
	ParamsData SubscribeRequestParams `json:"params"`
}
type SubscribeRequestParams struct {
	Uri string `json:"uri"`
}

func (r *SubscribeRequestSchema) Method() string {
	return r.MethodName
}
func (r *SubscribeRequestSchema) Params() any {
	return r.ParamsData
}

// resources/unsubscribe
type UnsubscribeRequestSchema struct {
	MethodName string                   `json:"method"`
	ParamsData UnsubscribeRequestParams `json:"params"`
}
type UnsubscribeRequestParams struct {
	Uri string `json:"uri"`
}

func (r *UnsubscribeRequestSchema) Method() string {
	return r.MethodName
}
func (r *UnsubscribeRequestSchema) Params() any {
	return r.ParamsData
}

// completion/complete
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

// logging/setLevel
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

// prompts/list
type ListPromptsRequestSchema struct {
	MethodName string `json:"method"`
}

func (r *ListPromptsRequestSchema) Method() string {
	return r.MethodName
}
func (r *ListPromptsRequestSchema) Params() any {
	return nil
}

// prompts/get
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

// tools/list
type ListToolsRequestSchema struct {
	MethodName string `json:"method"`
}

func (r *ListToolsRequestSchema) Method() string {
	return r.MethodName
}
func (r *ListToolsRequestSchema) Params() any {
	return nil
}

// tools/call
type CallToolRequestSchema struct {
	MethodName string                `json:"method"`
	ParamsData CallToolRequestParams `json:"params"`
}
type CallToolRequestParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"` // 変数名と値のマップ
}

func (r *CallToolRequestSchema) Method() string {
	return r.MethodName
}
func (r *CallToolRequestSchema) Params() any {
	return r.ParamsData
}
