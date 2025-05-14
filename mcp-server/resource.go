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

type RegisteredResourceTemplate struct {
	resourceTemplate *ResourceTemplate
	metadata         schema.ResourceMetadata
	readCallback     ReadResourceCallback[schema.ResourceContentSchema]
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

type ResourceTemplate struct {
	uriTemp   *utilities.UriTemplate
	callBacks *struct {
		list     *ListResourcesCallback
		complete map[string]*CompleteResourceCallback
	}
}
type ListResourcesCallback func() schema.ListResourcesResultSchema
type CompleteResourceCallback func() []string

func NewResourceTemplate(uriTemplate string) (*ResourceTemplate, error) {
	uriTemp, err := utilities.NewUriTemplate(uriTemplate)
	if err != nil {
		return nil, err
	}
	return &ResourceTemplate{
		uriTemp: uriTemp,
	}, nil
}

func (r *ResourceTemplate) uriTemplate() *utilities.UriTemplate {
	return r.uriTemp
}

func (r *ResourceTemplate) ListCallback() *ListResourcesCallback {
	if r.callBacks == nil {
		return nil
	}
	return r.callBacks.list
}

func (r *ResourceTemplate) CompleteCallBack(variable string) *CompleteResourceCallback {
	if r.callBacks == nil {
		return nil
	}
	return r.callBacks.complete[variable]
}
