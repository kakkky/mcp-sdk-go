package mcpserver

import (
	"errors"
	"fmt"

	"github.com/kakkky/mcp-sdk-go/mcp-server/server"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

// リソース、ツール、プロンプトを扱うための、よりシンプルなAPIを提供する高レベルのMCPサーバー
// 通知の送信やカスタムリクエストハンドラーの設定など、より高度な使用を行いたい場合は、
// Serverプロパティ経由で利用できる下位の Server インスタンスを使用する必要がある
type McpServer struct {
	Server                          *server.Server
	registeredResources             map[string]*RegisteredResource
	registeredResourceTemplates     map[string]*RegisteredResourceTemplate
	registerdTools                  map[string]*RegisteredTool
	registeredPrompts               map[string]*RegisteredPrompt
	isResourceHandlersInitialized   bool
	isToolHandlersInitialized       bool
	isPromptHandlersInitialized     bool
	isCompletionHandlersInitialized bool
}

func NewMcpServer(serverInfo schema.Implementation, options *server.ServerOptions) *McpServer {
	return &McpServer{
		Server:                      server.NewServer(serverInfo, options),
		registeredResources:         make(map[string]*RegisteredResource),
		registeredResourceTemplates: make(map[string]*RegisteredResourceTemplate),
		registerdTools:              make(map[string]*RegisteredTool),
		registeredPrompts:           make(map[string]*RegisteredPrompt),
	}
}

func (m *McpServer) Connect(transport protocol.Transport) error {
	return m.Server.Connect(transport)
}

func (m *McpServer) Close() error {
	return m.Server.Close()
}

func (m *McpServer) isConnected() bool {
	return m.Server.Transport() != nil
}

// Resource はURIベースのリソースを登録します
func (m *McpServer) Resource(
	name string,
	uri string,
	metadata *schema.ResourceMetadata,
	readResourceCallBack ReadResourceCallback[schema.ResourceContentSchema],
) (*RegisteredResource, error) {
	if uri == "" {
		return nil, errors.New("uri is required")
	}
	if readResourceCallBack == nil {
		return nil, errors.New("readResourceCallBack is required")
	}
	if m.registeredResources[uri] != nil {
		return nil, fmt.Errorf("resource %s is already registered", uri)
	}

	uriPtr := &uri
	registeredResource := RegisteredResource{
		name:         name,
		metadata:     metadata,
		readCallback: readResourceCallBack,
		enabled:      true,
		Disable: func() {
			if _, ok := m.registeredResources[*uriPtr]; !ok {
				fmt.Println("resource not found")
				return
			}
			disabled := false
			m.registeredResources[*uriPtr].Update(ResourceUpdates{Enabled: &disabled})
		},
		Enable: func() {
			if _, ok := m.registeredResources[*uriPtr]; !ok {
				fmt.Println("resource not found")
				return
			}
			enabled := true
			m.registeredResources[*uriPtr].Update(ResourceUpdates{Enabled: &enabled})
		},
		Remove: func() { delete(m.registeredResources, *uriPtr) },
		Update: func(updates ResourceUpdates) {
			if _, ok := m.registeredResources[*uriPtr]; !ok {
				fmt.Println("resource not found")
				return
			}
			if updates.Uri != "" && updates.Uri != *uriPtr {
				resourceCopy := m.registeredResources[*uriPtr]
				delete(m.registeredResources, *uriPtr)
				m.registeredResources[updates.Uri] = resourceCopy

				// 参照する値を更新することで、uriの値が変わった場合でも、以降の処理 & Disable/Enable/Removeが正しく動作するようにする
				*uriPtr = updates.Uri
			}
			if updates.Name != "" {
				m.registeredResources[*uriPtr].name = updates.Name
			}
			if updates.Metadata != nil {
				m.registeredResources[*uriPtr].metadata = updates.Metadata
			}
			if updates.Callback != nil {
				m.registeredResources[*uriPtr].readCallback = *updates.Callback
			}
			if updates.Enabled != nil {
				m.registeredResources[*uriPtr].enabled = *updates.Enabled
			}
			_ = m.Server.SendResourceListChanged()
		},
	}
	m.registeredResources[*uriPtr] = &registeredResource

	_ = m.setResourceRequestHandlers()

	_ = m.Server.SendResourceListChanged
	return &registeredResource, nil
}

// ResourceTemplate はテンプレートベースのリソースを登録します
func (m *McpServer) ResourceTemplate(
	name string,
	template *ResourceTemplate,
	metadata *schema.ResourceMetadata,
	readResourceTemplateCallBack ReadResourceTemplateCallback[schema.ResourceContentSchema],
) (*RegisteredResourceTemplate, error) {
	if template == nil {
		return nil, errors.New("template is required")
	}
	if readResourceTemplateCallBack == nil {
		return nil, errors.New("readResourceTemplateCallBack is required")
	}
	if m.registeredResourceTemplates[name] != nil {
		return nil, fmt.Errorf("resource template %s is already registered", name)
	}

	namePtr := &name
	registeredResourceTemplate := RegisteredResourceTemplate{
		resourceTemplate: template,
		metadata:         metadata,
		readCallback:     readResourceTemplateCallBack,
		enabled:          true,
		Disable: func() {
			if _, ok := m.registeredResourceTemplates[*namePtr]; !ok {
				fmt.Println("resource template not found")
				return
			}
			disabled := false
			m.registeredResourceTemplates[*namePtr].Update(ResourceTemplateUpdates{Enabled: &disabled})
		},
		Enable: func() {
			if _, ok := m.registeredResourceTemplates[*namePtr]; !ok {
				fmt.Println("resource template not found")
				return
			}
			enabled := true
			m.registeredResourceTemplates[*namePtr].Update(ResourceTemplateUpdates{Enabled: &enabled})
		},
		Remove: func() { delete(m.registeredResourceTemplates, *namePtr) },
		Update: func(updates ResourceTemplateUpdates) {
			if _, ok := m.registeredResourceTemplates[*namePtr]; !ok {
				fmt.Println("resource template not found")
				return
			}
			if updates.Name != "" && updates.Name != *namePtr {
				resourceTemplateCopy := m.registeredResourceTemplates[*namePtr]
				delete(m.registeredResourceTemplates, *namePtr)
				m.registeredResourceTemplates[updates.Name] = resourceTemplateCopy

				// 参照する値を更新することで、nameの値が変わった場合でも、以降の処理 & Disable/Enable/Removeが正しく動作するようにする
				*namePtr = updates.Name
			}
			if updates.Template != nil {
				m.registeredResourceTemplates[*namePtr].resourceTemplate = updates.Template
			}
			if updates.Metadata != nil {
				m.registeredResourceTemplates[*namePtr].metadata = updates.Metadata
			}
			if updates.Callback != nil {
				m.registeredResourceTemplates[*namePtr].readCallback = *updates.Callback
			}
			if updates.Enabled != nil {
				m.registeredResourceTemplates[*namePtr].enabled = *updates.Enabled
			}
			_ = m.Server.SendResourceListChanged()
		},
	}
	m.registeredResourceTemplates[*namePtr] = &registeredResourceTemplate
	_ = m.setResourceRequestHandlers()
	m.sendResourceListChanged()
	return &registeredResourceTemplate, nil
}

func (m *McpServer) Tool(
	name string,
	description string,
	propertySchema schema.PropertySchema,
	annotations *schema.ToolAnotationsSchema,
	callback ToolCallback,
) (*RegisteredTool, error) {
	if m.registerdTools[name] != nil {
		return nil, fmt.Errorf("tool %s is already registered", name)
	}
	namePtr := &name
	registeredTool := RegisteredTool{
		description:    description,
		propertySchema: propertySchema,
		annotations:    annotations,
		callback:       callback,
		enabled:        true,
		Disable: func() {
			m.registerdTools[*namePtr].Update(ToolUpdates{Enabled: &[]bool{false}[0]})
		},
		Enable: func() {
			m.registerdTools[*namePtr].Update(ToolUpdates{Enabled: &[]bool{true}[0]})
		},
		Remove: func() { delete(m.registerdTools, *namePtr) },
		Update: func(updates ToolUpdates) {
			if _, ok := m.registerdTools[*namePtr]; !ok {
				fmt.Println("tool not found")
				return
			}
			if updates.Name != "" && updates.Name != *namePtr {
				toolCopy := m.registerdTools[*namePtr]
				delete(m.registerdTools, *namePtr)
				m.registerdTools[updates.Name] = toolCopy

				// 参照する値を更新することで、nameの値が変わった場合でも、以降の処理 & Disable/Enable/Removeが正しく動作するようにする
				*namePtr = updates.Name
			}
			if updates.Description != "" {
				m.registerdTools[*namePtr].description = updates.Description
			}
			if updates.ParamsSchema != nil {
				m.registerdTools[*namePtr].propertySchema = updates.ParamsSchema
			}
			if updates.callback != nil {
				m.registerdTools[*namePtr].callback = updates.callback
			}
			if updates.Annotations != nil {
				m.registerdTools[*namePtr].annotations = updates.Annotations
			}
			if updates.Enabled != nil {
				m.registerdTools[*namePtr].enabled = *updates.Enabled
			}
			m.sendToolListChanged()
		},
	}
	m.registerdTools[*namePtr] = &registeredTool

	_ = m.setToolRequestHandlers()
	m.sendToolListChanged()
	return &registeredTool, nil
}

func (m *McpServer) Prompt(
	name string,
	description string,
	argsSchema []schema.PromptAugmentSchema,
	callback PromptCallback,
) (*RegisteredPrompt, error) {
	if m.registeredPrompts[name] != nil {
		return nil, fmt.Errorf("prompt %s is already registered", name)
	}
	namePtr := &name
	registeredPrompt := RegisteredPrompt{
		description: description,
		argsSchema:  argsSchema,
		callback:    callback,
		enabled:     true,
		Disable: func() {
			if _, ok := m.registeredPrompts[*namePtr]; !ok {
				fmt.Println("prompt not found")
				return
			}
			disabled := false
			m.registeredPrompts[*namePtr].Update(PromptUpdates{Enabled: &disabled})
		},
		Enable: func() {
			if _, ok := m.registeredPrompts[*namePtr]; !ok {
				fmt.Println("prompt not found")
				return
			}
			enabled := true
			m.registeredPrompts[*namePtr].Update(PromptUpdates{Enabled: &enabled})
		},
		Remove: func() { delete(m.registeredPrompts, *namePtr) },
		Update: func(updates PromptUpdates) {
			if _, ok := m.registeredPrompts[*namePtr]; !ok {
				fmt.Println("prompt not found")
				return
			}
			if updates.Name != "" && updates.Name != *namePtr {
				promptCopy := m.registeredPrompts[*namePtr]
				delete(m.registeredPrompts, *namePtr)
				m.registeredPrompts[updates.Name] = promptCopy
				// 参照する値を更新することで、nameの値が変わった場合でも、以降の処理 & Disable/Enable/Removeが正しく動作するようにする
				*namePtr = updates.Name
			}
			if updates.Description != "" {
				m.registeredPrompts[*namePtr].description = updates.Description
			}
			if updates.ArgsSchema != nil {
				m.registeredPrompts[*namePtr].argsSchema = updates.ArgsSchema
			}
			if updates.Callback != nil {
				m.registeredPrompts[*namePtr].callback = updates.Callback
			}
			if updates.Enabled != nil {
				m.registeredPrompts[*namePtr].enabled = *updates.Enabled
			}
			m.sendPromptListChanged()
		},
	}
	m.registeredPrompts[*namePtr] = &registeredPrompt
	_ = m.setPromptRequestHandlers()
	m.sendPromptListChanged()
	return &registeredPrompt, nil
}
