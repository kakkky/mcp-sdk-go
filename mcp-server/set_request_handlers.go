package mcpserver

import (
	"fmt"
	"net/url"

	mcperr "github.com/kakkky/mcp-sdk-go/shared/mcp-err"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func (m *McpServer) setResourceRequestHandlers() error {
	if m.isResourceHandlersInitialized {
		return nil
	}
	// リクエストハンドラは重複して登録できない
	var resourceMethodlist = []string{"resources/list", "resources/templates/list", "resources/read"}
	for _, method := range resourceMethodlist {
		if err := m.Server.ValidateCanSetRequestHandler(method); err != nil {
			return err
		}
	}
	_ = m.Server.RegisterCapabilities(schema.ServerCapabilities{
		Resources: &schema.Resources{
			ListChanged: true,
		},
	})
	m.Server.SetRequestHandler(&schema.ListResourceRequestSchema{MethodName: "resources/list"}, func(req schema.JsonRpcRequest) (schema.Result, error) {
		var resources []schema.ResourceSchema
		for uri, registerdResource := range m.registeredResources {
			if registerdResource.enabled {
				resources = append(resources, schema.ResourceSchema{
					Uri:              uri,
					Name:             registerdResource.name,
					ResourceMetadata: registerdResource.metadata,
				})
			}
		}

		var templateResources []schema.ResourceSchema
		for _, registerdResourceTemplate := range m.registeredResourceTemplates {
			if registerdResourceTemplate.resourceTemplate.ListCallback() == nil {
				continue
			}
			result := registerdResourceTemplate.resourceTemplate.ListCallback()()
			for _, resource := range result.Resources {
				resource = schema.ResourceSchema{
					Uri:              resource.Uri,
					Name:             resource.Name,
					ResourceMetadata: registerdResourceTemplate.metadata,
				}
				templateResources = append(templateResources, resource)
			}

		}
		resources = append(resources, templateResources...)
		return &schema.ListResourcesResultSchema{
			Resources: resources,
		}, nil
	})

	m.Server.SetRequestHandler(&schema.ListResourceTemplatesRequestSchema{MethodName: "resources/templates/list"}, func(req schema.JsonRpcRequest) (schema.Result, error) {
		var resourceTemplates []schema.ResourceTemplateSchema
		for name, registerdResourceTemplate := range m.registeredResourceTemplates {
			resourceTemplate := schema.ResourceTemplateSchema{
				Name:             name,
				UriTemplate:      registerdResourceTemplate.resourceTemplate.uriTemp.ToString(),
				ResourceMetadata: registerdResourceTemplate.metadata,
			}
			resourceTemplates = append(resourceTemplates, resourceTemplate)
		}
		return &schema.ListResourceTemplatesResultSchema{
			ResourceTemplates: resourceTemplates,
		}, nil
	})

	m.Server.SetRequestHandler(&schema.ReadResourceRequestSchema{MethodName: "resources/read"}, func(req schema.JsonRpcRequest) (schema.Result, error) {
		request, ok := req.Request.(*schema.ReadResourceRequestSchema)
		if !ok {
			return nil, mcperr.NewMcpErr(mcperr.INVALID_REQUEST, "invalid request", nil)
		}
		uri, err := url.Parse(request.ParamsData.Uri)
		if err != nil {
			return nil, mcperr.NewMcpErr(mcperr.INVALID_PARAMS, fmt.Sprintf("invalid uri %s", uri.String()), nil)
		}
		resource, ok := m.registeredResources[uri.String()]

		// paramsのuriからリソースを取得できなかった場合、リソーステンプレートを確認する
		if !ok {
			for _, registerdResourceTemplate := range m.registeredResourceTemplates {
				variables, err := registerdResourceTemplate.resourceTemplate.uriTemp.Match(uri.String())
				if err != nil {
					return nil, mcperr.NewMcpErr(mcperr.INVALID_PARAMS, fmt.Sprintf("invalid uri template %s", request.ParamsData.Uri), nil)
				}
				if variables != nil {
					result, err := registerdResourceTemplate.readCallback(*uri, variables)
					if err != nil {
						return nil, mcperr.NewMcpErr(mcperr.INTERNAL_ERROR, fmt.Sprintf("failed to read resource %s", uri.String()), err)
					}
					return &result, nil
				}
			}
			return nil, mcperr.NewMcpErr(mcperr.INVALID_PARAMS, fmt.Sprintf("resource %s not found", uri.String()), nil)
		}

		if !resource.enabled {
			return nil, mcperr.NewMcpErr(mcperr.INVALID_PARAMS, fmt.Sprintf("resource %s disabled", uri.String()), nil)
		}
		result, err := resource.readCallback(*uri)
		if err != nil {
			return nil, mcperr.NewMcpErr(mcperr.INTERNAL_ERROR, fmt.Sprintf("failed to read resource %s", uri.String()), err)
		}
		return &result, nil
	})

	_ = m.setCompletionRequestHandlers()

	m.isResourceHandlersInitialized = true
	return nil
}

func (m *McpServer) setCompletionRequestHandlers() error {
	if m.isCompletionHandlersInitialized {
		return nil
	}
	if err := m.Server.ValidateCanSetRequestHandler("completion/complete"); err != nil {
		return err
	}
	m.Server.SetRequestHandler(&schema.CompleteRequestSchema{MethodName: "completion/complete"}, func(req schema.JsonRpcRequest) (schema.Result, error) {
		request, ok := req.Request.(*schema.CompleteRequestSchema)
		if !ok {
			return nil, mcperr.NewMcpErr(mcperr.INVALID_REQUEST, "invalid request", nil)
		}
		params := request.Params().(schema.CompleteRequestParams)
		switch params.Ref.Type() {
		case "ref/prompt":
			// TODO: implement
		case "ref/resource":
			ref, ok := params.Ref.(*schema.ResourceReferenceSchema)
			if !ok {
				return nil, mcperr.NewMcpErr(mcperr.INVALID_PARAMS, "invalid ref type", nil)
			}
			return m.handleResourceCompletion(*request, *ref)
		default:
			return nil, mcperr.NewMcpErr(mcperr.INVALID_PARAMS, fmt.Sprintf("invalid completion reference : %s", params.Ref), nil)
		}
		return nil, nil
	})
	m.isCompletionHandlersInitialized = true
	return nil
}
