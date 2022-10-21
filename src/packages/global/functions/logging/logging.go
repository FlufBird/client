package logging

import (
	"github.com/FlufBird/client/src/packages/global/variables"

	"fmt"
	"time"
)

func Information(section string, message string, arguments ...any) {
	Log("INFORMATION", section, message, arguments...)
}

func Warning(section string, message string, arguments ...any) {
	Log("WARNING", section, message, arguments...)
}

func Error(section string, message string, arguments ...any) {
	Log("ERROR", section, message, arguments...)
}

func Fatal(section string, message string, arguments ...any) {
	Log("FATAL", section, message, arguments...)
}

func Log(_type string, section string, message string, arguments ...any) {
	if !variables.DevelopmentMode {
		return
	}

	_time := time.Now().Format("15:04:05")

	fmt.Printf(
		"[%s] [%s] [%s] %s\n",

		_time,
		_type,
		section,
		fmt.Sprintf(message, arguments...),
	)
}