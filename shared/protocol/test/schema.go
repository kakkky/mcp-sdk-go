package test

type TestRequest struct {
	MethodName string `json:"method"`
	ParamsData any    `json:"params"`
}

func (r *TestRequest) Method() string {
	return r.MethodName
}

func (r *TestRequest) Params() any {
	return r.ParamsData
}

type TestResult struct {
	Status string `json:"status"`
}

func (r *TestResult) Result() any {
	return r.Status
}
