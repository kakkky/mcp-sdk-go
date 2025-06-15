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
	isResourceHandlersInitialized   bool
	isCompletionHandlersInitialized bool
}

func NewMcpServer(serverInfo schema.Implementation, options *server.ServerOptions) *McpServer {
	return &McpServer{
		Server:                      server.NewServer(serverInfo, options),
		registeredResources:         make(map[string]*RegisteredResource),
		registeredResourceTemplates: make(map[string]*RegisteredResourceTemplate),
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
