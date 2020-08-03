package example

import (
	"context"
	"strconv"
	"time"

	"github.com/tal-tech/loggerX/builders"
	"github.com/tal-tech/loggerX/logtrace"
	logger "github.com/wgffgw/loggerX"
)

func example_main() {
	//初始化logger main函数指定
	config := logger.NewLogConfig()
	logger.InitLogWithConfig(config)
	//或使用xml配置 logger.InitLogger("conf/log.xml")
	defer logger.Close()
	builder := new(builders.TraceBuilder)
	builder.SetTraceDepartment("HS-Golang")
	builder.SetTraceVersion("0.1")
	logger.SetBuilder(builder)
}

type Args struct {
}

type Reply struct {
}

func example_server_func(ctx context.Context, args *Args, reply *Reply) error {
	ctx = logtrace.ExtractTraceNodeToXexContext(ctx)
	logger.Ix(ctx, "Server", "start call %v", args)
	//Do func
	logger.Ix(ctx, "Server", "finish call %v", args)
	return nil
}

func example_client_func() {

	//初始化trace信息 一次完整调用前执行
	ctx := context.WithValue(context.Background(), "logid", strconv.FormatInt(logger.Id(), 10))
	ctx = context.WithValue(ctx, "start", time.Now())
	ctx = context.WithValue(ctx, logtrace.GetMetadataKey(), logtrace.GenLogTraceMetadata())

	logger.Ix(ctx, "Client", "before call %v", "args")
	ctx = logtrace.InjectTraceNodeToRpcx(ctx)
	//rpc call
	//xclient.Call(ctx, "CallFunc", &args, &reply)
	logger.Ix(ctx, "Client", "after call %v", "args")

}
