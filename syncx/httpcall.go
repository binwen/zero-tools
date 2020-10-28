package syncx

import (
	"sync"

	"github.com/binwen/zero-tools/lang"
	"github.com/binwen/zero-tools/threading"
)

const (
	defaultWorkers    = 10
	defaultContainers = 1024
)

type (
	LimitOption func(options *limitOptions)

	HttpCallFunc func(item interface{})
	callWrap     struct {
		fn  HttpCallFunc // 请求处理方法
		req interface{}  // 请求参数
	}

	HttpLimitCall struct {
		pool      chan lang.PlaceholderType // 并发限制
		container chan callWrap             // 存放请求体的容器
		lock      sync.Mutex
		guarded   bool
	}

	limitOptions struct {
		workers    int
		containers int
	}
)

func newLimitOptions() limitOptions {
	return limitOptions{
		workers:    defaultWorkers,
		containers: defaultContainers,
	}
}

func WithHttpCallWorkers(workers int) LimitOption {
	return func(options *limitOptions) {
		options.workers = workers
	}
}

func WithHttpCallContainer(containers int) LimitOption {
	return func(options *limitOptions) {
		options.containers = containers
	}
}

func NewHttpLimitCall(opts ...LimitOption) *HttpLimitCall {
	options := newLimitOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &HttpLimitCall{
		pool:      make(chan lang.PlaceholderType, options.workers),
		container: make(chan callWrap, options.containers),
	}
}

// fn为回调方法 reqData是fn的参数
func (hc *HttpLimitCall) Call(fn HttpCallFunc, reqData interface{}) {
	hc.lock.Lock()
	defer func() {
		var start bool
		if !hc.guarded {
			hc.guarded = true
			start = true
		}
		hc.lock.Unlock()
		if start {
			hc.backgroundCall()
		}
	}()

	hc.container <- callWrap{fn: fn, req: reqData}
}

func (hc *HttpLimitCall) backgroundCall() {
	threading.GoSafe(func() {
		var wg sync.WaitGroup
		for {
			hc.pool <- lang.Placeholder
			callWrap, ok := <-hc.container
			if !ok {
				<-hc.pool
				hc.lock.Lock()
				hc.guarded = false
				hc.lock.Unlock()
				break
			}
			wg.Add(1)
			threading.GoSafe(func() {
				defer func() {
					wg.Done()
					<-hc.pool
				}()
				callWrap.fn(callWrap.req)
			})
		}
		wg.Wait()
	})
}

func (hc *HttpLimitCall) Surplus() int {
	return len(hc.container)
}
