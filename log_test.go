package logx

import (
	"bytes"
	"strings"
	"sync"
	"testing"
)

func TestGlobalLogger(t *testing.T) {
	// Test global log functions
	buffer := &bytes.Buffer{}

	// Reset global logger settings
	std = nil
	stdOnce = sync.Once{}

	// Set output to buffer
	SetOutput(buffer)

	// Test global Info function
	Info("global info message")
	output := buffer.String()

	if !strings.Contains(output, "global info message") {
		t.Fatal("Expected global log message to contain 'global info message', got:", output)
	}
}

func TestGlobalSetOutput(t *testing.T) {
	// Test setting global log output
	buffer1 := &bytes.Buffer{}
	buffer2 := &bytes.Buffer{}

	// Reset global logger
	std = nil
	stdOnce = sync.Once{}

	// Set output to buffer1
	SetOutput(buffer1)
	Info("message to buffer1")

	if !strings.Contains(buffer1.String(), "message to buffer1") {
		t.Fatal("Expected global log in buffer1")
	}

	// Set output to buffer2
	SetOutput(buffer2)
	Info("message to buffer2")

	if !strings.Contains(buffer2.String(), "message to buffer2") {
		t.Fatal("Expected global log in buffer2 after SetOutput")
	}
}

func TestGlobalSetPrefix(t *testing.T) {
	// Test setting global log prefix
	buffer := &bytes.Buffer{}

	// Reset global logger
	std = nil
	stdOnce = sync.Once{}

	SetOutput(buffer)

	// Test no default prefix
	Info("no prefix")
	buffer.Reset()

	// Set prefix
	SetPrefix("GLOBAL")
	Info("with global prefix")
	output2 := buffer.String()

	if !strings.Contains(output2, "GLOBAL:") {
		t.Fatal("Expected global log to contain prefix 'GLOBAL:', got:", output2)
	}
}

func TestAllGlobalLevels(t *testing.T) {
	// Test all global log levels
	buffer := &bytes.Buffer{}

	// Reset global logger
	std = nil
	stdOnce = sync.Once{}

	SetOutput(buffer)

	// Test all levels
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")

	output := buffer.String()

	// Verify all messages exist
	if !strings.Contains(output, "debug message") {
		t.Fatal("Expected debug message in output")
	}
	if !strings.Contains(output, "info message") {
		t.Fatal("Expected info message in output")
	}
	if !strings.Contains(output, "warn message") {
		t.Fatal("Expected warn message in output")
	}
	if !strings.Contains(output, "error message") {
		t.Fatal("Expected error message in output")
	}

	// Verify all level markers exist
	if !strings.Contains(output, "DEBUG") {
		t.Fatal("Expected DEBUG level in output")
	}
	if !strings.Contains(output, "INFO") {
		t.Fatal("Expected INFO level in output")
	}
	if !strings.Contains(output, "WARN") {
		t.Fatal("Expected WARN level in output")
	}
	if !strings.Contains(output, "ERROR") {
		t.Fatal("Expected ERROR level in output")
	}
}

func TestGlobalSetFormatter(t *testing.T) {
	// Test setting global log formatter
	buffer := &bytes.Buffer{}

	// Reset global logger
	std = nil
	stdOnce = sync.Once{}

	SetOutput(buffer)

	// Set custom formatter
	customFormatter := func(entry LogEntry) []byte {
		return []byte("GLOBAL_CUSTOM:" + entry.Message + "\n")
	}

	SetFormatter(customFormatter)
	Info("test with custom formatter")

	output := buffer.String()
	if !strings.Contains(output, "GLOBAL_CUSTOM:test with custom formatter") {
		t.Fatal("Expected global log with custom format, got:", output)
	}
}

func TestGlobalLogFunction(t *testing.T) {
	// Test Log function
	buffer := &bytes.Buffer{}

	// Reset global logger
	std = nil
	stdOnce = sync.Once{}

	SetOutput(buffer)

	// Use Log function to log different levels
	err := Log(LevelInfo, "log function test")
	if err != nil {
		t.Fatal("Expected no error from Log function, got:", err)
	}

	output := buffer.String()
	if !strings.Contains(output, "log function test") {
		t.Fatal("Expected message from Log function, got:", output)
	}
}
