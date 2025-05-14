package utilities

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jtacoma/uritemplates"
)

type UriTemplate struct {
	*uritemplates.UriTemplate
}

func NewUriTemplate(template string) (*UriTemplate, error) {
	uritemplate, err := uritemplates.Parse(template)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URI template: %w", err)
	}
	return &UriTemplate{
		UriTemplate: uritemplate,
	}, nil
}

func (u *UriTemplate) ToString() string {
	return u.UriTemplate.String()
}

// テンプレートの展開
func (u *UriTemplate) Expand(value map[string]any) (string, error) {
	return u.UriTemplate.Expand(value)
}

// テンプレート内で使用されている全ての変数名を取得
func (u *UriTemplate) VariableNames() []string {
	return u.UriTemplate.Names()
}

func (u *UriTemplate) Match(uri string) (map[string]any, error) {
	// テンプレートを正規表現パターンに変換
	pattern, names, err := u.templateToRegex()
	if err != nil {
		return nil, fmt.Errorf("unable to convert template to regex: %w", err)
	}

	// 正規表現でマッチング
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}
	// マッチした値を取得
	matches := regex.FindStringSubmatch(uri)
	if matches == nil {
		return nil, nil // マッチしなかった場合はnilを返す
	}

	// 変数名とマッチして抽出した値のマップを構築
	result := make(map[string]any)
	for i, name := range names {
		if i+1 < len(matches) {
			value := matches[i+1]

			// exploded変数かどうかをチェック
			exploded := strings.HasSuffix(name, "*")
			cleanName := strings.TrimSuffix(name, "*")

			// exploded変数で、カンマを含む場合は配列に分割
			if exploded && strings.Contains(value, ",") {
				result[cleanName] = strings.Split(value, ",")
			} else {
				result[cleanName] = value
			}
		}
	}

	return result, nil
}

// テンプレートから、正規表現パターン文字列と変数名のリストを返す
func (u *UriTemplate) templateToRegex() (pattern string, names []string, err error) {
	// テンプレート文字列を取得
	templateStr := u.String()

	// 正規表現パターンを構築するための文字列ビルダー
	var patternBuilder strings.Builder
	// 先頭を表すメタ文字を追加
	patternBuilder.WriteString("^")

	// テンプレート内の変数を正規表現パターンに置き換え
	segments := strings.Split(templateStr, "{")
	// 最初のセグメントをエスケープして追加
	patternBuilder.WriteString(regexp.QuoteMeta(segments[0]))

	for i := 1; i < len(segments); i++ {
		if !strings.Contains(segments[i], "}") {
			return "", nil, fmt.Errorf("invalid template format: missing closing brace")
		}
		// hoge}/fuga/ のような文字列を、hoge と　/fuga/ に分割
		parts := strings.SplitN(segments[i], "}", 2)
		varExpr := parts[0] // 例：hoge
		literal := parts[1] // 例：/fuga/

		// 変数式の解析
		varNames, varPattern := parseVarExpression(varExpr)
		names = append(names, varNames...)
		patternBuilder.WriteString(varPattern)

		// リテラル部分の追加
		patternBuilder.WriteString(regexp.QuoteMeta(literal))
	}
	// 末尾を表すメタ文字を追加s
	patternBuilder.WriteString("$")
	return patternBuilder.String(), names, nil
}

// 変数の配列と、付加されている演算子に対応した正規表現パターンを返す
// 例：{userId} -> userId, "([^/]+)"
func parseVarExpression(expr string) (names []string, pattern string) {
	// 演算子のチェック
	op := ""
	if len(expr) > 0 {
		switch expr[0] {
		case '+', '.', '/', ';', '?', '&', '#':
			op = string(expr[0])
			// 演算子を除去
			expr = expr[1:]
		}
	}

	// 変数のリストを取得
	varList := strings.Split(expr, ",")
	var patternBuilder strings.Builder

	switch op {
	// 例：/api{?page,limit}　-> ["page","limit"] , \\?page=([^&]+)&limit=([^&]+)
	case "?", "&": // クエリパラメータ用の特別処理
		for i, v := range varList {
			names = append(names, v)
			// 特殊文字をエスケープ
			if i == 0 {
				patternBuilder.WriteString("\\" + op)
			} else {
				// 2つ目以降は?を&に変換する
				// 例：?page=2&limit=10 -> ?page=2&limit=10
				op = strings.Replace(op, "?", "&", 1)
				patternBuilder.WriteString(op)
			}
			// メタ文字をエスケープして追加
			patternBuilder.WriteString(regexp.QuoteMeta(v))
			patternBuilder.WriteString("=([^&]+)")
		}

	case "#": // フラグメント用の特別処理
		patternBuilder.WriteString("#")
		for _, v := range varList {
			names = append(names, v)
			patternBuilder.WriteString("(.+)")
		}

	case "/": // パスセグメント用の特別処理
		for _, v := range varList {
			names = append(names, v)
			// explodedかどうかをチェック
			exploded := strings.HasSuffix(v, "*")

			patternBuilder.WriteString("/")
			if exploded {
				patternBuilder.WriteString("([^/]+(?:,[^/]+)*)")
			} else {
				patternBuilder.WriteString("([^/,]+)")
			}
		}

	case ".": // ドット区切り用の特別処理
		for _, v := range varList {
			names = append(names, v)
			patternBuilder.WriteString("\\.([^/,]+)")
		}

	case "+": // 予約されていない文字用の特別処理
		for _, v := range varList {
			names = append(names, v)
			patternBuilder.WriteString("(.+)")
		}

	default: // 標準のパスパラメータ処理
		for _, v := range varList {
			// 変数名を保存
			names = append(names, v)

			// exploded変数かどうかをチェック
			exploded := strings.HasSuffix(v, "*")

			// 単純なキャプチャグループパターン
			if exploded {
				patternBuilder.WriteString("([^/]+(?:,[^/]+)*)")
			} else {
				patternBuilder.WriteString("([^/,]+)")
			}
		}
	}

	return names, patternBuilder.String()
}
