> The Russian localization of this README can be found in the file: `README.ru.md`

# barelog

**barelog** is a minimal, fast, and dependency-free logger for Go. It is suitable for libraries, microservices, and any Go application where you want fast logging without extra dependencies.

---

## Installation

```bash
go get github.com/buraev/barelog@latest
```

---

## Quick Start

```go
package main

import "github.com/buraev/barelog"

func main() {
    barelog.Init() // Initialize the logger

    barelog.Info("Starting application", "port", 8080)
    barelog.Debug("Debugging info", "trace_id")
}
```

---

## Log Levels

- `DEBUG`
- `INFO`
- `WARN`
- `ERROR`

Log levels are set globally via `barelog.Init()`. Example:

```go
barelog.Debug("This is a debug message")
barelog.Info("User login", "user", "alice")
barelog.Warn("Attempt", "attempt", 2)
barelog.Error("Connection refused", "err", "connection refused")
```

---

## Contextual Logging

barelog supports middleware, request-scoped logging, etc. For example:

```go
ctx := barelog.WithContext(context.Background(), barelog.New(barelog.DEBUG))
log := barelog.FromContext(ctx)
log.Info("Request started", "req_id")
```

---

## Custom Log Levels

You can define your own log levels if needed:

| Level           | Description             | Code Example             |
|-----------------|------------------------|--------------------------|
| `BARELOG_LEVEL` | Custom log level        | `debug`, `info`          |

Example:

```bash
BARELOG_LEVEL=debug ./yourApp
```

---

## License

MIT

---

[github.com/buraev](https://github.com/buraev)
