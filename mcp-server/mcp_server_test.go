package mcpserver

import (
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kakkky/mcp-sdk-go/mcp-server/server"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
	utilities "github.com/kakkky/mcp-sdk-go/shared/utilities/uri-template"
	"github.com/stretchr/testify/assert"
)

// URIベースのリソース登録テスト
func TestMcpServer_Resource(t *testing.T) {
	type args struct {
		name                 string
		uri                  string
		metadata             *schema.ResourceMetadata
		readResourceCallBack ReadResourceCallback[schema.ResourceContentSchema]
	}
	tests := []struct {
		name             string
		args             args
		expectedResource *RegisteredResource
		isExpectedErr    bool
	}{
		// URIリソースの正常登録
		{
			name: "normal : uri is provided",
			args: args{
				name: "test",
				uri:  "file:///test.txt",
				metadata: &schema.ResourceMetadata{
					Description: "test description",
					MimeType:    "text/plain",
				},
				readResourceCallBack: func(url url.URL) (schema.ReadResourceResultSchema, error) {
					return schema.ReadResourceResultSchema{
						Contents: []schema.ResourceContentSchema{
							&schema.TextResourceContentsSchema{
								UriData:      "file:///test.txt",
								MimeTypeData: "test/plain",
								ContentData:  "test contents",
							},
						},
					}, nil
				},
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
		// URIリソース登録時にコールバック未指定（エラー）
		{
			name: "error : uri without callback",
			args: args{
				name: "missing-callback",
				uri:  "file://no-callback.txt",
				metadata: &schema.ResourceMetadata{
					Description: "Missing callback",
					MimeType:    "text/plain",
				},
				readResourceCallBack: nil, // コールバック未指定
			},
			expectedResource: nil,
			isExpectedErr:    true,
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

			gotResource, err := sut.Resource(
				tt.args.name,
				tt.args.uri,
				tt.args.metadata,
				tt.args.readResourceCallBack)

			if (err != nil) != tt.isExpectedErr {
				t.Errorf("McpServer.Resource() error = %v, wantErr %v", err, tt.isExpectedErr)
				return
			}

			cmpOpts := cmp.Options{
				cmpopts.IgnoreFields(RegisteredResource{}, "readCallback", "Enable", "Disable", "Update", "Remove"),
				cmp.AllowUnexported(RegisteredResource{}),
			}

			if diff := cmp.Diff(gotResource, tt.expectedResource, cmpOpts); diff != "" {
				t.Errorf("McpServer.Resource() gotResource = %v, want %v, diff %s", gotResource, tt.expectedResource, diff)
			}

			// リソースが登録されているか確認
			if tt.expectedResource != nil {
				if diff := cmp.Diff(sut.registeredResources[tt.args.uri], gotResource, cmpOpts); diff != "" {
					t.Errorf("McpServer.Resource() registeredResources = %v, want %v, diff %s", sut.registeredResources[tt.args.uri], tt.expectedResource, diff)
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

			// リクエストハンドラが登録されていることを確認
			if gotResource != nil {
				protocol := reflect.ValueOf(sut.Server.Protocol).Elem()
				handlers := protocol.FieldByName("handlers").Elem()
				requestHandlers := handlers.FieldByName("requestHandlers")
				// リクエストハンドラーとして登録されているメソッド一覧を取得
				var methods []string
				for _, method := range requestHandlers.MapKeys() {
					methods = append(methods, method.String())
				}
				// 想定されるメソッド一覧と比較
				basicMethods := []string{"ping", "initialize"}
				expectedMethods := append([]string{"resources/list", "resources/read", "completion/complete", "resources/templates/list"}, basicMethods...)
				if !assert.ElementsMatch(t, methods, expectedMethods) {
					t.Errorf("McpServer.Resource() requestHandlers = %v, want %v", methods, expectedMethods)
				}
			}
		})
	}
}

// テンプレートベースのリソース登録テスト
func TestMcpServer_ResourceTemplate(t *testing.T) {
	type args struct {
		name                         string
		template                     *ResourceTemplate
		metadata                     *schema.ResourceMetadata
		readResourceTemplateCallBack ReadResourceTemplateCallback[schema.ResourceContentSchema]
	}
	tests := []struct {
		name                     string
		args                     args
		expectedResourceTemplate *RegisteredResourceTemplate
		isExpectedErr            bool
	}{
		// テンプレートリソースの正常登録
		{
			name: "normal : template is provided",
			args: args{
				name: "test-template",
				template: func() *ResourceTemplate {
					template, _ := NewResourceTemplate(
						"/api/users/{userId}",
						&ResourceTemplateCallbacks{
							List: func() schema.ListResourcesResultSchema {
								return schema.ListResourcesResultSchema{
									Resources: []schema.ResourceSchema{
										{
											Name: "test-template-123",
											Uri:  "/api/users/123",
										},
										{
											Name: "test-template-456",
											Uri:  "/api/users/456",
										},
										{
											Name: "test-template-789",
											Uri:  "/api/users/789",
										},
									},
								}
							},
							Complete: map[string]CompleteResourceCallback{
								"userId": func(value string) []string {
									return []string{"123", "456", "789"}
								},
							},
						},
					)
					return template
				}(),
				metadata: &schema.ResourceMetadata{
					Description: "User resource template",
					MimeType:    "application/json",
				},
				readResourceTemplateCallBack: func(url url.URL, vars map[string]any) (schema.ReadResourceResultSchema, error) {
					userId, _ := vars["userId"].(string)
					return schema.ReadResourceResultSchema{
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
					template, _ := NewResourceTemplate(
						"/api/users/{userId}",
						&ResourceTemplateCallbacks{
							List: func() schema.ListResourcesResultSchema {
								return schema.ListResourcesResultSchema{
									Resources: []schema.ResourceSchema{
										{
											Name: "test-template-123",
											Uri:  "/api/users/123",
										},
										{
											Name: "test-template-456",
											Uri:  "/api/users/456",
										},
										{
											Name: "test-template-789",
											Uri:  "/api/users/789",
										},
									},
								}
							},
							Complete: map[string]CompleteResourceCallback{
								"userId": func(value string) []string {
									return []string{"123", "456", "789"}
								},
							},
						},
					)
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
		// テンプレートリソース登録時にコールバック未指定（エラー）
		{
			name: "error : template without callback",
			args: args{
				name: "template-missing-callback",
				template: func() *ResourceTemplate {
					template, _ := NewResourceTemplate("/api/data/{id}", nil)
					return template
				}(),
				metadata: &schema.ResourceMetadata{
					Description: "Template missing callback",
					MimeType:    "application/json",
				},
				readResourceTemplateCallBack: nil, // コールバック未指定
			},
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

			gotResourceTemplate, err := sut.ResourceTemplate(
				tt.args.name,
				tt.args.template,
				tt.args.metadata,
				tt.args.readResourceTemplateCallBack)

			if (err != nil) != tt.isExpectedErr {
				t.Errorf("McpServer.ResourceTemplate() error = %v, wantErr %v", err, tt.isExpectedErr)
				return
			}

			cmpOpts := cmp.Options{
				cmpopts.IgnoreFields(RegisteredResourceTemplate{}, "readCallback", "Enable", "Disable", "Update", "Remove"),
				cmp.AllowUnexported(RegisteredResourceTemplate{}, ResourceTemplate{}),
				cmpopts.IgnoreFields(ResourceTemplateCallbacks{}, "List", "Complete"), //TODO: コールバックの出力比較
				cmp.Comparer(func(a, b utilities.UriTemplate) bool {
					return a.ToString() == b.ToString()
				}),
			}

			if diff := cmp.Diff(gotResourceTemplate, tt.expectedResourceTemplate, cmpOpts); diff != "" {
				t.Errorf("McpServer.ResourceTemplate() gotResourceTemplate = %v, want %v, diff %s", gotResourceTemplate, tt.expectedResourceTemplate, diff)
			}

			// リソーステンプレートが登録されているか確認
			if tt.expectedResourceTemplate != nil {
				if diff := cmp.Diff(sut.registeredResourceTemplates[tt.args.name], gotResourceTemplate, cmpOpts); diff != "" {
					t.Errorf("McpServer.ResourceTemplate() registeredResourceTemplates = %v, want %v, diff %s", sut.registeredResourceTemplates[tt.args.name], tt.expectedResourceTemplate, diff)
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
					t.Errorf("McpServer.ResourceTemplate() registeredResourceTemplates = %v, want %v, diff %s", sut.registeredResourceTemplates["updated"], gotResourceTemplate, diff)
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
					t.Errorf("McpServer.ResourceTemplate() registeredResourceTemplates = %v, want %v, diff %s", sut.registeredResourceTemplates[tt.args.name], gotResourceTemplate, diff)
				}
			}

			// リクエストハンドラが登録されていることを確認
			if gotResourceTemplate != nil {
				protocol := reflect.ValueOf(sut.Server.Protocol).Elem()
				handlers := protocol.FieldByName("handlers").Elem()
				requestHandlers := handlers.FieldByName("requestHandlers")
				// リクエストハンドラーとして登録されているメソッド一覧を取得
				var methods []string
				for _, method := range requestHandlers.MapKeys() {
					methods = append(methods, method.String())
				}
				// 想定されるメソッド一覧と比較
				basicMethods := []string{"ping", "initialize"}
				expectedMethods := append([]string{"resources/templates/list", "resources/read", "completion/complete", "resources/list"}, basicMethods...)
				if !assert.ElementsMatch(t, methods, expectedMethods) {
					t.Errorf("McpServer.ResourceTemplate() requestHandlers = %v, want %v", methods, expectedMethods)
				}
			}
		})
	}
}

func TestMcpServer_Tool(t *testing.T) {
	type args struct {
		name           string
		description    string
		propertySchema schema.PropertySchema
		annotations    *schema.ToolAnotationsSchema
		callback       ToolCallback
	}
	tests := []struct {
		name          string
		args          args
		expectedTool  *RegisteredTool
		isExpectedErr bool
	}{
		// ツールの正常登録
		{
			name: "normal : tool is provided",
			args: args{
				name:        "test-tool",
				description: "Test tool description",
				propertySchema: schema.PropertySchema{
					"param1": schema.PropertyInfoSchema{
						Type:        "number",
						Description: "Test parameter",
					},
					"param2": schema.PropertyInfoSchema{
						Type:        "number",
						Description: "Another parameter",
					},
				},
				annotations: &schema.ToolAnotationsSchema{
					Title: "Test Tool",
				},
				callback: func(args map[string]any) (schema.CallToolResultSchema, error) {
					params1 := args["param1"].(int)
					params2 := args["param2"].(int)

					// ここでは単純に結果を返すだけのコールバック
					return schema.CallToolResultSchema{
						Content: []schema.ToolContentSchema{
							&schema.TextContentSchema{
								Type: strconv.Itoa(params1 + params2), // intをstringに変換
								Text: "result",
							},
						},
					}, nil
				},
			},
			expectedTool: &RegisteredTool{
				description: "Test tool description",
				propertySchema: schema.PropertySchema{
					"param1": schema.PropertyInfoSchema{
						Type:        "number",
						Description: "Test parameter",
					},
					"param2": schema.PropertyInfoSchema{
						Type:        "number",
						Description: "Another parameter",
					},
				},
				annotations: &schema.ToolAnotationsSchema{
					Title: "Test Tool",
				},
				callback: func(args map[string]any) (schema.CallToolResultSchema, error) {
					params1 := args["param1"].(int)
					params2 := args["param2"].(int)

					// ここでは単純に結果を返すだけのコールバック
					return schema.CallToolResultSchema{
						Content: []schema.ToolContentSchema{
							&schema.TextContentSchema{
								Type: strconv.Itoa(params1 + params2), // intをstringに変換
								Text: "result",
							},
						},
					}, nil
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
					Tools: &schema.Tools{
						ListChanged: true,
					},
				},
			})

			gotTool, err := sut.Tool(
				tt.args.name,
				tt.args.description,
				tt.args.propertySchema,
				tt.args.annotations,
				tt.args.callback)

			if (err != nil) != tt.isExpectedErr {
				t.Errorf("McpServer.Tool() error = %v, wantErr %v", err, tt.isExpectedErr)
				return
			}

			cmpOpts := cmp.Options{
				cmpopts.IgnoreFields(RegisteredTool{}, "Enable", "Disable", "Update", "Remove", "callback"),
				cmp.AllowUnexported(RegisteredTool{}),
			}

			if diff := cmp.Diff(gotTool, tt.expectedTool, cmpOpts); diff != "" {
				t.Errorf("McpServer.Tool() gotTool = %v, want %v, diff %s", gotTool, tt.expectedTool, diff)
			}

			// ツールが登録されているか確認
			if tt.expectedTool != nil {
				if diff := cmp.Diff(sut.registerdTools[tt.args.name], gotTool, cmpOpts); diff != "" {
					t.Errorf("McpServer.Tool() registerdTools = %v, want %v, diff %s", sut.registerdTools[tt.args.name], tt.expectedTool, diff)
				}
			}

			// ツールを更新できる
			if gotTool != nil {
				// nameを更新する
				updatedName := "updated-tool"
				gotTool.Update(ToolUpdates{
					Name: updatedName,
				})
				// ツールが更新されていることを確認
				if diff := cmp.Diff(sut.registerdTools[updatedName], gotTool, cmpOpts); diff != "" {
					t.Errorf("McpServer.Tool() registerdTools = %v, want %v, diff %s", sut.registerdTools[updatedName], gotTool, diff)
				}

				// ツールを無効化できる
				gotTool.Disable()
				if sut.registerdTools[updatedName].enabled {
					t.Errorf("McpServer.Tool() registerdTools = %v, want %v", sut.registerdTools[updatedName].enabled, false)
				}

				// ツールを削除できる
				gotTool.Remove()
				if _, ok := sut.registerdTools[updatedName]; ok {
					t.Errorf("McpServer.Tool() registerdTools = %v, want %v", sut.registerdTools[updatedName], nil)
				}
			}
		})
	}
}
