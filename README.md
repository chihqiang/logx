# logx

`logx` 是一个Go语言实现的简单、灵活的日志模块，提供了多级日志支持和彩色输出功能。

## 功能特点

- 支持4个基础日志级别（Debug、Info、Warn、Error）
- 支持自定义日志级别（可在基础级别上进行偏移）
- 彩色日志输出，不同级别显示不同颜色
- 灵活的日志配置（可设置输出目标、前缀、格式化函数等）
- 全局日志实例，方便快速使用
- 完全并发安全的设计
- 支持自定义日志格式化器

## 文件结构

- `level.go`: 定义日志级别枚举和相关方法
- `log.go`: 提供全局日志实例和简化的日志函数
- `logger.go`: 定义Logger结构体和核心实现
- `formatter.go`: 定义日志格式化器接口和默认实现

## 日志级别

logx 支持以下基础日志级别：

|   级别   | 枚举值 |   颜色    |           描述           |
| :------: | :----: | :-------: | :----------------------: |
|  Error   |   8    | 亮红色加粗 |  错误信息，影响功能运行  |
|   Warn   |   4    |  粗黄色   |      警告，潜在问题      |
|   Info   |   0    |    绿色   |       普通运行信息       |
|  Debug   |  -4    |    蓝色   |     调试信息，最详细     |

**自定义级别**：logx 支持在基础级别上进行偏移，例如 `LevelInfo+1` 或 `LevelWarn-1`，这样可以实现更细粒度的日志控制。自定义级别的字符串表示会在基础级别后加上偏移量，如 `INFO+1`、`WARN-1` 等。

## 使用示例

### 全局日志实例

```go
import "github.com/chihqiang/logx"

// 使用不同级别的日志
logx.Info("这是一条信息日志")
logx.Warn("这是一条警告日志")
logx.Error("这是一条错误日志")
logx.Debug("这是一条调试日志")

// 设置全局日志输出目标
logx.SetOutput(os.Stdout)

// 设置全局日志前缀
logx.SetPrefix("[APP]")

// 使用自定义级别
logx.Log(logx.LevelInfo+1, "这是一条自定义级别的信息日志")
```

### 创建自定义日志实例

```go
package main

import (
	"github.com/chihqiang/logx"
	"os"
)

func main() {
	// 创建自定义日志实例，输出到标准错误
	logger := logx.New(os.Stderr)
	
	// 设置日志前缀
	logger.SetPrefix("[MYAPP]")
	
	// 使用自定义日志实例
	logger.Info("这是使用自定义日志实例的信息")
	logger.Debug("这是调试信息")
	logger.Warn("这是警告信息")
	logger.Error("这是错误信息")
	
	// 使用自定义日志级别
	logger.Log(logx.LevelDebug-1, "这是比Debug更详细的日志")
}
```

### 自定义日志格式化器

```go
package main

import (
	"encoding/json"
	"github.com/chihqiang/logx"
	"os"
)

func main() {
	logger := logx.New(os.Stdout)
	// 设置自定义格式化器
	logger.SetFormatter(func(entry logx.LogEntry) []byte {
		marshal, _ := json.Marshal(entry)
		return marshal
	})
	// 输出日志
	logger.Info("使用自定义格式化器的日志")
}
```

## 核心API

### 全局日志函数

- `Debug(format string, v ...any)` - 记录Debug级别日志
- `Info(format string, v ...any)` - 记录Info级别日志  
- `Warn(format string, v ...any)` - 记录Warn级别日志
- `Error(format string, v ...any)` - 记录Error级别日志
- `Log(level Level, format string, v ...any)` - 记录指定级别的日志
- `SetOutput(w io.Writer)` - 设置日志输出目标
- `SetPrefix(p string)` - 设置日志前缀
- `SetFormatter(fn Formatter)` - 设置日志格式化函数

### Logger结构体方法

- `New(w io.Writer) *Logger` - 创建新的日志实例
- `(*Logger) SetOutput(w io.Writer)` - 设置日志输出目标
- `(*Logger) SetPrefix(p string)` - 设置日志前缀
- `(*Logger) SetFormatter(fn Formatter)` - 设置日志格式化函数
- `(*Logger) Debug(format string, v ...any)` - 输出Debug级别日志
- `(*Logger) Info(format string, v ...any)` - 输出Info级别日志
- `(*Logger) Warn(format string, v ...any)` - 输出Warn级别日志
- `(*Logger) Error(format string, v ...any)` - 输出Error级别日志
- `(*Logger) Log(level Level, format string, v ...any) error` - 输出指定级别的日志

## 依赖

- `github.com/fatih/color`: 提供终端彩色输出功能
- Go标准库 `fmt`, `io`, `os`, `runtime`, `sync`, `time`

## 注意事项

- 所有日志方法都是线程安全的，可以在并发环境中使用
- 默认日志格式化器会显示时间戳、日志级别、调用文件和行号以及日志内容
- 可以通过自定义Formatter实现完全个性化的日志格式
- 日志输出目标可以是任意实现了io.Writer接口的对象，如标准输出、文件等
