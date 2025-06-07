package mcpserver

import (
	"errors"
	"fmt"

	"github.com/kakkky/mcp-sdk-go/mcp-server/server"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

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

// uriかtemplateのどちらかに値を渡す
// uriを渡す場合は、readResourceCallBackを渡す
// templateを渡す場合は、readResourceTemplateCallBackを渡す
func (m *McpServer) Resource(
	name string,
	uri string,
	template *ResourceTemplate,
	metadata *schema.ResourceMetadata,
	readResourceCallBack ReadResourceCallback[schema.ResourceContentSchema], // uriを渡す場合
	readResourceTemplateCallBack ReadResourceTemplateCallback[schema.ResourceContentSchema], // templateを渡す場合
) (*RegisteredResource, *RegisteredResourceTemplate, error) {
	if (uri == "" && template == nil) || (uri != "" && template != nil) {
		return nil, nil, errors.New("please provide a value for either uri or template")
	}
	// uriが渡された場合
	if uri != "" {
		uriPtr := &uri
		if readResourceCallBack == nil {
			return nil, nil, errors.New("readResourceCallBack is required when uri is provided")
		}
		if m.registeredResources[*uriPtr] != nil {
			return nil, nil, fmt.Errorf("resource %s is already registered", uri)
		}

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
				m.Server.SendResourceListChanged()
			},
		}
		m.registeredResources[*uriPtr] = &registeredResource
		m.setResourceRequestHandlers()
		m.sendResourceListChanged()
		return &registeredResource, nil, nil
	}
	// templateが渡された場合
	if template != nil {
		namePtr := &name
		if m.registeredResourceTemplates[*namePtr] != nil {
			return nil, nil, fmt.Errorf("resource template %s is already registered", name)
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
				m.Server.SendResourceListChanged()
			},
		}
		m.registeredResourceTemplates[*namePtr] = &registeredResourceTemplate
		m.setResourceRequestHandlers()
		m.sendResourceListChanged()
		return nil, &registeredResourceTemplate, nil
	}
	return nil, nil, errors.New("unexpected error")
}
