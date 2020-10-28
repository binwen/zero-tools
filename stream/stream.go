package stream

import (
	"sort"
	"sync"

	"github.com/binwen/zero-tools/lang"
	"github.com/binwen/zero-tools/threading"
)

/**
流数据处理利器
*/
const (
	defaultWorkers = 16
	minWorkers     = 1
)

type (
	rxOption struct {
		unlimitedWorkers bool
		workers          int
	}

	FilterFunc   func(item interface{}) bool
	ForAllFunc   func(pipe <-chan interface{})
	ForEachFunc  func(item interface{})
	GenerateFunc func(source chan<- interface{})
	KeyFunc      func(item interface{}) interface{}
	LessFunc     func(a, b interface{}) bool
	MapFunc      func(item interface{}) interface{}
	Option       func(opts *rxOption)
	ParallelFunc func(item interface{})
	ReduceFunc   func(pipe <-chan interface{}) (interface{}, error)
	WalkFunc     func(item interface{}, pipe chan<- interface{})

	Stream struct {
		source <-chan interface{}
	}
)

func Range(source <-chan interface{}) Stream {
	return Stream{
		source: source,
	}
}

// 产生数据流 Stream
func From(generate GenerateFunc) Stream {
	source := make(chan interface{})
	threading.GoSafe(func() {
		defer close(source)
		generate(source)
	})

	return Range(source)
}

func Just(items ...interface{}) Stream {
	source := make(chan interface{}, len(items))
	for _, item := range items {
		source <- item
	}
	close(source)
	return Range(source)
}

func (s Stream) Buffer(n int) Stream {
	if n < 0 {
		n = 0
	}
	source := make(chan interface{}, n)
	go func() {
		for item := range s.source {
			source <- item
		}
		close(source)
	}()

	return Range(source)
}

func (s Stream) Distinct(fn KeyFunc) Stream {
	source := make(chan interface{})

	threading.GoSafe(func() {
		defer close(source)

		keys := make(map[interface{}]lang.PlaceholderType)
		for item := range s.source {
			key := fn(item)
			if _, ok := keys[key]; !ok {
				source <- item
				keys[key] = lang.Placeholder
			}
		}
	})
	return Range(source)
}

func (s Stream) Filter(fn FilterFunc, opts ...Option) Stream {
	return s.Walk(func(item interface{}, pipe chan<- interface{}) {
		if fn(item) {
			pipe <- item
		}
	}, opts...)
}

func (s Stream) Walk(fn WalkFunc, opts ...Option) Stream {
	option := buildOptions(opts...)
	if option.unlimitedWorkers {
		return s.walkUnlimited(fn, option)
	} else {
		return s.walkLimited(fn, option)
	}
}

func (s Stream) walkLimited(fn WalkFunc, option *rxOption) Stream {
	pipe := make(chan interface{}, option.workers)

	go func() {
		var wg sync.WaitGroup
		pool := make(chan lang.PlaceholderType, option.workers)
		for {
			pool <- lang.Placeholder
			item, ok := <-s.source
			if !ok {
				<-pool
				break
			}
			wg.Add(1)
			threading.GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()
				fn(item, pipe)
			})
		}
		wg.Wait()
		close(pipe)
	}()

	return Range(pipe)
}

func (s Stream) walkUnlimited(fn WalkFunc, option *rxOption) Stream {
	pipe := make(chan interface{}, defaultWorkers)
	go func() {
		var wg sync.WaitGroup
		for {
			item, ok := <-s.source
			if !ok {
				break
			}
			wg.Add(1)
			threading.GoSafe(func() {
				defer wg.Done()
				fn(item, pipe)
			})
		}
		wg.Wait()
		close(pipe)
	}()
	return Range(pipe)
}

func (s Stream) Group(fn KeyFunc) Stream {
	groups := make(map[interface{}][]interface{})
	for item := range s.source {
		key := fn(item)
		groups[key] = append(groups[key], item)
	}
	source := make(chan interface{})
	go func() {
		for _, group := range groups {
			source <- group
		}
		close(source)
	}()

	return Range(source)
}

func (s Stream) Head(n int64) Stream {
	if n < 1 {
		panic("n must be greater than 0")
	}
	source := make(chan interface{})
	go func() {
		for item := range s.source {
			n--
			if n >= 0 {
				source <- item
			}
			if n == 0 {
				close(source)
			}
		}
		if n > 0 {
			close(source)
		}
	}()
	return Range(source)
}

func (s Stream) Split(n int) Stream {
	if n < 1 {
		panic("n should be greater than 0")
	}
	source := make(chan interface{})
	go func() {
		var chunk []interface{}
		for item := range s.source {
			chunk = append(chunk, item)
			if len(chunk) == n {
				source <- chunk
				chunk = nil
			}
		}
		if chunk != nil {
			source <- chunk
		}
		close(source)
	}()

	return Range(source)
}

func (s Stream) Map(fn MapFunc, opts ...Option) Stream {
	return s.Walk(func(item interface{}, pipe chan<- interface{}) {
		pipe <- fn(item)
	}, opts...)
}

func (s Stream) Merge() Stream {
	var items []interface{}
	for item := range s.source {
		items = append(items, item)
	}
	source := make(chan interface{}, 1)
	source <- items
	close(source)

	return Range(source)
}

func (s Stream) Reverse() Stream {
	var items []interface{}
	for item := range s.source {
		items = append(items, item)
	}
	for i := len(items)/2 - 1; i >= 0; i-- {
		opp := len(items) - 1 - i
		items[i], items[opp] = items[opp], items[i]
	}

	return Just(items...)
}

func (s Stream) Sort(less LessFunc) Stream {
	var items []interface{}
	for item := range s.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})

	return Just(items...)
}

func (s Stream) Done() {
	for range s.source {
	}
}

func (s Stream) ForAll(fn ForAllFunc) {
	fn(s.source)
}

func (s Stream) ForEach(fn ForEachFunc) {
	for item := range s.source {
		fn(item)
	}
}

func (s Stream) Count() (count int) {
	for range s.source {
		count++
	}
	return
}

func (s Stream) Parallel(fn ParallelFunc, opts ...Option) {
	s.Walk(func(item interface{}, pipe chan<- interface{}) {
		fn(item)
	}, opts...).Done()
}

func (s Stream) Reduce(fn ReduceFunc) (interface{}, error) {
	return fn(s.source)
}

func UnlimitedWorkers() Option {
	return func(opts *rxOption) {
		opts.unlimitedWorkers = true
	}
}

func WithWorkers(workers int) Option {
	return func(opts *rxOption) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

func newOptions() *rxOption {
	return &rxOption{
		workers: defaultWorkers,
	}
}

func buildOptions(opts ...Option) *rxOption {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}
