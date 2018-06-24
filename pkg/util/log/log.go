package log

import (
	"fmt"
	"os"

	"gitlab.com/nickwanninger/geode/pkg/util/color"
)

// LogLevel -
type LogLevel int

// Some global constant types for log levels
const (
	LevelInfo LogLevel = iota
	LevelVerbose
	LevelError
)

// LevelMap is a mapping from human names to the integer repr
var LevelMap = map[string]LogLevel{
	"info":    LevelInfo,
	"error":   LevelError,
	"verbose": LevelVerbose,
}

var currentLevel LogLevel
var enabledTags map[string]bool
var enableAll bool

func init() {
	currentLevel = LevelInfo
	enabledTags = make(map[string]bool)
	enableAll = false
}

type anything interface{}

func log(msg string) {
	fmt.Printf("%s", msg)

	f, _ := os.OpenFile("/tmp/geode.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.Write([]byte(msg))
}

// Printf -
func Printf(format string, args ...interface{}) {
	tolog := fmt.Sprintf(format, args...)
	log(tolog)
}

// Debug -
func Debug(format string, args ...interface{}) {
	tolog := color.Green(fmt.Sprintf(format, args...))
	log(tolog)
}

// Info -
func Info(format string, args ...interface{}) {
	tolog := color.Cyan("info: ") + fmt.Sprintf(format, args...)
	log(tolog)
}

// Error -
func Error(format string, args ...interface{}) {
	tolog := color.Red("error: ") + fmt.Sprintf(format, args...)
	log(tolog)
}

// Fatal -
func Fatal(format string, args ...interface{}) {
	tolog := color.Red("fatal: ") + fmt.Sprintf(format, args...)
	log(tolog)
	os.Exit(1)
}
