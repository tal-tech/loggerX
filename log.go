package logger

import (
	"context"

	"github.com/tal-tech/loggerX/builders"
	"github.com/tal-tech/loggerX/logutils"
	"github.com/tal-tech/loggerX/plugin"
)

var builder MessageBuilder = new(builders.DefaultBuilder)

type MessageBuilder interface {
	LoggerX(ctx context.Context, lvl string, tag string, args interface{}, v ...interface{})
	Build(ctx context.Context, args interface{}, v ...interface{}) (position, message string)
}

func SetBuilder(b MessageBuilder) {
	builder = b
}

func I(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "INFO", tag, args, v...)
}
func Ix(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "INFO", tag, args, v...)
}

func T(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "TRACE", tag, args, v...)
}
func Tx(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "TRACE", tag, args, v...)
}

func D(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "DEBUG", tag, args, v...)
}

func Dx(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "DEBUG", tag, args, v...)
}

func W(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "WARNING", tag, args, v...)
}

func Wx(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "WARNING", tag, args, v...)
}

func E(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "ERROR", tag, args, v...)
}

func Ex(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "ERROR", tag, args, v...)
}

func C(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "CRITICAL", tag, args, v...)
}

func F(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "FATAL", tag, args, v...)
}

func Id() int64 {
	return logutils.GenLoggerId()
}

func RegisterPerfPlugin(perfFunc plugin.PerfPlugin) {
	plugin.PerfPluginer = &perfFunc
	return
}
