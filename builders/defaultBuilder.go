package builders

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/tal-tech/loggerX/log4go"
	"github.com/tal-tech/loggerX/logutils"
	"github.com/tal-tech/loggerX/plugin"
	"github.com/tal-tech/loggerX/stackerr"
)

//default builder
type DefaultBuilder struct {
}

var hostname string

//logger function
//support 7 level log
//use context to transfer loggerid and start time
func (this *DefaultBuilder) LoggerX(ctx context.Context, lvl string, tag string, args interface{}, v ...interface{}) {
	if !logutils.ValidLevel(lvl) {
		return
	}
	if tag == "" {
		tag = "NoTagError"
	}

	if ctx == nil {
		//generate new loggerid
		id := strconv.FormatInt(logutils.GenLoggerId(), 10)
		ctx = context.WithValue(context.Background(), "logid", id)
	}

	tag = logutils.Filter(tag)
	//build logger
	position, message := this.Build(ctx, args, v...)

	//compute the logtime from request acceived to now
	if startValue := ctx.Value("start"); startValue != nil {
		if start, ok := startValue.(time.Time); ok {
			cost := time.Now().Sub(start)
			message = message + " COST:" + fmt.Sprintf("%.2f", cost.Seconds()*1e3)
		}
	}

	//log level match
	switch lvl {
	case "DEBUG":
		log4go.Log(log4go.DEBUG, position, tag+"\t"+message)
	case "TRACE":
		log4go.Log(log4go.TRACE, position, tag+"\t"+message)
	case "INFO":
		log4go.Log(log4go.INFO, position, tag+"\t"+message)
	case "WARNING":
		log4go.Log(log4go.WARNING, position, tag+"\t"+message)
	case "ERROR":
		log4go.Log(log4go.ERROR, position, tag+"\t"+message)
		//ERROR level have 5 sub levels,1-5, level 5 no need care
		if !strings.Contains(message, "[level[5]]") {
			//monitor error log in falcon
			plugin.DoPerfPlugin(tag)
		}
	case "CRITICAL":
		log4go.Log(log4go.CRITICAL, position, tag+"\t"+message)
		//monitor critical log in falcon
		plugin.DoPerfPlugin(tag)
	case "FATAL":
		log4go.Log(log4go.CRITICAL, position, tag+"\t"+message)
		//monitor fatal error in falcon
		plugin.DoPerfPlugin(tag)
		//fatal error need suspend the request
		panic(message)
	}
}

//logger build
//used to customize the log format
func (this *DefaultBuilder) Build(ctx context.Context, args interface{}, v ...interface{}) (position string, message string) {
	id := ctx.Value("logid")
	logid := cast.ToString(id)
	//type match and get message
	switch t := args.(type) {
	case *stackerr.StackErr:
		message = t.Info
		position = t.Position
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

	//get the log position in which file and which row
	if position == "" {
		_, file, line, ok := runtime.Caller(3)
		if ok {
			position = filepath.Base(file) + ":" + strconv.Itoa(line) + ":" + logid
		} else {
			position = "EMPTY"
		}
	}

	if hostname == "" {
		//get server hostname
		hostname, _ = os.Hostname()
	}

	position = position + "\t" + hostname
	message = logutils.Filter(message)
	return
}
