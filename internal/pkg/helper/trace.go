package helper

import (
	"fmt"
	"runtime"
)

func Caller() string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func StackTrace() string {
	stackBuf := make([]byte, 1024)
	n := runtime.Stack(stackBuf, false)
	return string(stackBuf[:n])
}
