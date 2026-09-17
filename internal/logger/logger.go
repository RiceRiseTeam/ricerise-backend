package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var logger *slog.Logger = nil

type loggerHandler struct {
	mu    *sync.Mutex
	level slog.Level
	out   io.Writer
}

func (h loggerHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h loggerHandler) Handle(_ context.Context, record slog.Record) error {
	var builder strings.Builder
	builder.Grow(128)

	builder.WriteByte('[')
	builder.WriteString(record.Level.String())
	builder.WriteString("] ")

	if h.level <= slog.LevelDebug {
		builder.WriteString(getSourceFile())
		builder.WriteString(" ")
	}

	builder.WriteString(record.Time.Format(time.DateTime))
	builder.WriteString(" ")

	builder.WriteString(record.Message)

	builder.WriteByte('\n')
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := io.WriteString(h.out, builder.String())
	return err
}

func (h loggerHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h loggerHandler) WithGroup(_ string) slog.Handler {
	return h
}

func getSourceFile() string {
	_, file, line, _ := runtime.Caller(5)
	return filepath.Base(file) + ":" + strconv.Itoa(line)
}

func init() {
	if gin.Mode() == gin.DebugMode {
		handler := &loggerHandler{
			mu:    &sync.Mutex{},
			level: slog.LevelDebug,
			out:   os.Stdout,
		}
		logger = slog.New(handler)
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	}
}

func Info(message string, args ...any) {
	logger.Info(fmt.Sprintf(message, args...))
}

func Debug(message string, args ...any) {
	logger.Debug(fmt.Sprintf(message, args...))
}

func Warn(message string, args ...any) {
	logger.Warn(fmt.Sprintf(message, args...))
}

func Error(message string, args ...any) {
	logger.Error(fmt.Errorf(message, args...).Error())
}
