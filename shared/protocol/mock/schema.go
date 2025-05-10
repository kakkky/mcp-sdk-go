package mock

// テスト用リクエストスキーマ
type TestRequestSchema struct {
	MethodName string `json:"method"`
	ParamsData any    `json:"params"`
}

func (r *TestRequestSchema) Method() string {
	return r.MethodName
}

func (r *TestRequestSchema) Params() any {
	return r.ParamsData
}

// テスト用レスポンススキーマ
type TestResultShema struct {
	Status string `json:"status"`
}

func (r *TestResultShema) Result() any {
	return r.Status
}

// テスト用通知スキーマ
type TestNotificationSchema struct {
	MethodName string `json:"method"`
	ParamsData any    `json:"params"`
}

func (n *TestNotificationSchema) Method() string {
	return n.MethodName
}

func (n *TestNotificationSchema) Params() any {
	return n.ParamsData
}
