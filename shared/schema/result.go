package schema

type Result interface {
	Result() any
}

// initialize
type InitializeResultSchema struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      Implementation     `json:"serverInfo"`
	Instructions    string             `json:"instructions,omitempty"`
}

func (r *InitializeResultSchema) Result() any {
	return r
}

// 何も返さない場合に使用する
type EmptyResultSchema struct{}

func (r *EmptyResultSchema) Result() any {
	return r
}

// sampling/createMessage
type CreateMessageResultSchema[T ContentSchema] struct {
	Model      string `json:"model"`
	StopReason string `json:"stopReason,omitempty"` // endTurn or stopSequence or maxTokens
	Role       string `json:"role"`                 // user or assistant
	Content    T      `json:"content"`
}

func (r *CreateMessageResultSchema[T]) Result() any {
	return r
}

// roots/list
type ListRootsResultSchema struct {
	Roots []RootSchema `json:"roots"`
}

func (r *ListRootsResultSchema) Result() any {
	return r
}

// resources/read
type ReadResourceResultSchema struct {
	Contents []ResourceContentSchema `json:"contents"`
}

func (r *ReadResourceResultSchema) Result() any {
	return r
}

// resources/list
type ListResourcesResultSchema struct {
	Resources []ResourceSchema `json:"resources"`
}

func (r *ListResourcesResultSchema) Result() any {
	return r
}

// resources/templates/list
type ListResourceTemplatesResultSchema struct {
	ResourceTemplates []ResourceTemplateSchema `json:"resourceTemplates"`
}

func (r *ListResourceTemplatesResultSchema) Result() any {
	return r
}

// completion/complete
type CompleteResultSchema struct {
	Completion CompletionSchema `json:"completion"`
}

func (r *CompleteResultSchema) Result() any {
	return r
}

// prompts/list
type GetPromptResultSchema struct {
	Description string                `json:"description,omitempty"`
	Messages    []PromptMessageSchema `json:"messages"`
}

func (r *GetPromptResultSchema) Result() any {
	return r
}

// prompts/list
type ListPromptsResultSchema struct {
	Prompts []PromptSchema `json:"prompts"`
}

func (r *ListPromptsResultSchema) Result() any {
	return r
}

// tools/list
type ListToolsResultSchema struct {
	Tools []ToolSchema `json:"tools"`
}

func (r *ListToolsResultSchema) Result() any {
	return r
}

// tools/call
type CallToolResultSchema struct {
	Content                           []ToolContentSchema `json:"content"`
	IsError                           bool                `json:"isError,omitempty"`
	CompatibilityCallToolResultSchema                     // Deprecated: use Content instead
}

func (r *CallToolResultSchema) Result() any {
	return r
}

// CallToolResultSchema extended with backwards compatibility to protocol version 2024-10-07.
type CompatibilityCallToolResultSchema struct {
	ToolResult any `json:"toolResult,omitempty"` // Deprecated: use ToolResults instead
}
