package builders

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tal-tech/loggerX/logutils"
	"github.com/tal-tech/loggerX/stackerr"
)

type logrusBuilder struct {
	LogrusI *logrus.Logger
}

func NewLogrusBuilder(logger *logrus.Logger) *logrusBuilder {
	return &logrusBuilder{LogrusI: logger}
}

func (this *logrusBuilder) LoggerX(ctx context.Context, lvl string, tag string, args interface{}, v ...interface{}) {

	if tag == "" {
		tag = "NoTagError"
	}

	tag = logutils.Filter(tag)
	_, message := this.Build(ctx, args, v...)

	field := this.LogrusI.WithFields(logrus.Fields{
		"tag": tag,
	})
	switch lvl {
	case "DEBUG":
		field.Debug(message)
	case "TRACE":
		field.Trace(message)
	case "INFO":
		field.Info(message)
	case "WARNING":
		field.Warn(message)
	case "ERROR":
		field.Error(message)
	case "FATAL":
		field.Panic(message)
	}
}

func (this *logrusBuilder) Build(ctx context.Context, args interface{}, v ...interface{}) (position string, message string) {

	switch t := args.(type) {
	case *stackerr.StackErr:
		message = t.Info
	case error:
		message = t.Error()
	case string:
		if len(v) > 0 {
			message = fmt.Sprintf(t, v...)
		} else {
			message = t
		}
	default:
		message = fmt.Sprint(t)
	}
	message = logutils.Filter(message)
	return
}
