package protocol

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func TestMergeCapabilities(t *testing.T) {
	tests := []struct {
		name                 string
		baseServerCapa       schema.ServerCapabilities
		additionalServerCapa schema.ServerCapabilities
		expectedServerCapa   schema.ServerCapabilities
	}{
		{
			name: "normal : able to merge server capabilities",
			baseServerCapa: schema.ServerCapabilities{
				Tools: &schema.Tools{
					ListChanged: false,
				},
				Resources: &schema.Resources{
					ListChanged: false,
				},
			},
			additionalServerCapa: schema.ServerCapabilities{
				Prompts: &schema.Prompts{
					ListChanged: true,
				},
			},
			expectedServerCapa: schema.ServerCapabilities{
				Tools: &schema.Tools{
					ListChanged: false,
				},
				Resources: &schema.Resources{
					ListChanged: false,
				},
				Prompts: &schema.Prompts{
					ListChanged: true,
				},
			},
		},
		{
			name: "normal : able to overwrite server capabilities",
			baseServerCapa: schema.ServerCapabilities{
				Tools: &schema.Tools{
					ListChanged: false,
				},
				Resources: &schema.Resources{
					ListChanged: false,
				},
				Prompts: &schema.Prompts{
					ListChanged: true,
				},
			},
			additionalServerCapa: schema.ServerCapabilities{
				Tools: &schema.Tools{
					ListChanged: true,
				},
				Resources: &schema.Resources{
					ListChanged: false,
					Subscribe:   true,
				},
				Prompts: &schema.Prompts{
					ListChanged: true,
				},
				Logging: &schema.Logging{},
			},
			expectedServerCapa: schema.ServerCapabilities{
				Tools: &schema.Tools{
					ListChanged: true,
				},
				Resources: &schema.Resources{
					ListChanged: false,
					Subscribe:   true,
				},
				Prompts: &schema.Prompts{
					ListChanged: true,
				},
				Logging: &schema.Logging{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotServerCapa := MergeCapabilities(tt.baseServerCapa, tt.additionalServerCapa)
			if diff := cmp.Diff(gotServerCapa, tt.expectedServerCapa); diff != "" {
				t.Errorf("MergeCapabilities() server capabilities mismatch (-got +expected):\n%s", diff)
			}
		})
	}

}
