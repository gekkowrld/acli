package errors

import (
	"fmt"
	slog "log"
	"os"

	"github.com/charmbracelet/log"
)

var pref string = "acli "
var logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportTimestamp: true,
	Prefix:          pref,
})

// Print out error messages
func Error(content interface{}, dispColour bool) {
	str := fmt.Sprintf("%v", content)
	if !dispColour {
		str = fmt.Sprintf("ERRO %s: %s\n", pref, str)
		slog.Printf(str)
		os.Exit(1)
	} else {
		logger.Error(str)
		os.Exit(1)
	}
}

// Print out warning messages
func Warn(content interface{}, dispColour bool) {
	str := fmt.Sprintf("%v", content)
	if !dispColour {
		str = fmt.Sprintf("WARN %s: %s\n", pref, str)
		slog.Printf(str)
	} else {
		logger.Warn(str)
	}
}

// Print out info messages
func Info(content interface{}, dispColour bool) {
	str := fmt.Sprintf("%v", content)
	if !dispColour {
		str = fmt.Sprintf("INFO %s: %s\n", pref, str)
		slog.Printf(str)
	} else {
		logger.Info(str)
	}
}
