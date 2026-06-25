// Package logger provides a small, opinionated application logger built on
// pterm, exposed as a log/slog logger with a swappable package-level default.
//
// It is intentionally separate from log/slog's global default logger, which
// root.go discards to silence noisy library logs (e.g. blaze's per-query
// ranking output). This package owns its own *slog.Logger instance, so its
// records — command status and MCP tool-call tracing — remain visible.
//
// Records render through pterm via a small slog handler that sorts attributes
// by key (for stable, trackable column order) and applies per-key colors. Pass
// structured key/value pairs (and slog.Group for nested sets like tool
// input/output) rather than pre-serialized blobs.
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sort"
	"sync"

	"github.com/pterm/pterm"
)

// Logger is the application logger type. It is a log/slog logger backed by a
// pterm slog handler.
type Logger = slog.Logger

// keyStyles assigns a distinct color per top-level log key so they stay
// visually trackable regardless of position.
var keyStyles = map[string]pterm.Style{
	"input":   *pterm.NewStyle(pterm.FgMagenta),
	"output":  *pterm.NewStyle(pterm.FgYellow),
	"elapsed": *pterm.NewStyle(pterm.FgCyan),
	"err":     *pterm.NewStyle(pterm.FgRed, pterm.Bold),
}

// Config controls logger construction. The zero value is valid and yields an
// info-level logger writing to stderr.
type Config struct {
	// Writer is the log sink. Defaults to os.Stderr when nil.
	Writer io.Writer
	// Level is the minimum level to emit: "debug", "info", "warn", or "error".
	// Defaults to "info" when empty or unrecognized.
	Level string
}

// New builds a configured Logger from cfg.
func New(cfg Config) *Logger {
	w := cfg.Writer
	if w == nil {
		w = os.Stderr
	}

	pl := pterm.DefaultLogger.
		WithWriter(w).
		WithLevel(parseLevel(cfg.Level)).
		WithKeyStyles(keyStyles)

	return slog.New(&sortedHandler{pl: pl})
}

// sortedHandler is a slog.Handler that renders through a pterm.Logger, sorting
// attributes by key first so the output column order is stable. pterm's own
// NewSlogHandler flattens attributes into a map, which randomizes their order;
// this handler preserves a deterministic (sorted) order instead.
type sortedHandler struct {
	pl    *pterm.Logger
	attrs []slog.Attr
}

func (h *sortedHandler) Enabled(_ context.Context, level slog.Level) bool {
	switch level {
	case slog.LevelDebug:
		return h.pl.CanPrint(pterm.LogLevelDebug)
	case slog.LevelWarn:
		return h.pl.CanPrint(pterm.LogLevelWarn)
	case slog.LevelError:
		return h.pl.CanPrint(pterm.LogLevelError)
	default:
		return h.pl.CanPrint(pterm.LogLevelInfo)
	}
}

func (h *sortedHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make([]slog.Attr, 0, record.NumAttrs()+len(h.attrs))
	attrs = append(attrs, h.attrs...)
	record.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})

	sort.SliceStable(attrs, func(i, j int) bool { return attrs[i].Key < attrs[j].Key })

	args := make([]pterm.LoggerArgument, len(attrs))
	for i, a := range attrs {
		args[i] = pterm.LoggerArgument{Key: a.Key, Value: a.Value}
	}
	wrapped := [][]pterm.LoggerArgument{args}

	switch record.Level {
	case slog.LevelDebug:
		h.pl.Debug(record.Message, wrapped...)
	case slog.LevelWarn:
		h.pl.Warn(record.Message, wrapped...)
	case slog.LevelError:
		h.pl.Error(record.Message, wrapped...)
	default:
		h.pl.Info(record.Message, wrapped...)
	}
	return nil
}

func (h *sortedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	merged := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	merged = append(merged, h.attrs...)
	merged = append(merged, attrs...)
	return &sortedHandler{pl: h.pl, attrs: merged}
}

// WithGroup is a no-op: this app passes groups as record attributes (rendered
// as values), not via handler-level grouping.
func (h *sortedHandler) WithGroup(string) slog.Handler { return h }

// parseLevel maps a level string to a pterm.LogLevel, defaulting to info.
func parseLevel(level string) pterm.LogLevel {
	switch level {
	case "debug":
		return pterm.LogLevelDebug
	case "warn":
		return pterm.LogLevelWarn
	case "error":
		return pterm.LogLevelError
	default:
		return pterm.LogLevelInfo
	}
}

var (
	mu  sync.RWMutex
	def = New(Config{})
)

// Default returns the package-level logger.
func Default() *Logger {
	mu.RLock()
	defer mu.RUnlock()
	return def
}

// SetDefault replaces the package-level logger. Safe to call during startup
// once flags are parsed (e.g. to set the level or writer).
func SetDefault(l *Logger) {
	mu.Lock()
	defer mu.Unlock()
	def = l
}

// Package-level helpers delegate to the Default logger, mirroring log/slog.

// Debug logs at debug level via the default logger.
func Debug(msg string, args ...any) { Default().Debug(msg, args...) }

// Info logs at info level via the default logger.
func Info(msg string, args ...any) { Default().Info(msg, args...) }

// Warn logs at warn level via the default logger.
func Warn(msg string, args ...any) { Default().Warn(msg, args...) }

// Error logs at error level via the default logger.
func Error(msg string, args ...any) { Default().Error(msg, args...) }

// With returns a sub-logger of the default logger carrying args on every record.
func With(args ...any) *Logger { return Default().With(args...) }
