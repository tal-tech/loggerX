package logtrace

import "sync"

type TraceNode struct {
	metadata map[string]string
	lock     *sync.RWMutex
}

func NewTraceNode() *TraceNode {
	t := new(TraceNode)
	t.metadata = make(map[string]string, 5)
	t.lock = new(sync.RWMutex)
	return t
}

func (this *TraceNode) Get(key string) string {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.metadata[key]
}

func (this *TraceNode) Set(key, val string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.metadata[key] = val
	return
}

func (this *TraceNode) ForkMap() map[string]string {
	ret := make(map[string]string, 5)
	this.lock.RLock()
	defer this.lock.RUnlock()
	for k, v := range this.metadata {
		ret[k] = v
	}
	return ret
}
