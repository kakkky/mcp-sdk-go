package mcpserver

import (
	"github.com/kakkky/mcp-sdk-go/shared/schema"
	utilities "github.com/kakkky/mcp-sdk-go/shared/utilities/uri-template"
)

type ResourceTemplate struct {
	uriTemp   *utilities.UriTemplate
	callBacks *struct {
		list     ListResourcesCallback
		complete map[string]CompleteResourceCallback
	}
}

// uriTemp（テンプレート）に具体的な値を埋め込んだURIを持つリソースをリストで返すコールバック
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

func (r *ResourceTemplate) ListCallback() ListResourcesCallback {
	if r.callBacks == nil {
		return nil
	}
	return r.callBacks.list
}

func (r *ResourceTemplate) CompleteCallBack(variable string) CompleteResourceCallback {
	if r.callBacks == nil {
		return nil
	}
	return r.callBacks.complete[variable]
}
