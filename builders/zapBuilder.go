package builders

import (
	"context"
	"fmt"

	"github.com/tal-tech/loggerX/logutils"
	"github.com/tal-tech/loggerX/stackerr"
	"go.uber.org/zap"
)

type zapBuilder struct {
	zapI *zap.Logger
}

func NewZapBuilder(logger *zap.Logger) *zapBuilder {
	return &zapBuilder{zapI: logger}
}

func (this *zapBuilder) LoggerX(ctx context.Context, lvl string, tag string, args interface{}, v ...interface{}) {

	if tag == "" {
		tag = "NoTagError"
	}

	tag = logutils.Filter(tag)
	_, message := this.Build(ctx, args, v...)

	suger := this.zapI.Sugar()
	switch lvl {
	case "DEBUG":
		suger.Debugw(message, "tag", tag)
	case "TRACE":
		suger.Debugw(message, "tag", tag)
	case "INFO":
		suger.Infow(message, "tag", tag)
	case "WARNING":
		suger.Warnw(message, "tag", tag)
	case "ERROR":
		suger.Errorw(message, "tag", tag)
	case "FATAL":
		suger.Panicw(message, "tag", tag)
	}
}

func (this *zapBuilder) Build(ctx context.Context, args interface{}, v ...interface{}) (position string, message string) {

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
