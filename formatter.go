package logx

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"time"
)

// LogEntry represents a log entry structure
type LogEntry struct {
	Time       time.Time `json:"time" xml:"time"`       // Time when the log occurred
	Level      Level     `json:"level" xml:"level"`     // Log level (e.g., TRACE, INFO, ERROR, etc.)
	Prefix     string    `json:"prefix" xml:"prefix"`   // Log prefix for distinguishing modules or subsystems, can be empty
	File       string    `json:"file" xml:"file"`       // File path where the log is located (relative or formatted path)
	Line       int       `json:"line" xml:"line"`       // Line number in the file where the log is located
	Message    string    `json:"message" xml:"message"` // Log message content
	CallerSkip int       `json:"-" xml:"-"`             // Stack depth for determining the call source location (file and line number)
}

// Formatter defines a function type for formatting log entries
// Input: Log entry
// Output: Formatted log string
type Formatter func(entry LogEntry) []byte

var DefaultFormatter Formatter = func(entry LogEntry) []byte {
	// Time format
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	// Log level in uppercase
	level := entry.Level.String()

	fileLine := fmt.Sprintf("[%s:%d]", TrimCallerPath(entry.File, 1), entry.Line)
	// Log prefix
	prefix := ""
	if entry.Prefix != "" {
		prefix = entry.Prefix + ": "
	}
	// Custom default output format
	logStr := fmt.Sprintf("%s %s %s %s%s",
		timestamp,
		entry.Level.Color().Sprint(level),
		color.New(color.FgHiBlack).Sprint(fileLine),
		color.New(color.FgHiBlack).Add(color.Bold).Sprint(prefix),
		entry.Level.Color().Sprint(entry.Message))
	return []byte(logStr + "\n")
}

func TrimCallerPath(path string, n int) string {
	// lovely borrowed from zap
	// nb. To make sure we trim the path correctly on Windows too, we
	// counter-intuitively need to use '/' and *not* os.PathSeparator here,
	// because the path given originates from Go stdlib, specifically
	// runtime.Caller() which (as of Mar/17) returns forward slashes even on
	// Windows.
	//
	// See https://github.com/golang/go/issues/3335
	// and https://github.com/golang/go/issues/18151
	//
	// for discussion on the issue on Go side.
	// Return the full path if n is 0.
	if n <= 0 {
		return path
	}
	// Find the last separator.
	idx := strings.LastIndexByte(path, '/')
	if idx == -1 {
		return path
	}
	for i := 0; i < n-1; i++ {
		// Find the penultimate separator.
		idx = strings.LastIndexByte(path[:idx], '/')
		if idx == -1 {
			return path
		}
	}
	return path[idx+1:]
}
