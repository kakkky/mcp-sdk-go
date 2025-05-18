package mcpserver

import (
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kakkky/mcp-sdk-go/mcp-server/server"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
	utilities "github.com/kakkky/mcp-sdk-go/shared/utilities/uri-template"
)

func TestMcpServer_Resource(t *testing.T) {
	type args struct {
		name                         string
		uri                          string
		template                     *ResourceTemplate
		metadata                     *schema.ResourceMetadata
		readResourceCallBack         ReadResourceCallback[schema.ResourceContentSchema]
		readResourceTemplateCallBack ReadResourceTemplateCallback[schema.ResourceContentSchema]
	}
	tests := []struct {
		name                     string
		args                     args
		expectedResource         *RegisteredResource
		expectedResourceTemplate *RegisteredResourceTemplate
		isExpectedErr            bool
	}{
		{
			name: "normal : uri is provided",
			args: args{
				name:     "test",
				uri:      "file:///test.txt",
				template: nil,
				metadata: &schema.ResourceMetadata{
					Description: "test description",
					MimeType:    "text/plain",
				},
				readResourceCallBack: func(url url.URL) (schema.ReadResourceResultSchema[schema.ResourceContentSchema], error) {
					return schema.ReadResourceResultSchema[schema.ResourceContentSchema]{
						Contents: []schema.ResourceContentSchema{
							&schema.TextResourceContentsSchema{
								UriData:      "file:///test.txt",
								MimeTypeData: "test/plain",
								ContentData:  "test contents",
							},
						},
					}, nil
				},
				readResourceTemplateCallBack: nil,
			},
			expectedResource: &RegisteredResource{
				name: "test",
				metadata: &schema.ResourceMetadata{
					Description: "test description",
					MimeType:    "text/plain",
				},
				enabled: true,
			},
			isExpectedErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewMcpServer(schema.Implementation{}, &server.ServerOptions{
				Capabilities: schema.ServerCapabilities{
					Resources: &schema.Resources{
						ListChanged: true,
					},
				},
			})

			gotResource, gotResourceTemplate, err := sut.Resource(
				tt.args.name,
				tt.args.uri,
				tt.args.template,
				tt.args.metadata,
				tt.args.readResourceCallBack,
				tt.args.readResourceTemplateCallBack)

			if (err != nil) != tt.isExpectedErr {
				t.Errorf("McpServer.Resource() error = %v, wantErr %v", err, tt.isExpectedErr)
				return
			}

			var cmpOpts = cmp.Options{
				cmpopts.IgnoreFields(RegisteredResource{}, "readCallback", "Enable", "Disable", "Update", "Remove"),
				cmp.AllowUnexported(RegisteredResource{}),
				cmpopts.IgnoreFields(RegisteredResourceTemplate{}, "readCallback", "Enable", "Disable", "Update", "Remove"),
				cmp.AllowUnexported(RegisteredResourceTemplate{}),
			}

			if diff := cmp.Diff(gotResource, tt.expectedResource, cmpOpts); diff != "" {
				t.Errorf("McpServer.Resource() gotResource = %v, want %v, diff %s", gotResource, tt.expectedResource, diff)
			}
			if diff := cmp.Diff(gotResourceTemplate, tt.expectedResourceTemplate, cmpOpts); diff != "" {
				t.Errorf("McpServer.Resource() gotResourceTemplate = %v, want %v, diff %s", gotResourceTemplate, tt.expectedResourceTemplate, diff)
			}
			// リソースが登録されているか確認
			if tt.expectedResource != nil {
				if diff := cmp.Diff(sut.registeredResources[tt.args.uri], gotResource, cmpOpts); diff != "" {
					t.Errorf("McpServer.Resource() registeredResources = %v, want %v, diff %s", sut.registeredResources[tt.args.uri], tt.expectedResource, diff)
				}
			}
			// リソーステンプレートが登録されているか確認
			if tt.expectedResourceTemplate != nil {
				if diff := cmp.Diff(sut.registeredResourceTemplates[tt.args.name], gotResourceTemplate, cmpOpts); diff != "" {
					t.Errorf("McpServer.Resource() registeredResourceTemplates = %v, want %v, diff %s", sut.registeredResourceTemplates[tt.args.name], tt.expectedResourceTemplate, diff)
				}
			}

			// リソースを更新
			if gotResource != nil {
				updatedUri := "file://test2.txt"
				// URIを更新する
				gotResource.Update(ResourceUpdates{
					Uri: updatedUri,
				})
				// リソースが更新されていることを確認
				if diff := cmp.Diff(sut.registeredResources[updatedUri], gotResource, cmpOpts); diff != "" {
					t.Errorf("McpServer.Resource() registeredResources = %v, want %v, diff %s", sut.registeredResources["file://test2.txt"], gotResource, diff)
				}
				// リソースを無効化できる
				gotResource.Disable()
				if sut.registeredResources[updatedUri].enabled {
					t.Errorf("McpServer.Resource() registeredResources = %v, want %v", sut.registeredResources[updatedUri].enabled, false)
				}

				// リソースを削除できる
				gotResource.Remove()
				if _, ok := sut.registeredResources[updatedUri]; ok {
					t.Errorf("McpServer.Resource() registeredResources = %v, want %v", sut.registeredResources[updatedUri], nil)
				}
			}
			// リソーステンプレートを更新
			if gotResourceTemplate != nil {
				// nameを更新する
				updatedName := "updated"
				gotResourceTemplate.Update(ResourceTemplateUpdates{
					Name: updatedName,
				})
				// リソーステンプレートが更新されていることを確認
				if diff := cmp.Diff(sut.registeredResourceTemplates[updatedName], gotResourceTemplate, cmpOpts); diff != "" {
					t.Errorf("McpServer.Resource() registeredResourceTemplates = %v, want %v, diff %s", sut.registeredResourceTemplates["updated"], gotResourceTemplate, diff)
				}

				updatedUriTemp, _ := utilities.NewUriTemplate("file:///updated/{test}")
				updatedTemplate := &ResourceTemplate{
					uriTemp: updatedUriTemp,
				}
				// テンプレートを更新する
				gotResourceTemplate.Update(ResourceTemplateUpdates{
					Template: updatedTemplate,
				})
				// リソーステンプレートが更新されていることを確認
				if diff := cmp.Diff(sut.registeredResourceTemplates[updatedName], gotResourceTemplate, cmpOpts); diff != "" {
					t.Errorf("McpServer.Resource() registeredResourceTemplates = %v, want %v, diff %s", sut.registeredResourceTemplates[tt.args.name], gotResourceTemplate, diff)
				}
			}
		})
	}
}
