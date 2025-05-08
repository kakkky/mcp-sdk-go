package schema

type Notification interface {
	Method() string
	Params() any
}
