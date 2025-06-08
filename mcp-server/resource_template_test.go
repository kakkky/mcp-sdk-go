package mcpserver

import (
	"reflect"
	"testing"

	utilities "github.com/kakkky/mcp-sdk-go/shared/utilities/uri-template"
)

func TestResourceTemplate_uriTemplate(t *testing.T) {
	type fields struct {
		uriTemp   *utilities.UriTemplate
		callBacks *ResourceTemplateCallbacks
	}
	tests := []struct {
		name   string
		fields fields
		want   *utilities.UriTemplate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ResourceTemplate{
				uriTemp:   tt.fields.uriTemp,
				callBacks: tt.fields.callBacks,
			}
			if got := r.uriTemplate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResourceTemplate.uriTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
