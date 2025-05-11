package server

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kakkky/mcp-sdk-go/mcp-server/server/mock"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
	"go.uber.org/mock/gomock"
)

func TestServer_CreateMessage(t *testing.T) {
	type args struct {
		params      any
		contentType string
	}
	tests := []struct {
		name          string
		args          args
		expected      schema.Result
		isExpectedErr bool
		mockFn        func(mp *mock.MockProtocol, contentType string, params any, resultSchema schema.Result)
	}{
		{
			name: "normal : able to create message (text)",
			args: args{
				params: schema.CreateMessageRequestParams[schema.TextContentSchema]{
					Messages: []schema.SamplingMessageSchema[schema.TextContentSchema]{
						{
							Role: "user",
							Content: schema.TextContentSchema{
								Type: "text",
								Text: "Hello",
							},
						},
					},
				},
				contentType: "text",
			},
			expected: &schema.CreateMessageResultSchema[schema.TextContentSchema]{
				Model: "test-model",
				Role:  "assistant",
				Content: schema.TextContentSchema{
					Type: "text",
					Text: "Hello , user",
				},
			},
			isExpectedErr: false,
			mockFn: func(mp *mock.MockProtocol, contentType string, params any, resultSchema schema.Result) {
				mp.EXPECT().
					Request(
						&schema.CreateMessageRequestSchema[schema.TextContentSchema]{
							MethodName: "sampling/createMessage",
							ParamsData: schema.CreateMessageRequestParams[schema.TextContentSchema]{
								Messages: []schema.SamplingMessageSchema[schema.TextContentSchema]{
									{
										Role: "user",
										Content: schema.TextContentSchema{
											Type: "text",
											Text: "Hello",
										},
									},
								},
							},
						},
						&schema.CreateMessageResultSchema[schema.TextContentSchema]{},
					).
					Return(
						&schema.CreateMessageResultSchema[schema.TextContentSchema]{
							Model: "test-model",
							Role:  "assistant",
							Content: schema.TextContentSchema{
								Type: contentType,
								Text: params.(schema.CreateMessageRequestParams[schema.TextContentSchema]).Messages[0].Content.Text + " , " + params.(schema.CreateMessageRequestParams[schema.TextContentSchema]).Messages[0].Role,
							},
						},
						nil,
					)
			},
		},
		{
			name: "semi normal : returns an error if the params type is incorrect",
			args: args{
				params:      "invalid params",
				contentType: "text",
			},
			isExpectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			sut := NewServer(schema.Implementation{}, nil)
			// Protocolをモックに差し替え
			mockProtocol := mock.NewMockProtocol(ctrl)
			sut.Protocol = mockProtocol
			// モックの動作を設定
			if tt.mockFn != nil {
				tt.mockFn(mockProtocol, tt.args.contentType, tt.args.params, tt.expected)
			}

			got, err := sut.CreateMessage(tt.args.params, tt.args.contentType)
			if (err != nil) != tt.isExpectedErr {
				t.Errorf("CreateMessage() error = %v, wantErr %v", err, tt.isExpectedErr)
				return
			}
			if diff := cmp.Diff(got, tt.expected); diff != "" {
				t.Errorf("CreateMessage() got = %v, want %v, diff %s", got, tt.expected, diff)
			}
		})
	}
}
