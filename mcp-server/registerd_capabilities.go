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

type RegisteredTool struct {
	description    string
	propertySchema schema.PropertySchema
	annotations    *schema.ToolAnotationsSchema
	callback       ToolCallback
	enabled        bool
	Enable         func()
	Remove         func()
	Disable        func()
	Update         func(ToolUpdates)
}

type ToolUpdates struct {
	Name         string
	Description  string
	ParamsSchema schema.PropertySchema
	callback     ToolCallback
	Annotations  *schema.ToolAnotationsSchema
	Enabled      *bool
}

type ToolCallback func(args map[string]any) (schema.CallToolResultSchema, error)

type RegisteredPrompt struct {
	description string
	argsSchema  []schema.PromptAugmentSchema
	callback    PromptCallback
	enabled     bool
	Enable      func()
	Disable     func()
	Remove      func()
	Update      func(PromptUpdates)
}

type PromptUpdates struct {
	Name        string
	Description string
	ArgsSchema  []schema.PromptAugmentSchema
	Callback    PromptCallback
	Enabled     *bool
}

type PromptCallback func(args []schema.PromptAugmentSchema) (schema.GetPromptResultSchema, error)
