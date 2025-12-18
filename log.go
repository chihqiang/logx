package logx

import (
	"io"
	"os"
	"sync"
)

var (
	std     *Logger   // Global default Logger instance
	stdOnce sync.Once // Ensures global Logger is initialized only once (thread-safe)
)

// _std returns the global Logger instance (singleton pattern)
// Initializes Logger on first call and sets callDepth
func _std() *Logger {
	stdOnce.Do(func() {
		// Create a Logger that outputs to standard error
		_log := New(os.Stderr)
		_log.callerSkip = 3
		std = _log
	})
	return std
}

// SetOutput sets the output destination for the global Logger (thread-safe)
func SetOutput(w io.Writer) {
	_std().SetOutput(w)
}

// SetPrefix sets the log prefix for the global Logger (thread-safe)
func SetPrefix(p string) {
	_std().SetPrefix(p)
}

// SetFormatter sets the log formatting function for the global Logger (thread-safe)
func SetFormatter(fn Formatter) {
	_std().SetFormatter(fn)
}

// Debug logs at Debug level
func Debug(format string, v ...interface{}) {
	_std().Debug(format, v...)
}

// Info logs at Info level
func Info(format string, v ...interface{}) {
	_std().Info(format, v...)
}

// Warn logs at Warn level
func Warn(format string, v ...interface{}) {
	_std().Warn(format, v...)
}

// Error logs at Error level
func Error(format string, v ...interface{}) {
	_std().Error(format, v...)
}

// Log logs at the specified Level
// Logs will be output if the level is higher than the Logger's minimum level
func Log(level Level, format string, v ...interface{}) error {
	return _std().Log(level, format, v...)
}
