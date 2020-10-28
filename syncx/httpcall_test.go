package syncx

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttpLimitCall(t *testing.T) {
	var values = make(map[int64]int)
	var lock sync.Mutex
	var wg sync.WaitGroup

	httpLimitCall := NewHttpLimitCall(
		WithHttpCallWorkers(5),
		WithHttpCallContainer(20),
	)
	for i := 0; i < 30; i++ {
		wg.Add(1)
		httpLimitCall.Call(func(item interface{}) {
			defer wg.Done()
			tn := time.Now().Unix()
			lock.Lock()
			values[tn] = values[tn] + 1
			lock.Unlock()
			time.Sleep(time.Second)

		}, i)
	}
	wg.Wait()
	assert.True(t, len(values) == 6)
	for _, v := range values {
		assert.Equal(t, v, 5)
	}
}
