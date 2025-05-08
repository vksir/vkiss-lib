package log

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
	"vkiss-lib/pkg/util/errutil"
)

const slogFields = "slog_fields"

var logger = NewLogger("")

func SetLevel(level string) error {
	return logger.SetLevel(level)
}

func With(args ...any) *Logger {
	return logger.With(args...)
}

func Debug(msg string, args ...any) {
	logger.log(context.Background(), slog.LevelDebug, msg, args...)
}

func DebugC(ctx context.Context, msg string, args ...any) {
	logger.log(ctx, slog.LevelDebug, msg, args...)
}

func Info(msg string, args ...any) {
	logger.log(context.Background(), slog.LevelInfo, msg, args...)
}

func InfoC(ctx context.Context, msg string, args ...any) {
	logger.log(ctx, slog.LevelInfo, msg, args...)
}

func Warn(msg string, args ...any) {
	logger.log(context.Background(), slog.LevelWarn, msg, args...)
}

func WarnC(ctx context.Context, msg string, args ...any) {
	logger.log(ctx, slog.LevelWarn, msg, args...)
}

func Error(msg string, args ...any) {
	logger.log(context.Background(), slog.LevelError, msg, args...)
}

func ErrorC(ctx context.Context, msg string, args ...any) {
	logger.log(ctx, slog.LevelError, msg, args...)
}

type Logger struct {
	logger *slog.Logger
	level  *slog.LevelVar
}

func (l *Logger) SetLevel(level string) error {
	err := l.level.UnmarshalText([]byte(level))
	if err != nil {
		return errutil.Wrap(err)
	}
	return nil
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{logger: l.logger.With(args...), level: l.level}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(context.Background(), slog.LevelDebug, msg, args...)
}

func (l *Logger) DebugC(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(context.Background(), slog.LevelInfo, msg, args...)
}

func (l *Logger) InfoC(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(context.Background(), slog.LevelWarn, msg, args...)
}

func (l *Logger) WarnC(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log(context.Background(), slog.LevelError, msg, args...)
}

func (l *Logger) ErrorC(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelError, msg, args...)
}

func (l *Logger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if !l.logger.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.logger.Handler().Handle(ctx, r)
}

func NewLogger(path string) *Logger {
	var w io.Writer = os.Stderr
	if path != "" {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o640)
		cobra.CheckErr(err)
		w = io.MultiWriter(w, file)
	}

	lvl := &slog.LevelVar{}

	h := &ContextHandler{
		slog.NewTextHandler(w, &slog.HandlerOptions{
			AddSource: true,
			Level:     lvl,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.SourceKey {
					source := a.Value.Any().(*slog.Source)
					source.File = source.Function[strings.LastIndex(source.Function, "/")+1:]
				}
				return a
			},
		}),
	}
	lg := slog.New(h)
	return &Logger{logger: lg, level: lvl}
}

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	args, ok := ctx.Value(slogFields).([]any)
	if ok {
		r.Add(args...)
	}
	return h.Handler.Handle(ctx, r)
}

func AppendCtx(ctx context.Context, args ...any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	curArgs, ok := ctx.Value(slogFields).([]any)
	if ok {
		args = append(curArgs, args...)
	}
	return context.WithValue(ctx, slogFields, args)
}

func Init(path string, level string) {
	logger = NewLogger(path)
	err := logger.SetLevel(level)
	cobra.CheckErr(err)
	cobra.WriteStringAndCheck(os.Stderr, fmt.Sprintf("set log level: %s\n", level))
}
