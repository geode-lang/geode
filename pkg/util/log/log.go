package log

import (
	"fmt"
	"os"
	"time"

	"github.com/geode-lang/geode/pkg/util/color"
)

// ShowTimers determines if the compiler should show timers or not
var ShowTimers = false

// PrintVerbose determinies if the compiler should show non-error/warning messages
// like info and debug
var PrintVerbose = false

func log(msg string) {
	fmt.Printf("%s", msg)
}

// Printf -
func Printf(format string, args ...interface{}) {
	tolog := fmt.Sprintf(format, args...)
	log(tolog)
}

// Debug -
func Debug(format string, args ...interface{}) {
	if PrintVerbose {
		tolog := color.Yellow("[debug] ") + fmt.Sprintf(format, args...)
		log(tolog)
	}
}

// Info -
func Info(format string, args ...interface{}) {
	tolog := color.Cyan("[info] ") + fmt.Sprintf(format, args...)
	if PrintVerbose {
		log(tolog)
	}

}

// Deprecated -
func Deprecated(format string, args ...interface{}) {
	tolog := color.Bold("[deprecated] ") + fmt.Sprintf(format, args...)
	log(tolog)
}

// Syntax -
func Syntax(format string, args ...interface{}) {
	tolog := color.Yellow("[syntax] ") + fmt.Sprintf(format, args...)
	log(tolog)
}

// Error -
func Error(format string, args ...interface{}) {
	tolog := color.Red("[error] ") + fmt.Sprintf(format, args...)
	log(tolog)
}

// Fatal -
func Fatal(format string, args ...interface{}) {
	tolog := color.Red("[fatal] ") + fmt.Sprintf(format, args...)
	log(tolog)
	os.Exit(1)
}

// Verbose is a verbose printing style
func Verbose(format string, args ...interface{}) {

	if PrintVerbose {
		tolog := color.Magenta("[verbose] ") + fmt.Sprintf(format, args...)
		log(tolog)
	}

}

// Timed -
func Timed(title string, fn func()) {

	start := time.Now()

	fn()

	duration := time.Since(start)
	if PrintVerbose {
		t := fmt.Sprintf(color.Green("[%s]"), duration)
		Printf("%s %s\n", t, title)
	}

}
