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

func (l *Logger) log(level Level, args ...any) {
	if level < l.level {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	tag := fmt.Sprintf("%s%-5s%s", levelColors[level], levelNames[level], resetColor)
	msg := strings.TrimSpace(fmt.Sprintln(args...)) // пробелы + без \n
	line := fmt.Sprintf("%s [%s] %s", tag, timestamp, msg)
	fmt.Fprintln(l.out, line)
}

// --- Уровневые методы ---

func (l *Logger) Debug(args ...any) { l.log(DEBUG, args...) }
func (l *Logger) Info(args ...any)  { l.log(INFO, args...) }
func (l *Logger) Warn(args ...any)  { l.log(WARN, args...) }
func (l *Logger) Error(args ...any) { l.log(ERROR, args...) }

// --- Глобальные обёртки ---

func SetGlobal(logger *Logger) {
	if logger != nil {
		globalLogger = logger
	}
}

func Debug(args ...any) { globalLogger.Debug(args...) }
func Info(args ...any)  { globalLogger.Info(args...) }
func Warn(args ...any)  { globalLogger.Warn(args...) }
func Error(args ...any) { globalLogger.Error(args...) }

// --- Поддержка контекста ---

func WithContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, ctxKey, logger)
}

func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(ctxKey).(*Logger); ok {
		return logger
	}
	return globalLogger
}

// --- Init из окружения ---

func Init() {
	levelStr := strings.ToLower(os.Getenv("BARELOG_LEVEL"))
	level := INFO

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
		// default
	default:
		fmt.Fprintf(os.Stderr, "barelog: неизвестный уровень: %q, используем INFO\n", levelStr)
	}

	SetGlobal(New(level))
}
