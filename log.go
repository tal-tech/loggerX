package logger

import (
	"context"

	"github.com/tal-tech/loggerX/builders"
	"github.com/tal-tech/loggerX/logutils"
	"github.com/tal-tech/loggerX/plugin"
)

//global builer
var builder MessageBuilder = new(builders.DefaultBuilder)

//builer interface
type MessageBuilder interface {
	//log build
	LoggerX(ctx context.Context, lvl string, tag string, args interface{}, v ...interface{})
	//log format
	Build(ctx context.Context, args interface{}, v ...interface{}) (position, message string)
}

//builder select
func SetBuilder(b MessageBuilder) {
	builder = b
}

//INFO level
func I(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "INFO", tag, args, v...)
}

//INFO with context
func Ix(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "INFO", tag, args, v...)
}

//TRACE level
func T(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "TRACE", tag, args, v...)
}

//TRACE level with context
func Tx(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "TRACE", tag, args, v...)
}

//DEBUG level
func D(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "DEBUG", tag, args, v...)
}

//DEBUG level with context
func Dx(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "DEBUG", tag, args, v...)
}

//WARNING level
func W(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "WARNING", tag, args, v...)
}

//WARNING level with context
func Wx(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "WARNING", tag, args, v...)
}

//ERROR level
func E(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "ERROR", tag, args, v...)
}

//ERROR level with context
func Ex(ctx context.Context, tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(ctx, "ERROR", tag, args, v...)
}

//CRITICAL level
func C(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "CRITICAL", tag, args, v...)
}

//FATAL level
func F(tag string, args interface{}, v ...interface{}) {
	builder.LoggerX(nil, "FATAL", tag, args, v...)
}

//id generate
func Id() int64 {
	return logutils.GenLoggerId()
}

//perf monitor
func RegisterPerfPlugin(perfFunc plugin.PerfPlugin) {
	plugin.PerfPluginer = &perfFunc
	return
}
