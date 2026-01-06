package logx

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	// Test creating Logger instance
	buffer := &bytes.Buffer{}
	logger := New(buffer)

	if logger == nil {
		t.Fatal("Expected a Logger instance, got nil")
	}

	// Default output should be os.Stderr, but New function can accept parameters
	logger.SetOutput(buffer)
	logger.Info("test message")

	output := buffer.String()
	if !strings.Contains(output, "test message") {
		t.Fatal("Expected log message to contain 'test message', got:", output)
	}
}

func TestSetOutput(t *testing.T) {
	// Test setting output target
	buffer1 := &bytes.Buffer{}
	buffer2 := &bytes.Buffer{}

	logger := New(buffer1)
	logger.Info("message to buffer1")

	if !strings.Contains(buffer1.String(), "message to buffer1") {
		t.Fatal("Expected message in buffer1")
	}

	if buffer2.Len() != 0 {
		t.Fatal("Expected buffer2 to be empty")
	}

	// Change output target
	logger.SetOutput(buffer2)
	logger.Info("message to buffer2")

	if !strings.Contains(buffer2.String(), "message to buffer2") {
		t.Fatal("Expected message in buffer2 after SetOutput")
	}
}

func TestSetPrefix(t *testing.T) {
	// Test setting log prefix
	buffer := &bytes.Buffer{}
	logger := New(buffer)

	// No prefix
	logger.Info("no prefix")
	buffer.Reset()

	// Set prefix
	logger.SetPrefix("TEST")
	logger.Info("with prefix")
	output2 := buffer.String()

	if !strings.Contains(output2, "TEST:") {
		t.Fatal("Expected log message to contain prefix 'TEST:', got:", output2)
	}
}

func TestSetFormatter(t *testing.T) {
	// Test setting custom formatter
	buffer := &bytes.Buffer{}
	logger := New(buffer)

	// Use custom formatter
	customFormatter := func(entry LogEntry) []byte {
		return []byte("CUSTOM:" + entry.Message + "\n")
	}

	logger.SetFormatter(customFormatter)
	logger.Info("custom format test")

	output := buffer.String()
	if !strings.Contains(output, "CUSTOM:custom format test") {
		t.Fatal("Expected custom formatted message, got:", output)
	}
}

func TestLogLevels(t *testing.T) {
	// Test different log levels
	buffer := &bytes.Buffer{}
	logger := New(buffer)

	// Reset custom formatter to use default formatter for testing levels
	logger.SetFormatter(DefaultFormatter)

	// Test Info level
	buffer.Reset()
	logger.Info("info message")
	if !strings.Contains(buffer.String(), "INFO") {
		t.Fatal("Expected INFO level in log")
	}

	// Test Debug level
	buffer.Reset()
	logger.Debug("debug message")
	if !strings.Contains(buffer.String(), "DEBUG") {
		t.Fatal("Expected DEBUG level in log")
	}

	// Test Warn level
	buffer.Reset()
	logger.Warn("warn message")
	if !strings.Contains(buffer.String(), "WARN") {
		t.Fatal("Expected WARN level in log")
	}

	// Test Error level
	buffer.Reset()
	logger.Error("error message")
	if !strings.Contains(buffer.String(), "ERROR") {
		t.Fatal("Expected ERROR level in log")
	}
}

func TestLogEntryFields(t *testing.T) {
	// Test if LogEntry fields are set correctly
	var capturedEntry LogEntry

	customFormatter := func(entry LogEntry) []byte {
		capturedEntry = entry
		return []byte("captured")
	}

	buffer := &bytes.Buffer{}
	logger := New(buffer)
	logger.SetFormatter(customFormatter)

	// Send test log
	logger.Info("test message")

	// Verify Time field
	if capturedEntry.Time.IsZero() {
		t.Fatal("Expected Time field to be set")
	}

	// Verify Level field
	if capturedEntry.Level != LevelInfo {
		t.Fatalf("Expected LevelInfo, got %v", capturedEntry.Level)
	}

	// Verify Message field
	if capturedEntry.Message != "test message" {
		t.Fatalf("Expected message 'test message', got '%s'", capturedEntry.Message)
	}

	// Verify File and Line fields (should be available)
	if capturedEntry.File == "???" || capturedEntry.Line == 0 {
		t.Fatalf("Expected valid file and line, got file=%s, line=%d", capturedEntry.File, capturedEntry.Line)
	}
}

func TestCallerSkip(t *testing.T) {
	// Test if callerSkip works correctly
	buffer := &bytes.Buffer{}
	logger := New(buffer)

	// Reset formatter to use default formatter to view file information
	logger.SetFormatter(DefaultFormatter)

	// Create a wrapper function to increase call depth
	wrapper := func() {
		logger.Info("test from wrapper")
	}

	// Save original callerSkip
	originalSkip := logger.callerSkip

	// Test default callerSkip
	buffer.Reset()
	wrapper()
	output1 := buffer.String()

	// Modify callerSkip
	logger.callerSkip = 3
	buffer.Reset()
	wrapper()
	output2 := buffer.String()

	// Restore original settings
	logger.callerSkip = originalSkip

	// The two outputs should show different file/line numbers (one from wrapper function, one from test function)
	if output1 == output2 {
		t.Fatal("Expected different file/line information with different callerSkip values")
	}
}
