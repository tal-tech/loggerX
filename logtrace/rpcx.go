package logtrace

import (
	"context"

	"github.com/smallnest/rpcx/share"
	"github.com/tal-tech/xtools/jsutil"
)

var rpcxMetaKey = share.ReqMetaDataKey

func InjectTraceNodeToRpcx(ctx context.Context) context.Context {
	//IncrementRpcId(ctx)
	meta := ExtractTraceNodeFromXesContext(ctx)
	traceRpcId := meta.Get("x_rpcid")
	if len(traceRpcId) == 0 {
		return ctx
	}
	metaStr, err := jsutil.Json.MarshalToString(meta.ForkMap())
	if err != nil {
		metaStr = ""
	}

	rpcxMeta := ctx.Value(rpcxMetaKey)
	if rpcxMeta == nil {
		//todo: why enter here
		logMeta := map[string]string{
			GetMetadataKey(): metaStr,
		}
		ctx = context.WithValue(ctx, rpcxMetaKey, logMeta)
	} else {
		if logMeta, ok := rpcxMeta.(map[string]string); ok {
			logMeta[GetMetadataKey()] = metaStr
			ctx = context.WithValue(ctx, rpcxMetaKey, logMeta)
		}
	}

	return ctx
}

func ExtractTraceNodeToXexContext(ctx context.Context) context.Context {
	reqMeta, ok := ctx.Value(rpcxMetaKey).(map[string]string)
	if ok {
		logMetaKey := GetMetadataKey()
		metaStr, ok1 := reqMeta[logMetaKey]
		if ok1 {
			var mapVal map[string]string
			if jsutil.Json.UnmarshalFromString(metaStr, &mapVal) == nil {
				tNode := ExtractTraceNodeFromXesContext(ctx)
				for k, v := range mapVal {
					tNode.Set(k, v)
				}
				ctx = context.WithValue(ctx, GetMetadataKey(), tNode)
			}
			AppendNewRpcId(ctx)
		}
	}
	return ctx
}
