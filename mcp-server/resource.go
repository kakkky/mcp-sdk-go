package mcpserver

import (
	"net/url"

	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type registeredResources struct {
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

type readResourceCallback[T schema.ResourceContentSchema] func(url url.URL) (schema.ReadResourceResultSchema[T], error)
