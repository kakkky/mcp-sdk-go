package mcpserver

import (
	"net/url"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type RegisteredResource struct {
	name         string
	metadata     schema.ResourceMetadata
	readCallback ReadResourceCallback[schema.ResourceContentSchema]
	enabled      bool
	enable       func()
	disable      func()
	update       func(updates struct {
		name     string
		uri      string
		metadata *schema.ResourceMetadata
		callback ReadResourceCallback[schema.ResourceContentSchema]
		enabled  bool
	})
	remove func()
}

type ReadResourceCallback[T schema.ResourceContentSchema] func(url url.URL) (schema.ReadResourceResultSchema[T], error)

type RegisteredResourceTemplate struct {
	resourceTemplate *ResourceTemplate
	metadata         schema.ResourceMetadata
	readCallback     ReadResourceTemplateCallback[schema.ResourceContentSchema]
	enabled          bool
	enable           func()
	disable          func()
	update           func(updates struct {
		name     string
		template *ResourceTemplate
		metadata *schema.ResourceMetadata
		callback ReadResourceTemplateCallback[schema.ResourceContentSchema]
		enabled  bool
	})
	remove func()
}

type ReadResourceTemplateCallback[T schema.ResourceContentSchema] func(url url.URL, variables map[string]any) schema.ReadResourceResultSchema[T]
