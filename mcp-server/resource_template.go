package mcpserver

import (
	"github.com/kakkky/mcp-sdk-go/shared/schema"
	utilities "github.com/kakkky/mcp-sdk-go/shared/utilities/uri-template"
)

type ResourceTemplate struct {
	uriTemp   *utilities.UriTemplate
	callBacks *ResourceTemplateCallbacks
}

type ResourceTemplateCallbacks struct {
	List     ListResourcesCallback
	Complete map[string]CompleteResourceCallback
}

// uriTemp（テンプレート）に具体的な値を埋め込んだURIを持つリソースをリストで返すコールバック
type ListResourcesCallback func() schema.ListResourcesResultSchema
type CompleteResourceCallback func(value string) []string

func NewResourceTemplate(uriTemplate string, callbacks *ResourceTemplateCallbacks) (*ResourceTemplate, error) {
	uriTemp, err := utilities.NewUriTemplate(uriTemplate)
	if err != nil {
		return nil, err
	}
	return &ResourceTemplate{
		uriTemp:   uriTemp,
		callBacks: callbacks,
	}, nil
}

func (r *ResourceTemplate) uriTemplate() *utilities.UriTemplate {
	return r.uriTemp
}

func (r *ResourceTemplate) ListCallback() ListResourcesCallback {
	if r.callBacks == nil {
		return nil
	}
	return r.callBacks.List
}

func (r *ResourceTemplate) CompleteCallBack(variable string) CompleteResourceCallback {
	if r.callBacks == nil {
		return nil
	}
	return r.callBacks.Complete[variable]
}
