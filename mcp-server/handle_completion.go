package mcpserver

import (
	"fmt"
	"strings"

	mcperr "github.com/kakkky/mcp-sdk-go/shared/mcp-err"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func (m *McpServer) handleResourceCompletion(request schema.CompleteRequestSchema, ref schema.ResourceReferenceSchema) (*schema.CompleteResultSchema, error) {
	params := request.Params().(schema.CompleteRequestParams)
	// refで渡されたリソーステンプレートのURIと一致するテンプレートを探す
	var template *RegisteredResourceTemplate
	for _, registeredTemplate := range m.registeredResourceTemplates {
		if registeredTemplate.resourceTemplate.uriTemplate().ToString() == ref.UriOrName() {
			template = registeredTemplate
			break
		}
	}
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

func (m *McpServer) handlePromptCompletion(request schema.CompleteRequestSchema, ref schema.PromptReferenceSchema) (*schema.CompleteResultSchema, error) {
	params := request.Params().(schema.CompleteRequestParams)
	prompt, ok := m.registeredPrompts[ref.UriOrName()]
	if !ok {
		return nil, mcperr.NewMcpErr(
			mcperr.INVALID_PARAMS,
			fmt.Sprintf("prompt %s not found", params.Ref.UriOrName()),
			nil,
		)
	}
	if !prompt.enabled {
		return nil, mcperr.NewMcpErr(
			mcperr.INVALID_PARAMS,
			fmt.Sprintf("prompt %s disabled", params.Ref.UriOrName()),
			nil,
		)
	}
	if prompt.argsSchema == nil {
		return EmptyCompletionResult(), nil
	}
	searchingArgName := params.Argument.Name
	currentValue := params.Argument.Value
	var suggestions []string
	for _, arg := range prompt.argsSchema {
		if searchingArgName == arg.Name {
			for _, value := range arg.CompletionValues {
				if strings.HasPrefix(value, currentValue) {
					suggestions = append(suggestions, value)
				}
			}
		}
	}
	return createCompletionResult(suggestions), nil
}

func createCompletionResult(suggestions []string) *schema.CompleteResultSchema {
	suggestLengh := len(suggestions)
	isHasMore := suggestLengh > 100
	return &schema.CompleteResultSchema{
		Completion: schema.CompletionSchema{
			Values:  suggestions,
			Total:   suggestLengh,
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
			Total:   0,
			HasMore: &notHasMore,
		},
	}
}
