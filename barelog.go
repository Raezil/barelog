package barelog

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	levelNames  = [...]string{"DEBUG", "INFO", "WARN", "ERROR"}
	levelColors = [...]string{"\033[36m", "\033[32m", "\033[33m", "\033[31m"} // cyan, green, yellow, red
	resetColor  = "\033[0m"
	ctxKey      = &contextKey{}

	// глобальный логгер по умолчанию
	globalLogger = New(INFO)
)

type contextKey struct{}

type Logger struct {
	level Level
	out   *os.File
}

func New(level Level) *Logger {
	return &Logger{
		level: level,
		out:   os.Stdout,
	}
}

func (l *Logger) log(level Level, msg string, kv ...any) {
	if level < l.level {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	tag := fmt.Sprintf("%s%-5s%s", levelColors[level], levelNames[level], resetColor)
	line := fmt.Sprintf("%s [%s] %s", tag, timestamp, msg)

	if len(kv) > 0 {
		parts := make([]string, 0, len(kv)/2)
		for i := 0; i < len(kv)-1; i += 2 {
			parts = append(parts, fmt.Sprintf("%v=%v", kv[i], kv[i+1]))
		}
		line += " | " + strings.Join(parts, " ")
	}

	fmt.Fprintln(l.out, line)
}

// Методы уровня логгера
func (l *Logger) Debug(msg string, kv ...any) { l.log(DEBUG, msg, kv...) }
func (l *Logger) Info(msg string, kv ...any)  { l.log(INFO, msg, kv...) }
func (l *Logger) Warn(msg string, kv ...any)  { l.log(WARN, msg, kv...) }
func (l *Logger) Error(msg string, kv ...any) { l.log(ERROR, msg, kv...) }

// --- Глобальный логгер ---

func SetGlobal(logger *Logger) {
	if logger != nil {
		globalLogger = logger
	}
}

func Debug(msg string, kv ...any) { globalLogger.Debug(msg, kv...) }
func Info(msg string, kv ...any)  { globalLogger.Info(msg, kv...) }
func Warn(msg string, kv ...any)  { globalLogger.Warn(msg, kv...) }
func Error(msg string, kv ...any) { globalLogger.Error(msg, kv...) }

// --- Контекстная поддержка ---

func WithContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, ctxKey, logger)
}

func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(ctxKey).(*Logger); ok {
		return logger
	}
	return globalLogger
}

// Init настраивает глобальный логгер из переменных окружения.
// Например: BARELOG_LEVEL=debug
func Init() {
	levelStr := strings.ToLower(os.Getenv("BARELOG_LEVEL"))
	level := INFO // значение по умолчанию

	switch levelStr {
	case "debug":
		level = DEBUG
	case "info":
		level = INFO
	case "warn", "warning":
		level = WARN
	case "error":
		level = ERROR
	case "":
		//default
	default:
		fmt.Fprintf(os.Stderr, "barelog: неизвестный уровень: %q, используем INFO\n", levelStr)
	}

	SetGlobal(New(level))
}
