package mcpserver

import (
	"net/url"

	utilities "github.com/kakkky/mcp-sdk-go/mcp-server/utilities/uri-template"
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

type registered