package schema

type Request interface {
	Method() string
	Params() any
}
