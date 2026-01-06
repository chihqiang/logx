English | [简体中文](README_zh.md)

<div align="center">
<h1>logx</h1>

[![Auth](https://img.shields.io/badge/Auth-chihqiang-ff69b4)](https://github.com/chihqiang)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/chihqiang/logx)](https://github.com/chihqiang/logx/pulls)
[![Go Report Card](https://goreportcard.com/badge/github.com/chihqiang/logx)](https://goreportcard.com/report/github.com/chihqiang/logx)
[![Release](https://img.shields.io/github/release/chihqiang/logx.svg?style=flat-square)](https://github.com/chihqiang/logx/releases)
[![GitHub Pull Requests](https://img.shields.io/github/stars/chihqiang/logx)](https://github.com/chihqiang/logx/stargazers)
[![HitCount](https://views.whatilearened.today/views/github/chihqiang/logx.svg)](https://github.com/chihqiang/logx)
[![GitHub license](https://img.shields.io/github/license/chihqiang/logx)](https://github.com/chihqiang/logx/blob/main/LICENSE)

<p> logx is a simple and flexible logging module implemented in Go, providing multi-level logging support and colored output functionality.</p>

</div>

## Features

- Supports 4 basic log levels (Debug, Info, Warn, Error)
- Supports custom log levels (can be offset based on basic levels)
- Colored log output, different levels display different colors
- Flexible log configuration (can set output target, prefix, formatting function, etc.)
- Global log instance for quick use
- Fully concurrent-safe design
- Support for custom log formatters

## File Structure

- `level.go`: Defines log level enumeration and related methods
- `log.go`: Provides global log instance and simplified log functions
- `logger.go`: Defines Logger struct and core implementation
- `formatter.go`: Defines log formatter interface and default implementation

## Log Levels

logx supports the following basic log levels:

|  Level  | Value |    Color     |          Description          |
| :-----: | :---: | :----------: | :---------------------------: |
|  Error  |   8   | Bold bright red | Error information affecting function operation |
|   Warn  |   4   |   Bold yellow   |       Warning, potential problem       |
|   Info  |   0   |      Green      |       Normal running information       |
|  Debug  |  -4   |      Blue       |     Debug information, most detailed     |

**Custom Levels**: logx supports offsets based on basic levels, such as `LevelInfo+1` or `LevelWarn-1`, which enables more granular log control. The string representation of custom levels will include the offset after the basic level, like `INFO+1`, `WARN-1`, etc.

## Usage Examples

### Global Log Instance

```go
import "github.com/chihqiang/logx"

// Use different log levels
logx.Info("This is an info log")
logx.Warn("This is a warning log")
logx.Error("This is an error log")
logx.Debug("This is a debug log")

// Set global log output target
logx.SetOutput(os.Stdout)

// Set global log prefix
logx.SetPrefix("[APP]")

// Use custom level
logx.Log(logx.LevelInfo+1, "This is a custom level info log")
```

### Create Custom Log Instance

```go
package main

import (
	"github.com/chihqiang/logx"
	"os"
)

func main() {
	// Create custom log instance, output to standard error
	logger := logx.New(os.Stderr)
	
	// Set log prefix
	logger.SetPrefix("[MYAPP]")
	
	// Use custom log instance
	logger.Info("This is information using custom log instance")
	logger.Debug("This is debug information")
	logger.Warn("This is warning information")
	logger.Error("This is error information")
	
	// Use custom log level
	logger.Log(logx.LevelDebug-1, "This is more detailed than Debug level")
}
```

### Custom Log Formatter

```go
package main

import (
	"encoding/json"
	"github.com/chihqiang/logx"
	"os"
)

func main() {
	logger := logx.New(os.Stdout)
	// Set custom formatter
	logger.SetFormatter(func(entry logx.LogEntry) []byte {
		marshal, _ := json.Marshal(entry)
		return marshal
	})
	// Output log
	logger.Info("Log using custom formatter")
}
```

## Core API

### Global Log Functions

- `Debug(format string, v ...any)` - Record Debug level log
- `Info(format string, v ...any)` - Record Info level log  
- `Warn(format string, v ...any)` - Record Warn level log
- `Error(format string, v ...any)` - Record Error level log
- `Log(level Level, format string, v ...any)` - Record log at specified level
- `SetOutput(w io.Writer)` - Set log output target
- `SetPrefix(p string)` - Set log prefix
- `SetFormatter(fn Formatter)` - Set log formatting function

### Logger Struct Methods

- `New(w io.Writer) *Logger` - Create a new log instance
- `(*Logger) SetOutput(w io.Writer)` - Set log output target
- `(*Logger) SetPrefix(p string)` - Set log prefix
- `(*Logger) SetFormatter(fn Formatter)` - Set log formatting function
- `(*Logger) Debug(format string, v ...any)` - Output Debug level log
- `(*Logger) Info(format string, v ...any)` - Output Info level log
- `(*Logger) Warn(format string, v ...any)` - Output Warn level log
- `(*Logger) Error(format string, v ...any)` - Output Error level log
- `(*Logger) Log(level Level, format string, v ...any) error` - Output log at specified level

## Dependencies

- `github.com/fatih/color`: Provides terminal colored output functionality
- Go standard libraries: `fmt`, `io`, `os`, `runtime`, `sync`, `time`

## Performance

### Benchmark Results

| Test Case | Operations per Second | Time per Operation |
| :-------: | :------------------: | :----------------: |
| Logger Parallel | 450,409 ops/sec | 2,609 ns/op |
| Global Logger Parallel | 653,130 ops/sec | 2,865 ns/op |
| Logger Single | 231,024 ops/sec | 4,875 ns/op |
| Global Logger Single | 207,585 ops/sec | 5,038 ns/op |

### Throughput Test

- **Test Environment**: 100 goroutines, each writing 1000 log entries
- **Total Logs**: 100,000 entries
- **Total Time**: ~198 ms
- **Throughput**: **504,000 logs/second**

## Notes

- All log methods are thread-safe and can be used in concurrent environments
- The default log formatter displays timestamp, log level, calling file and line number, and log content
- Completely personalized log formats can be implemented through custom Formatter
- The log output target can be any object that implements the io.Writer interface, such as standard output, files, etc.