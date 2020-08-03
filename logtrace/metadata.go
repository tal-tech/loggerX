package logtrace

import (
	"context"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

//context key
func GetMetadataKey() string {
	return "logtraceutil_metadata_key"
}

//metadata转ctx
func AppendLogTraceMetadataContext(ctx context.Context, metadata map[string]string) context.Context {
	if metadata == nil {
		return ctx
	}
	t := NewTraceNode()
	for k, v := range metadata {
		t.Set(k, v)
	}
	ctx = context.WithValue(ctx, GetMetadataKey(), t)
	return ctx
}

//InitTraceNode
func GenLogTraceMetadata() *TraceNode {
	t := NewTraceNode()
	t.Set("x_trace_id", `"`+NewTraceId()+`"`)
	t.Set("x_rpcid", `"0.1"`)
	t.Set("x_name", `"request"`)
	return t
}

func NewTraceId() string {
	uuid := uuid.NewV4()
	return uuid.String()
}

//Get TraceNode
func ExtractTraceNodeFromXesContext(ctx context.Context) *TraceNode {
	meta := ctx.Value(GetMetadataKey())
	if meta == nil {
		return NewTraceNode()
	} else {
		if val, ok := meta.(*TraceNode); ok {
			return val
		} else {
			return NewTraceNode()
		}
	}
}

//TraceNode add other kv
func InjectMetadata(ctx context.Context, mapPtr *map[string]string) bool {
	meta := ExtractTraceNodeFromXesContext(ctx)
	traceRpcId := meta.Get("x_rpcid")
	if len(traceRpcId) == 0 {
		return false
	}
	for k, v := range meta.ForkMap() {
		(*mapPtr)[k] = v
	}
	return true
}

//Incr RpcId 1.1.1=>1.1.2
func IncrementRpcId(ctx context.Context) bool {
	meta := ExtractTraceNodeFromXesContext(ctx)
	traceRpcId := meta.Get("x_rpcid")
	if len(traceRpcId) == 0 {
		return false
	}

	index := strings.LastIndex(traceRpcId, ".")
	if index == -1 {
		return false
	}
	index += 1

	// skip the last `"` char, traceRpcId is of string format
	id, err := strconv.Atoi(traceRpcId[index : len(traceRpcId)-1])
	if err != nil {
		return false
	}

	id += 1
	meta.Set("x_rpcid", traceRpcId[0:index]+strconv.Itoa(id)+`"`)
	return true
}

//Append RpcId  1.1.1=>1.1.1.0
func AppendNewRpcId(ctx context.Context) bool {
	meta := ExtractTraceNodeFromXesContext(ctx)
	traceRpcId := meta.Get("x_rpcid")
	if len(traceRpcId) == 0 {
		return false
	}

	//	skip the last `"` char
	//	note there is no detection for this char, simply skip it
	traceRpcId = traceRpcId[:len(traceRpcId)-1] + `.0"`
	meta.Set("x_rpcid", traceRpcId)
	return true
}

//Ctx add k,v
func AppendKeyValue(ctx context.Context, key, value string) bool {
	meta := ExtractTraceNodeFromXesContext(ctx)
	traceRpcId := meta.Get("x_rpcid")
	if len(traceRpcId) == 0 {
		return false
	}

	meta.Set(key, value)
	return true
}
