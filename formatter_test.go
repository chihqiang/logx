package logx

import (
	"strings"
	"testing"
	"time"
)

func TestDefaultFormatter(t *testing.T) {
	// Test default formatter
	entry := LogEntry{
		Time:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Level:   LevelInfo,
		Prefix:  "TEST",
		File:    "/path/to/file.go",
		Line:    42,
		Message: "test message",
	}

	formatted := string(DefaultFormatter(entry))

	// Verify basic fields exist
	if !strings.Contains(formatted, "2023-01-01 12:00:00") {
		t.Fatal("Expected formatted log to contain correct timestamp")
	}

	if !strings.Contains(formatted, "INFO") {
		t.Fatal("Expected formatted log to contain level 'INFO'")
	}

	if !strings.Contains(formatted, "TEST:") {
		t.Fatal("Expected formatted log to contain prefix 'TEST:'")
	}

	if !strings.Contains(formatted, "file.go:42") {
		t.Fatal("Expected formatted log to contain file and line")
	}

	if !strings.Contains(formatted, "test message") {
		t.Fatal("Expected formatted log to contain message 'test message'")
	}

	if !strings.HasSuffix(formatted, "\n") {
		t.Fatal("Expected formatted log to end with newline")
	}
}

func TestDefaultFormatterWithoutPrefix(t *testing.T) {
	// Test case without prefix
	entry := LogEntry{
		Time:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Level:   LevelInfo,
		Prefix:  "",
		File:    "/path/to/file.go",
		Line:    42,
		Message: "test message",
	}

	formatted := string(DefaultFormatter(entry))

	if strings.Contains(formatted, ":") && strings.Contains(formatted, "test message") {
		// Check for redundant colon
		parts := strings.Split(formatted, "test message")
		if strings.Contains(parts[0], ":") && !strings.Contains(parts[0], "file.go:42") {
			t.Fatal("Expected no prefix separator when prefix is empty")
		}
	}
}

func TestTrimCallerPath(t *testing.T) {
	// Test TrimCallerPath function
	tests := []struct {
		path     string
		n        int
		expected string
	}{
		{"/a/b/c/file.go", 0, "/a/b/c/file.go"},
		{"/a/b/c/file.go", 1, "file.go"},
		{"/a/b/c/file.go", 2, "c/file.go"},
		{"/a/b/c/file.go", 3, "b/c/file.go"},
		{"/a/b/c/file.go", 4, "a/b/c/file.go"},
		{"/a/b/c/file.go", 5, "/a/b/c/file.go"},
		{"file.go", 1, "file.go"},
		{"", 1, ""},
	}

	for _, tt := range tests {
		result := TrimCallerPath(tt.path, tt.n)
		if result != tt.expected {
			t.Errorf("TrimCallerPath(%q, %d) = %q, expected %q", tt.path, tt.n, result, tt.expected)
		}
	}
}

func TestAllLevelsInFormatter(t *testing.T) {
	// Test that all log levels format correctly
	levels := []Level{
		LevelDebug,
		LevelInfo,
		LevelWarn,
		LevelError,
	}

	entry := LogEntry{
		Time:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Prefix:  "",
		File:    "/path/to/file.go",
		Line:    42,
		Message: "test message",
	}

	for _, level := range levels {
		entry.Level = level
		formatted := string(DefaultFormatter(entry))

		if !strings.Contains(formatted, level.String()) {
			t.Errorf("Expected formatted log for %s level to contain %s", level, level.String())
		}
	}
}
