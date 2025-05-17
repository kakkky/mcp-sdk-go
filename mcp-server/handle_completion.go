package mcpserver

import (
	"fmt"

	mcperr "github.com/kakkky/mcp-sdk-go/shared/mcp-err"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func (m *McpServer) handleResourceCompletion(request schema.CompleteRequestSchema, ref schema.ResourceReferenceSchema) (*schema.CompleteResultSchema, error) {
	params := request.Params().(*schema.CompleteRequestParams)

	template := m.registeredResourceTemplates[ref.UriOrName()]
	// テンプレートが見つからなかったが、固定リソースが見つかった場合は空の補完を返す（しかし、リクエストエラーとすべきだろう）
	if template == nil {
		if m.registeredResources[ref.UriOrName()] != nil {
			return EmptyCompletionResult(), nil
		}
		return nil, mcperr.NewMcpErr(
			mcperr.INVALID_PARAMS,
			fmt.Sprintf("resource template %s not found", params.Ref.UriOrName()),
			nil,
		)
	}
	// テンプレートの変数名から補完コールバックを取得
	completer := template.resourceTemplate.CompleteCallBack(params.Argument.Name)
	if completer == nil {
		return EmptyCompletionResult(), nil
	}
	suggestions := completer(params.Argument.Value)
	return createCompletionResult(suggestions), nil
}

func createCompletionResult(suggestions []string) *schema.CompleteResultSchema {
	suggestLengh := len(suggestions)
	isHasMore := suggestLengh > 100
	return &schema.CompleteResultSchema{
		Completion: schema.CompletionSchema{
			Values:  suggestions,
			Total:   &suggestLengh,
			HasMore: &isHasMore,
		},
	}
}

// 空の補完結果を返す
func EmptyCompletionResult() *schema.CompleteResultSchema {
	notHasMore := false
	return &schema.CompleteResultSchema{
		Completion: schema.CompletionSchema{
			Values:  []string{},
			Total:   nil,
			HasMore: &notHasMore,
		},
	}
}
