package mcpserver

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kakkky/mcp-sdk-go/mcp-server/server"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
	utilities "github.com/kakkky/mcp-sdk-go/shared/utilities/uri-template"
	"github.com/stretchr/testify/assert"
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
		// URIリソースの正常登録
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
		// テンプレートリソースの正常登録
		{
			name: "normal : template is provided",
			args: args{
				name: "test-template",
				uri:  "",
				template: func() *ResourceTemplate {
					template, _ := NewResourceTemplate("/api/users/{userId}")
					return template
				}(),
				metadata: &schema.ResourceMetadata{
					Description: "User resource template",
					MimeType:    "application/json",
				},
				readResourceCallBack: nil,
				readResourceTemplateCallBack: func(url url.URL, vars map[string]any) (schema.ReadResourceResultSchema[schema.ResourceContentSchema], error) {
					userId, _ := vars["userId"].(string)
					return schema.ReadResourceResultSchema[schema.ResourceContentSchema]{
						Contents: []schema.ResourceContentSchema{
							&schema.TextResourceContentsSchema{
								UriData:      url.String(),
								MimeTypeData: "application/json",
								ContentData:  `{"userId":"` + userId + `", "name":"Test User"}`,
							},
						},
					}, nil
				},
			},
			expectedResourceTemplate: &RegisteredResourceTemplate{
				resourceTemplate: func() *ResourceTemplate {
					template, _ := NewResourceTemplate("/api/users/{userId}")
					return template
				}(),
				metadata: &schema.ResourceMetadata{
					Description: "User resource template",
					MimeType:    "application/json",
				},
				enabled: true,
			},
			isExpectedErr: false,
		},

		// URIとテンプレートの両方を指定した場合（エラー）
		{
			name: "error : both uri and template are provided",
			args: args{
				name: "invalid-resource",
				uri:  "file://data.txt",
				template: func() *ResourceTemplate {
					template, _ := NewResourceTemplate("/api/data/{id}")
					return template
				}(),
				metadata: &schema.ResourceMetadata{
					Description: "Invalid resource",
					MimeType:    "text/plain",
				},
				readResourceCallBack: func(url url.URL) (schema.ReadResourceResultSchema[schema.ResourceContentSchema], error) {
					return schema.ReadResourceResultSchema[schema.ResourceContentSchema]{}, nil
				},
				readResourceTemplateCallBack: func(url url.URL, vars map[string]any) (schema.ReadResourceResultSchema[schema.ResourceContentSchema], error) {
					return schema.ReadResourceResultSchema[schema.ResourceContentSchema]{}, nil
				},
			},
			expectedResource:         nil,
			expectedResourceTemplate: nil,
			isExpectedErr:            true,
		},

		// URIとテンプレートの両方を省略した場合（エラー）
		{
			name: "error : neither uri nor template is provided",
			args: args{
				name:     "missing-resource",
				uri:      "",
				template: nil,
				metadata: &schema.ResourceMetadata{
					Description: "Missing resource",
					MimeType:    "text/plain",
				},
				readResourceCallBack:         nil,
				readResourceTemplateCallBack: nil,
			},
			expectedResource:         nil,
			expectedResourceTemplate: nil,
			isExpectedErr:            true,
		},

		// URIリソース登録時にコールバック未指定（エラー）
		{
			name: "error : uri without callback",
			args: args{
				name:     "missing-callback",
				uri:      "file://no-callback.txt",
				template: nil,
				metadata: &schema.ResourceMetadata{
					Description: "Missing callback",
					MimeType:    "text/plain",
				},
				readResourceCallBack:         nil, // コールバック未指定
				readResourceTemplateCallBack: nil,
			},
			expectedResource:         nil,
			expectedResourceTemplate: nil,
			isExpectedErr:            true,
		},

		// テンプレートリソース登録時にコールバック未指定（エラー）
		{
			name: "error : template without callback",
			args: args{
				name: "template-missing-callback",
				uri:  "",
				template: func() *ResourceTemplate {
					template, _ := NewResourceTemplate("/api/data/{id}")
					return template
				}(),
				metadata: &schema.ResourceMetadata{
					Description: "Template missing callback",
					MimeType:    "application/json",
				},
				readResourceCallBack:         nil,
				readResourceTemplateCallBack: nil, // コールバック未指定
			},
			expectedResource:         nil,
			expectedResourceTemplate: nil,
			isExpectedErr:            true,
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

			cmpOpts := cmp.Options{
				cmpopts.IgnoreFields(RegisteredResource{}, "readCallback", "Enable", "Disable", "Update", "Remove"),
				cmp.AllowUnexported(RegisteredResource{}),
				cmpopts.IgnoreFields(RegisteredResourceTemplate{}, "readCallback", "Enable", "Disable", "Update", "Remove"),
				cmp.AllowUnexported(RegisteredResourceTemplate{}, ResourceTemplate{}),
				cmp.Comparer(func(a, b utilities.UriTemplate) bool {
					return a.ToString() == b.ToString()
				}),
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

			// リソースを更新できる
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
			// リソーステンプレートを更新できる
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
			// リクエストハンドラが登録されていることを確認
			if gotResource != nil || gotResourceTemplate != nil {
				protocol := reflect.ValueOf(sut.server.Protocol).Elem()
				handlers := protocol.FieldByName("handlers").Elem()
				requestHandlers := handlers.FieldByName("requestHandlers")
				// リクエストハンドラーとして登録されているメソッド一覧を取得
				var methods []string
				for _, method := range requestHandlers.MapKeys() {
					methods = append(methods, method.String())
				}
				// 想定されるメソッド一覧と比較
				basicMethods := []string{"ping", "initialize"}
				expectedMethods := append([]string{"resources/list", "resources/templates/list", "resources/read", "completion/complete"}, basicMethods...)
				if len(methods) != len(expectedMethods) {
					t.Errorf("McpServer.Resource() requestHandlers = %v, want %v", len(methods), len(expectedMethods))
				}
				if !assert.ElementsMatch(t, methods, expectedMethods) {
					t.Errorf("McpServer.Resource() requestHandlers = %v, want %v", methods, expectedMethods)
				}
			}
		})
	}
}
