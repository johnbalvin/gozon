package trace

import "fmt"

func NewOrAdd(step int, pkgName, funName string, err error, desc string) Traceback {
	//on any scenario err could be nil
	var currentTrace Traceback
	if err == nil {
		return Traceback{
			MainErr: ErrorData{
				PackageName: pkgName,
				FuncName:    funName,
				Step:        step,
				Description: fmt.Sprintf("*%s*", desc),
			}}
	} else {
		var ok bool
		currentTrace, ok = err.(Traceback)
		if !ok {
			return Traceback{
				MainErr: ErrorData{
					PackageName: pkgName,
					FuncName:    funName,
					Step:        step,
					Err:         err,
					Description: fmt.Sprintf("*%s*", desc),
					ErrString:   err.Error(),
				}}
		}
	}
	errData := ErrorData{
		Step:        step,
		PackageName: pkgName,
		FuncName:    funName,
		Description: desc,
	}
	currentTrace.Trace = append(currentTrace.Trace, errData)
	return currentTrace
}
