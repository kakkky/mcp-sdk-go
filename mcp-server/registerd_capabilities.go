package mcpserver

import (
	"net/url"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type RegisteredResource struct {
	name         string
	metadata     *schema.ResourceMetadata
	readCallback ReadResourceCallback[schema.ResourceContentSchema]
	enabled      bool
	Enable       func()
	Disable      func()
	Update       func(ResourceUpdates)
	Remove       func()
}
type ResourceUpdates struct {
	Name     string
	Uri      string
	Metadata *schema.ResourceMetadata
	Callback *ReadResourceCallback[schema.ResourceContentSchema]
	Enabled  *bool
}

type ReadResourceCallback[T schema.ResourceContentSchema] func(url url.URL) (schema.ReadResourceResultSchema, error)

type RegisteredResourceTemplate struct {
	resourceTemplate *ResourceTemplate
	metadata         *schema.ResourceMetadata
	readCallback     ReadResourceTemplateCallback[schema.ResourceContentSchema]
	enabled          bool
	Enable           func()
	Disable          func()
	Update           func(ResourceTemplateUpdates)
	Remove           func()
}

type ResourceTemplateUpdates struct {
	Name     string
	Template *ResourceTemplate
	Metadata *schema.ResourceMetadata
	Callback *ReadResourceTemplateCallback[schema.ResourceContentSchema]
	Enabled  *bool
}

type ReadResourceTemplateCallback[T schema.ResourceContentSchema] func(url url.URL, variables map[string]any) (schema.ReadResourceResultSchema, error)
