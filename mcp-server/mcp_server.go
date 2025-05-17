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
	name string,
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
		if m.registeredResources[*uri] == nil {
			return nil, nil, fmt.Errorf("resource %s is already registered", *uri)
		}

		registeredResource := RegisteredResource{
			name:         name,
			metadata:     metadata,
			readCallback: readResourceCallBack,
			enabled:      true,
			disable:      func() { m.registeredResources[*uri].update(resourceUpdates{enabled: false}) },
			enable:       func() { m.registeredResources[*uri].update(resourceUpdates{enabled: true}) },
			remove:       func() { m.registeredResources[*uri].update(resourceUpdates{uri: nil}) },
			update: func(updates resourceUpdates) {
				// uriが更新されても以降の処理を正しく行えるようにするため
				currentUri := uri

				if updates.uri != nil && *updates.uri != *uri {
					resourceCopy := m.registeredResources[*uri]
					delete(m.registeredResources, *uri)
					m.registeredResources[*updates.uri] = resourceCopy
					// 以降で参照するuriを更新する
					currentUri = updates.uri
				}
				if updates.name != nil {
					m.registeredResources[*currentUri].name = *updates.name
				}
				if updates.metadata != nil {
					m.registeredResources[*currentUri].metadata = updates.metadata
				}
				if updates.callback != nil {
					m.registeredResources[*currentUri].readCallback = *updates.callback
				}
				if updates.enabled {
					m.registeredResources[*currentUri].enabled = updates.enabled
				}
				m.server.SendResourceListChanged()
			},
		}
		m.registeredResources[*uri] = &registeredResource
		m.setResourceRequestHandlers()
		m.server.SendResourceListChanged()
		return &registeredResource, nil, nil
	}
	// templateが渡された場合
	if template != nil {
		if m.registeredResourceTemplates[name] != nil {
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
			disable:          func() { m.registeredResourceTemplates[name].update(resourceTemplateUpdates{enabled: false}) },
			enable:           func() { m.registeredResourceTemplates[name].update(resourceTemplateUpdates{enabled: true}) },
			remove:           func() { m.registeredResourceTemplates[name].update(resourceTemplateUpdates{name: nil}) },
			update: func(updates resourceTemplateUpdates) {
				// nameが更新されても以降の処理を正しく行えるようにするため
				currentName := &name

				if updates.name != nil && *updates.name != name {
					resourceTemplateCopy := m.registeredResourceTemplates[name]
					delete(m.registeredResourceTemplates, name)
					m.registeredResourceTemplates[*updates.name] = resourceTemplateCopy
					// 以降で参照するnameを更新する
					currentName = updates.name
				}
				if updates.template != nil {
					m.registeredResourceTemplates[*currentName].resourceTemplate = updates.template
				}
				if updates.metadata != nil {
					m.registeredResourceTemplates[*currentName].metadata = updates.metadata
				}
				if updates.callback != nil {
					m.registeredResourceTemplates[*currentName].readCallback = *updates.callback
				}
				if updates.enabled {
					m.registeredResourceTemplates[*currentName].enabled = updates.enabled
				}
				m.server.SendResourceListChanged()
			},
		}
		m.registeredResourceTemplates[name] = &registeredResourceTemplate
		m.setResourceRequestHandlers()
		m.sendResourceListChanged()
		return nil, &registeredResourceTemplate, nil
	}
	return nil, nil, errors.New("unexpected error")
}
