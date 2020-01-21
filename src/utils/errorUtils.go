package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

const maxStackLength = 50

// Error is the type that implements the error interface.
// It contains the underlying err and its stacktrace.
type StackError struct {
	Err        error
	StackTrace string
}

func (m StackError) Error() string {
	return m.Err.Error() + m.StackTrace
}

// Wrap annotates the given error with a stack trace
func WrapError(err error) StackError {
	return StackError{Err: err, StackTrace: getStackTrace()}
}

func Log(level logrus.Level, err error, msg string) {
	if err != nil {
		logrus.New().Log(level, msg+"\n", WrapError(err))
		if level == logrus.FatalLevel {
			os.Exit(1)
		}
	} else {
		logrus.New().Log(level, msg)
	}
}

func PrintBanner(msg string) {
	logrus.Info("====================")
	logrus.Info(msg)
	logrus.Info("====================")
}

func getStackTrace() string {
	stackBuf := make([]uintptr, maxStackLength)
	length := runtime.Callers(3, stackBuf[:])
	stack := stackBuf[:length]

	trace := ""
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			trace = trace + fmt.Sprintf("\n\t %s \n\t\t %s:%d", frame.Function, frame.File, frame.Line)
		}
		if !more {
			break
		}
	}
	return trace
}
