package mcpserver

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func TestResourceTemplate_uriTemplate(t *testing.T) {
	type args struct {
		uriTemplate string
		callbacks   *ResourceTemplateCallbacks
	}
	tests := []struct {
		name            string
		args            args
		wantUriTemplate string
		isExpectedErr   bool
	}{
		{
			name: "normal: URI template with single variable",
			args: args{
				uriTemplate: "file://sample/{variable}",
				callbacks: &ResourceTemplateCallbacks{
					List: func() schema.ListResourcesResultSchema {
						return schema.ListResourcesResultSchema{}
					},
					Complete: map[string]CompleteResourceCallback{
						"variable": func(value string) []string {
							return []string{"value1", "value2"}
						},
					},
				},
			},
			wantUriTemplate: "file://sample/{variable}",
			isExpectedErr:   false,
		},
		{
			name: "normal: URI template with multiple variables",
			args: args{
				uriTemplate: "file://project/{project}/user/{username}",
				callbacks: &ResourceTemplateCallbacks{
					List: func() schema.ListResourcesResultSchema {
						return schema.ListResourcesResultSchema{}
					},
					Complete: map[string]CompleteResourceCallback{
						"project": func(value string) []string {
							return []string{"project1", "project2"}
						},
						"username": func(value string) []string {
							return []string{"user1", "user2"}
						},
					},
				},
			},
			wantUriTemplate: "file://project/{project}/user/{username}",
			isExpectedErr:   false,
		},
		{
			name: "normal: URI template without variables",
			args: args{
				uriTemplate: "file://static/resource",
				callbacks: &ResourceTemplateCallbacks{
					List: func() schema.ListResourcesResultSchema {
						return schema.ListResourcesResultSchema{}
					},
					Complete: map[string]CompleteResourceCallback{},
				},
			},
			wantUriTemplate: "file://static/resource",
			isExpectedErr:   false,
		},
		{
			name: "seminormal: URI template with nil callbacks",
			args: args{
				uriTemplate: "file://sample/{variable}",
				callbacks:   nil,
			},
			wantUriTemplate: "file://sample/{variable}",
			isExpectedErr:   false,
		},
		{
			name: "seminormal: Invalid URI template syntax",
			args: args{
				uriTemplate: "file://sample/{unclosed",
				callbacks:   nil,
			},
			wantUriTemplate: "",
			isExpectedErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut, err := NewResourceTemplate(
				tt.args.uriTemplate,
				tt.args.callbacks,
			)
			if (err != nil) != tt.isExpectedErr {
				t.Errorf("NewResourceTemplate() error = %v, isExpectedErr = %v", err, tt.isExpectedErr)
				return
			}

			// エラーが期待されている場合はテスト終了
			if tt.isExpectedErr {
				return
			}

			// 実際の値と期待値を比較
			if diff := cmp.Diff(tt.wantUriTemplate, sut.uriTemplate().ToString()); diff != "" {
				t.Errorf("ResourceTemplate.uriTemplate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestResourceTemplate_ListCallback(t *testing.T) {
	// テスト用のListResourcesResultSchema
	expectedResult := schema.ListResourcesResultSchema{
		Resources: []schema.ResourceSchema{
			{Uri: "file://sample/test1.txt"},
			{Uri: "file://sample/test2.txt"},
		},
	}

	tests := []struct {
		name             string
		callbacks        *ResourceTemplateCallbacks
		wantListCallback bool
		wantListResult   schema.ListResourcesResultSchema
	}{
		{
			name: "normal: List callback exists",
			callbacks: &ResourceTemplateCallbacks{
				List: func() schema.ListResourcesResultSchema {
					return expectedResult
				},
				Complete: map[string]CompleteResourceCallback{},
			},
			wantListCallback: true,
			wantListResult:   expectedResult,
		},
		{
			name:             "seminormal: callbacks is nil",
			callbacks:        nil,
			wantListCallback: false,
			wantListResult:   schema.ListResourcesResultSchema{},
		},
		{
			name: "seminormal: List function is nil",
			callbacks: &ResourceTemplateCallbacks{
				List:     nil,
				Complete: map[string]CompleteResourceCallback{},
			},
			wantListCallback: false,
			wantListResult:   schema.ListResourcesResultSchema{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ResourceTemplateを作成
			rt, err := NewResourceTemplate("file://sample/{variable}", tt.callbacks)
			if err != nil {
				t.Fatalf("Failed to create ResourceTemplate: %v", err)
			}

			// ListCallback()メソッドの結果を取得
			listCallback := rt.ListCallback()

			// コールバックの有無を確認
			if (listCallback != nil) != tt.wantListCallback {
				t.Errorf("ResourceTemplate.ListCallback() exists = %v, want %v", listCallback != nil, tt.wantListCallback)
			}

			// コールバックが存在する場合、結果を確認
			if listCallback != nil {
				result := listCallback()
				if diff := cmp.Diff(tt.wantListResult, result); diff != "" {
					t.Errorf("ListCallback() result mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestResourceTemplate_CompleteCallBack(t *testing.T) {
	tests := []struct {
		name                 string
		callbacks            *ResourceTemplateCallbacks
		variable             string
		wantCompleteCallback bool
		wantCompletions      []string
	}{
		{
			name: "normal: Complete callback exists",
			callbacks: &ResourceTemplateCallbacks{
				List: func() schema.ListResourcesResultSchema {
					return schema.ListResourcesResultSchema{}
				},
				Complete: map[string]CompleteResourceCallback{
					"variable": func(value string) []string {
						return []string{"value1", "value2", "value3"}
					},
				},
			},
			variable:             "variable",
			wantCompleteCallback: true,
			wantCompletions:      []string{"value1", "value2", "value3"},
		},
		{
			name: "seminormal: Variable name does not exist",
			callbacks: &ResourceTemplateCallbacks{
				List: func() schema.ListResourcesResultSchema {
					return schema.ListResourcesResultSchema{}
				},
				Complete: map[string]CompleteResourceCallback{
					"variable": func(value string) []string {
						return []string{"value1", "value2"}
					},
				},
			},
			variable:             "nonexistent",
			wantCompleteCallback: false,
			wantCompletions:      nil,
		},
		{
			name:                 "seminormal: callbacks is nil",
			callbacks:            nil,
			variable:             "variable",
			wantCompleteCallback: false,
			wantCompletions:      nil,
		},
		{
			name: "normal: Filtering based on input value",
			callbacks: &ResourceTemplateCallbacks{
				List: func() schema.ListResourcesResultSchema {
					return schema.ListResourcesResultSchema{}
				},
				Complete: map[string]CompleteResourceCallback{
					"variable": func(value string) []string {
						if value == "a" {
							return []string{"apple", "apricot", "avocado"}
						}
						return []string{"banana", "blueberry"}
					},
				},
			},
			variable:             "variable",
			wantCompleteCallback: true,
			wantCompletions:      []string{"apple", "apricot", "avocado"}, // "a"を入力した場合
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ResourceTemplateを作成
			rt, err := NewResourceTemplate("file://sample/{variable}", tt.callbacks)
			if err != nil {
				t.Fatalf("Failed to create ResourceTemplate: %v", err)
			}

			// CompleteCallBack()メソッドの結果を取得
			completeCallback := rt.CompleteCallBack(tt.variable)

			// コールバックの有無を確認
			if (completeCallback != nil) != tt.wantCompleteCallback {
				t.Errorf("ResourceTemplate.CompleteCallBack(%q) exists = %v, want %v",
					tt.variable, completeCallback != nil, tt.wantCompleteCallback)
			}

			// コールバックが存在する場合、結果を確認
			if completeCallback != nil {
				// テスト用の入力値を指定
				inputValue := "a"
				result := completeCallback(inputValue)
				if diff := cmp.Diff(tt.wantCompletions, result); diff != "" {
					t.Errorf("CompleteCallBack(%q) result mismatch (-want +got):\n%s", inputValue, diff)
				}
			}
		})
	}
}
