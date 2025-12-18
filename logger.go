package logx

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type ILogger interface {
	SetOutput(w io.Writer)
	SetPrefix(prefix string)
	SetFormatter(fn Formatter)
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Log(level Level, format string, v ...interface{}) error
}

// New creates a new Logger instance
// Parameter w specifies the log output destination (can be os.Stdout, os.Stderr, file, etc.)
func New(w io.Writer) *Logger {
	l := &Logger{}
	l.SetOutput(w)
	l.SetFormatter(DefaultFormatter) // Use default formatter function
	return l
}

// Logger represents a logging object
type Logger struct {
	mu         sync.RWMutex // Read-write lock for concurrent safety
	writer     io.Writer    // Log output destination
	prefix     string       // Log prefix
	formatter  Formatter    // Log formatting function
	callerSkip int          // runtime.Caller level offset for correctly displaying call file and line number
}

// SetOutput sets the log output destination (thread-safe)
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writer = w
}

// SetPrefix sets the log prefix (thread-safe)
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// SetFormatter sets the log formatting function (thread-safe)
func (l *Logger) SetFormatter(fn Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.formatter = fn
}

// Debug outputs Debug level logs
func (l *Logger) Debug(format string, v ...interface{}) {
	_ = l.log(LevelDebug, format, v...)
}

// Info outputs Info level logs
func (l *Logger) Info(format string, v ...interface{}) {
	_ = l.log(LevelInfo, format, v...)
}

// Warn outputs Warn level logs
func (l *Logger) Warn(format string, v ...interface{}) {
	_ = l.log(LevelWarn, format, v...)
}

// Error outputs Error level logs
func (l *Logger) Error(format string, v ...interface{}) {
	_ = l.log(LevelError, format, v...)
}

func (l *Logger) Log(level Level, format string, v ...interface{}) error {
	return l.log(level, format, v...)
}

// Log outputs logs at the specified level
// 1. Get call file and line number based on callDepth
// 2. Format log entry using Formatter
// 3. Write to log output destination (writer), default to os.Stdout if writer is nil
func (l *Logger) log(level Level, format string, v ...interface{}) error {
	// Read Logger current state with concurrent safety
	l.mu.RLock()
	prefix := l.prefix
	formatter := l.formatter
	writer := l.writer
	callerSkip := l.callerSkip
	if callerSkip == 0 {
		callerSkip = 2
	}
	l.mu.RUnlock()
	// Format log content
	msg := fmt.Sprintf(format, v...)
	// Get call file and line number
	_, file, line, ok := runtime.Caller(callerSkip)
	if !ok {
		file = "???" // Placeholder when unable to obtain
		line = 0
	}
	// Output log, default to stdout if writer is nil
	if writer == nil {
		writer = os.Stdout
	}
	_, err := writer.Write(formatter(LogEntry{
		Time:       time.Now(),
		Level:      level,
		Prefix:     prefix,
		CallerSkip: callerSkip,
		File:       file,
		Line:       line,
		Message:    msg,
	}))
	return err
}
