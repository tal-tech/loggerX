package example

import (
	"context"
	"strconv"
	"time"

	logger "github.com/tal-tech/loggerX"
	"github.com/tal-tech/loggerX/builders"
	"github.com/tal-tech/loggerX/logtrace"
)

func example_tracelog() {
	//初始化logger main函数指定
	config := logger.NewLogConfig()
	logger.InitLogWithConfig(config)
	//或使用xml配置 logger.InitLogger("conf/log.xml")
	defer logger.Close()
	builder := new(builders.TraceBuilder)
	builder.SetTraceDepartment("HS-Golang")
	builder.SetTraceVersion("0.1")
	logger.SetBuilder(builder)

	//初始化trace信息 一次完整调用前执行
	ctx := context.WithValue(context.Background(), "logid", strconv.FormatInt(logger.Id(), 10))
	ctx = context.WithValue(ctx, "start", time.Now())
	ctx = context.WithValue(ctx, logtrace.GetMetadataKey(), logtrace.GenLogTraceMetadata())

	//logger
	logger.Ix(ctx, "Example", "example log time:%v,module:%s", time.Now(), "test1")
	logger.Ex(ctx, "Example", "example log time:%v,module:%s", time.Now(), "test2")
	logger.Wx(ctx, "Example", "example log time:%v,module:%s", time.Now(), "test3")

}
