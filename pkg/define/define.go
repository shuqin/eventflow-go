package define

const (
	CONTINUE = 200
	TERMINAL = -1
)

type IllegalArgumentError struct {
	Msg string
	Err error
}

func (e *IllegalArgumentError) Error() string {
	return e.Msg
}

type FlowResult struct {
	Code int
}

// Do not need declared explicitly
//func NewFlowResult(code int) FlowResult {
//	return FlowResult{code}
//}

type EventData interface {
	GetData() interface{}
	GetType() string
	GetEventId() string
	SetOsType(i int) error
}
