package protocol

import (
	"fmt"

	"dario.cat/mergo"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func MergeCapabilities[T schema.ServerCapabilities | schema.ClientCapabilities](base, additional T) T {
	if err := mergo.Merge(&base, additional, mergo.WithOverride, mergo.WithTypeCheck, mergo.WithOverrideEmptySlice); err != nil {
		fmt.Println("Error merging capabilities:", err)
	}
	return base
}
