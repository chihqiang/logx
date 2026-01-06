package logx

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
)

// BenchmarkLoggerParallel benchmarks Logger performance under high concurrency
func BenchmarkLoggerParallel(b *testing.B) {
	buffer := &bytes.Buffer{}
	logger := New(buffer)

	// Reset timer
	b.ResetTimer()

	// Parallel test
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark message: %d", time.Now().UnixNano())
		}
	})
}

// BenchmarkGlobalLoggerParallel benchmarks global logger performance under high concurrency
func BenchmarkGlobalLoggerParallel(b *testing.B) {
	buffer := &bytes.Buffer{}

	// Reset global logger
	std = nil
	stdOnce = sync.Once{}
	SetOutput(buffer)

	// Reset timer
	b.ResetTimer()

	// Parallel test
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Info("benchmark message: %d", time.Now().UnixNano())
		}
	})
}

// BenchmarkLoggerSingle benchmarks Logger performance in single-threaded scenario
func BenchmarkLoggerSingle(b *testing.B) {
	buffer := &bytes.Buffer{}
	logger := New(buffer)

	// Reset timer
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message: %d", i)
	}
}

// BenchmarkGlobalLoggerSingle benchmarks global logger performance in single-threaded scenario
func BenchmarkGlobalLoggerSingle(b *testing.B) {
	buffer := &bytes.Buffer{}

	// Reset global logger
	std = nil
	stdOnce = sync.Once{}
	SetOutput(buffer)

	// Reset timer
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Info("benchmark message: %d", i)
	}
}

// safeWriter is a thread-safe writer
// It uses mutex to protect internal bytes.Buffer
// Used for high concurrency testing
type safeWriter struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

// Write implements io.Writer interface, thread-safe
func (w *safeWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}

// String returns the content of internal buffer
func (w *safeWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}

// TestLoggerThroughput tests log library throughput (logs/second)
func TestLoggerThroughput(t *testing.T) {
	safeBuf := &safeWriter{}
	logger := New(safeBuf)

	const goroutines = 100
	const logsPerGoroutine = 1000
	const totalLogs = goroutines * logsPerGoroutine

	var wg sync.WaitGroup
	wg.Add(goroutines)

	start := time.Now()

	// Start multiple goroutines for concurrent logging
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < logsPerGoroutine; j++ {
				logger.Info("throughput test: goroutine %d, log %d", id, j)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	elapsed := time.Since(start)
	throughput := float64(totalLogs) / elapsed.Seconds()

	t.Logf("Total logs: %d", totalLogs)
	t.Logf("Total time: %v", elapsed)
	t.Logf("Throughput: %.2f logs/second", throughput)

	// Verify output
	output := safeBuf.String()
	lines := strings.Count(output, "\n")
	if lines != totalLogs {
		t.Errorf("Expected %d log lines, got %d", totalLogs, lines)
	}
}

// TestGlobalLoggerThroughput tests global logger throughput
func TestGlobalLoggerThroughput(t *testing.T) {
	safeBuf := &safeWriter{}

	// Reset global logger
	std = nil
	stdOnce = sync.Once{}
	SetOutput(safeBuf)

	const goroutines = 100
	const logsPerGoroutine = 1000
	const totalLogs = goroutines * logsPerGoroutine

	var wg sync.WaitGroup
	wg.Add(goroutines)

	start := time.Now()

	// Start multiple goroutines for concurrent logging
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < logsPerGoroutine; j++ {
				Info("throughput test: goroutine %d, log %d", id, j)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	elapsed := time.Since(start)
	throughput := float64(totalLogs) / elapsed.Seconds()

	t.Logf("Global logger - Total logs: %d", totalLogs)
	t.Logf("Global logger - Total time: %v", elapsed)
	t.Logf("Global logger - Throughput: %.2f logs/second", throughput)

	// Verify output
	output := safeBuf.String()
	lines := strings.Count(output, "\n")
	if lines != totalLogs {
		t.Errorf("Expected %d log lines, got %d", totalLogs, lines)
	}
}
