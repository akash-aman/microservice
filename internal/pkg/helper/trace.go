package helper

import (
	"fmt"
	"runtime"
	"strings"
)

func Caller(i int) string {
	_, file, line, ok := runtime.Caller(i)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", file, line)
}

/**
 *
 * StackTrace returns a formatted stack trace as a string.
 * The stack trace starts from the caller specified by the skip parameter
 * and includes up to the number of frames specified by the depth parameter.
 * Parameters:
 * 	- skip: The number of stack frames to skip before recording the stack trace.
 * 	- depth: The maximum number of stack frames to record.
 *
 * Returns:
 * A string representing the formatted stack trace.
 */
func StackTrace(skip int, depth int) string {
	pcs := make([]uintptr, depth)
	n := runtime.Callers(skip, pcs[:])

	frames := runtime.CallersFrames(pcs[:n])
	var stack strings.Builder
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		fmt.Fprintf(&stack, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
	}
	return stack.String()
}
