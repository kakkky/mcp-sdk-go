package mcpserver

import (
	"net/url"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type RegisteredResources struct {
	name         string
	metadata     schema.ResourceMetadata
	readCallback readResourceCallback[schema.ResourceContentSchema]
	enabled      bool
	enable       func()
	disable      func()
	update       func(updates struct {
		name     string
		uri      string
		metadata schema.ResourceMetadata
		callback readResourceCallback[schema.ResourceContentSchema]
		enabled  bool
	})
}

type ReadResourceCallback[T schema.ResourceContentSchema] func(url url.URL) (schema.ReadResourceResultSchema[T], error)

type registered