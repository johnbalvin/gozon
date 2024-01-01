package trace

import (
	"encoding/json"
	"fmt"
)

func GetMainErr(err error) error {
	if err == nil {
		return nil
	}
	mainErr, ok := err.(Traceback)
	if ok {
		return mainErr.MainErr.Err
	}
	return err
}
func (tracebak Traceback) MainError() error {
	return tracebak.MainErr.Err
}
func (tb Traceback) Error() string {
	finalMsg, _ := json.Marshal(tb)
	path := fmt.Sprintf("trace.NewOrAdd(%d, \"%s\", \"%s\"", tb.MainErr.Step, tb.MainErr.PackageName, tb.MainErr.FuncName)
	finalFrmt := fmt.Sprintf("{PATH:  %s   VALUE: %s}", path, finalMsg)
	return finalFrmt
}
