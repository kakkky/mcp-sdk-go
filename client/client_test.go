package client

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kakkky/mcp-sdk-go/client/mock"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
	"go.uber.org/mock/gomock"
)

func TestClient_ValidateCapabilities(t *testing.T) {
	tests := []struct {
		name           string
		capability     any
		method         string
		serverCapSetup func(*schema.ServerCapabilities)
		wantErr        bool
		errorContains  string
	}{
		{
			name:       "normal: logging capability is supported",
			capability: &schema.Logging{},
			method:     "logging/message",
			serverCapSetup: func(sc *schema.ServerCapabilities) {
				sc.Logging = &schema.Logging{}
			},
			wantErr: false,
		},
		{
			name:       "semi normal: logging capability is not supported",
			capability: &schema.Logging{},
			method:     "logging/message",
			serverCapSetup: func(sc *schema.ServerCapabilities) {
				sc.Logging = nil
			},
			wantErr:       true,
			errorContains: "logging/message requires logging capability",
		},
		{
			name:       "normal: completion capability is supported",
			capability: &schema.Completion{},
			method:     "completion/complete",
			serverCapSetup: func(sc *schema.ServerCapabilities) {
				sc.Completion = &schema.Completion{}
			},
			wantErr: false,
		},
		{
			name:       "semi normal: completion capability is not supported",
			capability: &schema.Completion{},
			method:     "completion/complete",
			serverCapSetup: func(sc *schema.ServerCapabilities) {
				sc.Completion = nil
			},
			wantErr:       true,
			errorContains: "completion/complete requires completion capability",
		},
		{
			name:       "normal: prompts capability is supported",
			capability: &schema.Prompts{},
			method:     "prompts/list",
			serverCapSetup: func(sc *schema.ServerCapabilities) {
				sc.Prompts = &schema.Prompts{}
			},
			wantErr: false,
		},
		{
			name:       "semi normal: prompts capability is not supported",
			capability: &schema.Prompts{},
			method:     "prompts/list",
			serverCapSetup: func(sc *schema.ServerCapabilities) {
				sc.Prompts = nil
			},
			wantErr:       true,
			errorContains: "prompts/list requires prompts capability",
		},
		{
			name:       "normal: resources capability is supported",
			capability: &schema.Resources{},
			method:     "resources/read",
			serverCapSetup: func(sc *schema.ServerCapabilities) {
				sc.Resources = &schema.Resources{}
			},
			wantErr: false,
		},
		{
			name:       "semi normal: resources capability is not supported",
			capability: &schema.Resources{},
			method:     "resources/read",
			serverCapSetup: func(sc *schema.ServerCapabilities) {
				sc.Resources = nil
			},
			wantErr:       true,
			errorContains: "resources/read requires resources capability",
		},
		{
			name:       "normal: tools capability is supported",
			capability: &schema.Tools{},
			method:     "tools/list",
			serverCapSetup: func(sc *schema.ServerCapabilities) {
				sc.Tools = &schema.Tools{}
			},
			wantErr: false,
		},
		{
			name:       "semi normal: tools capability is not supported",
			capability: &schema.Tools{},
			method:     "tools/list",
			serverCapSetup: func(sc *schema.ServerCapabilities) {
				sc.Tools = nil
			},
			wantErr:       true,
			errorContains: "tools/list requires tools capability",
		},
		{
			name:       "semi normal: unknown capability type",
			capability: "unknown type",
			method:     "unknown/method",
			serverCapSetup: func(sc *schema.ServerCapabilities) {
				// 設定不要
			},
			wantErr:       true,
			errorContains: "unknown/method unknown capability type for method", // 修正したエラーメッセージ
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewClient(schema.Implementation{}, &ClientOptions{
				Capabilities: schema.ClientCapabilities{},
			})
			// テストケース固有のサーバー機能設定を適用
			tt.serverCapSetup(&sut.serverCapabilities)

			// テスト対象メソッドを実行
			err := sut.ValidateCapabilities(tt.capability, tt.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCapabilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// エラーメッセージの検証を追加
			if err != nil && tt.errorContains != "" {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("ValidateCapabilities() error = %v, want error containing %v", err, tt.errorContains)
				}
			}
		})
	}
}

func TestClient_Connect(t *testing.T) {
	tests := []struct {
		name                       string
		mockFn                     func(*mock.MockProtocol)
		expectedServerCapabilities schema.ServerCapabilities
		isExpectedErr              bool
	}{
		{
			name: "normal: successful connection",
			mockFn: func(mp *mock.MockProtocol) {
				mp.EXPECT().
					Connect(gomock.Any()).
					Do(
						func(protocol.Transport) {
							// 擬似的にトランスポートの接続確立を通知
							TransportStartedNotify <- struct{}{}

						},
					).Return(nil)
				mp.EXPECT().
					Request(
						gomock.Any(),
						gomock.Any(),
					).
					Return(
						&schema.InitializeResultSchema{
							ProtocolVersion: schema.LATEST_PROTOCOL_VERSION,
							Capabilities: schema.ServerCapabilities{
								Tools: &schema.Tools{
									ListChanged: true,
								},
							},
							ServerInfo: schema.Implementation{
								Name:    "test-server",
								Version: "1.0.0",
							},
						},
						nil,
					)
				mp.EXPECT().
					Notificate(
						gomock.Any(),
					).Return(nil)

			},
			expectedServerCapabilities: schema.ServerCapabilities{
				Tools: &schema.Tools{
					ListChanged: true,
				},
			},
		},
		{
			name: "semi normal: unsupported protocol version",
			mockFn: func(mp *mock.MockProtocol) {
				mp.EXPECT().
					Connect(gomock.Any()).
					Do(
						func(protocol.Transport) {
							// 擬似的にトランスポートの接続確立を通知
							TransportStartedNotify <- struct{}{}
						},
					).Return(nil)
				mp.EXPECT().
					Request(
						gomock.Any(),
						gomock.Any(),
					).
					Return(
						&schema.InitializeResultSchema{
							ProtocolVersion: "2021-01-01", // サポートされていないバージョン
							Capabilities: schema.ServerCapabilities{
								Tools: &schema.Tools{
									ListChanged: true,
								},
							},
							ServerInfo: schema.Implementation{
								Name:    "test-server",
								Version: "1.0.0",
							},
						},
						nil,
					)
				mp.EXPECT().
					Close().Return(nil)
			},
			isExpectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProtocol := mock.NewMockProtocol(ctrl)
			if tt.mockFn != nil {
				tt.mockFn(mockProtocol)
			}

			sut := NewClient(
				schema.Implementation{Name: "test-client", Version: "1.0.0"},
				&ClientOptions{
					Capabilities: schema.ClientCapabilities{},
				},
			)
			sut.Protocol = mockProtocol

			mockTransport := mock.NewMockTransport(ctrl)

			errCh := make(chan error)
			go func() {
				err := sut.Connect(mockTransport)
				if err != nil {
					fmt.Println("Connect() error:", err)
					errCh <- err
					return
				}
			}()
			select {
			case err := <-errCh:
				if (err != nil) != tt.isExpectedErr {
					t.Errorf("Connect() error = %v, isExpectedErr %v", err, tt.isExpectedErr)
					return
				}
				if tt.isExpectedErr {
					return
				}
			case <-OpetationPhaseStartNotify:
				// serverCapabilitiesを期待通りに設定されているか
				if diff := cmp.Diff(tt.expectedServerCapabilities, sut.serverCapabilities, cmp.AllowUnexported(Client{})); diff != "" {
					t.Errorf("Connect() client mismatch (-want +got):\n%s", diff)
				}
			}

		})
	}
}
