package mcpserver

import (
	"errors"
	"fmt"

	"github.com/kakkky/mcp-sdk-go/mcp-server/server"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type McpServer struct {
	server                        *server.Server
	registeredResources           map[string]*RegisteredResource
	registeredResourceTemplates   map[string]*RegisteredResourceTemplate
	isResourceHandlersInitialized bool
}

func NewMcpServer(serverInfo schema.Implementation, options *server.ServerOptions) *McpServer {
	return &McpServer{
		server:                      server.NewServer(serverInfo, options),
		registeredResources:         make(map[string]*RegisteredResource),
		registeredResourceTemplates: make(map[string]*RegisteredResourceTemplate),
	}
}

func (m *McpServer) Connect(transport protocol.Transport) error {
	return m.server.Connect(transport)
}

func (m *McpServer) Close() error {
	return m.server.Close()
}

func (m *McpServer) isConnected() bool {
	return m.server.Transport() != nil
}

// uriかtemplateのどちらかに値を渡す
// uriを渡す場合は、readResourceCallBackを渡す
// templateを渡す場合は、readResourceTemplateCallBackを渡す
func (m *McpServer) Resource(
	name *string,
	uri *string,
	template *ResourceTemplate,
	metadata *schema.ResourceMetadata,
	readResourceCallBack ReadResourceCallback[schema.ResourceContentSchema], //  uriを渡す場合
	readResourceTemplateCallBack ReadResourceTemplateCallback[schema.ResourceContentSchema], // templateを渡す場合
) (*RegisteredResource, *RegisteredResourceTemplate, error) {
	if (uri == nil && template == nil) || (uri != nil && template != nil) {
		return nil, nil, errors.New("please provide a value for either uri or template")
	}
	// uriが渡された場合
	if uri != nil {
		if readResourceCallBack == nil {
			return nil, nil, errors.New("readResourceCallBack is required when uri is provided")
		}
		if m.registeredResources[*uri] != nil {
			return nil, nil, fmt.Errorf("resource %s is already registered", *uri)
		}

		registeredResource := RegisteredResource{
			name:         *name,
			metadata:     metadata,
			readCallback: readResourceCallBack,
			enabled:      true,
			Disable: func() {
				if _, ok := m.registeredResources[*uri]; !ok {
					fmt.Println("resource not found")
					return
				}
				disabled := false
				m.registeredResources[*uri].Update(ResourceUpdates{Enabled: &disabled})
			},
			Enable: func() {
				if _, ok := m.registeredResources[*uri]; !ok {
					fmt.Println("resource not found")
					return
				}
				enabled := true
				m.registeredResources[*uri].Update(ResourceUpdates{Enabled: &enabled})
			},
			Remove: func() { delete(m.registeredResources, *uri) },
			Update: func(updates ResourceUpdates) {
				if _, ok := m.registeredResources[*uri]; !ok {
					fmt.Println("resource not found")
					return
				}
				if updates.Uri != nil && *updates.Uri != *uri {
					resourceCopy := m.registeredResources[*uri]
					delete(m.registeredResources, *uri)
					m.registeredResources[*updates.Uri] = resourceCopy

					// 参照を更新することで、uriの値が変わった場合でも、以降の処理 & Disable/Enable/Removeが正しく動作するようにする
					uri = updates.Uri
				}
				if updates.Name != nil {
					m.registeredResources[*uri].name = *updates.Name
				}
				if updates.Metadata != nil {
					m.registeredResources[*uri].metadata = updates.Metadata
				}
				if updates.Callback != nil {
					m.registeredResources[*uri].readCallback = *updates.Callback
				}
				if updates.Enabled != nil {
					m.registeredResources[*uri].enabled = *updates.Enabled
				}
				m.server.SendResourceListChanged()
			},
		}
		m.registeredResources[*uri] = &registeredResource
		m.setResourceRequestHandlers()
		m.sendResourceListChanged()
		return &registeredResource, nil, nil
	}
	// templateが渡された場合
	if template != nil {
		if m.registeredResourceTemplates[*name] != nil {
			return nil, nil, fmt.Errorf("resource template %s is already registered", *name)
		}
		if readResourceTemplateCallBack == nil {
			return nil, nil, errors.New("readResourceTemplateCallBack is required when template is provided")
		}

		registeredResourceTemplate := RegisteredResourceTemplate{
			resourceTemplate: template,
			metadata:         metadata,
			readCallback:     readResourceTemplateCallBack,
			enabled:          true,
			Disable: func() {
				if _, ok := m.registeredResourceTemplates[*name]; !ok {
					fmt.Println("resource template not found")
					return
				}
				disabled := false
				m.registeredResourceTemplates[*name].Update(ResourceTemplateUpdates{Enabled: &disabled})
			},
			Enable: func() {
				if _, ok := m.registeredResourceTemplates[*name]; !ok {
					fmt.Println("resource not found")
					return
				}
				enabled := true
				m.registeredResourceTemplates[*name].Update(ResourceTemplateUpdates{Enabled: &enabled})
			},
			Remove: func() { delete(m.registeredResourceTemplates, *name) },
			Update: func(updates ResourceTemplateUpdates) {
				if _, ok := m.registeredResourceTemplates[*name]; !ok {
					fmt.Println("resource not found")
					return
				}
				if updates.Name != nil && *updates.Name != *name {
					resourceTemplateCopy := m.registeredResourceTemplates[*name]
					delete(m.registeredResourceTemplates, *name)
					m.registeredResourceTemplates[*updates.Name] = resourceTemplateCopy

					// 参照を更新することで、nameの値が変わった場合でも、以降の処理 & Disable/Enable/Removeが正しく動作するようにする
					name = updates.Name
				}
				if updates.Template != nil {
					m.registeredResourceTemplates[*name].resourceTemplate = updates.Template
				}
				if updates.Metadata != nil {
					m.registeredResourceTemplates[*name].metadata = updates.Metadata
				}
				if updates.Callback != nil {
					m.registeredResourceTemplates[*name].readCallback = *updates.Callback
				}
				if updates.Enabled != nil {
					m.registeredResourceTemplates[*name].enabled = *updates.Enabled
				}
				m.server.SendResourceListChanged()
			},
		}
		m.registeredResourceTemplates[*name] = &registeredResourceTemplate
		m.setResourceRequestHandlers()
		m.sendResourceListChanged()
		return nil, &registeredResourceTemplate, nil
	}
	return nil, nil, errors.New("unexpected error")
}
