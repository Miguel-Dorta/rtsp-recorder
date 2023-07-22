package log

import (
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/http_push"
	"os"
)

func init() {
	logolang.DefaultLogger.Color = false
}

func Critical(message string) {
	logolang.Critical(message)
	http_push.Report(false, message)
	os.Exit(1)
}

func Criticalf(format string, a ...any) {
	Critical(fmt.Sprintf(format, a...))
}

func Error(message string) {
	logolang.Error(message)
	http_push.Report(false, message)
}

func Errorf(format string, a ...any) {
	Error(fmt.Sprintf(format, a...))
}
