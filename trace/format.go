package trace

type Traceback struct {
	MainErr ErrorData
	Trace   []ErrorData
}
type ErrorData struct {
	PackageName string
	FuncName    string
	Step        int
	Description string `json:"Description,omitempty"`
	Err         error  `json:"-"`
	ErrString   string `json:"Err,omitempty"`
}
