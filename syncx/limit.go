package syncx

import (
	"errors"

	"github.com/binwen/zero-tools/lang"
)

var ErrLimitReturn = errors.New("discarding limited token, resource pool is full, someone returned multiple times")

type Limit struct {
	pool chan lang.PlaceholderType
}

func NewLimit(n int) Limit {
	return Limit{pool: make(chan lang.PlaceholderType, n)}
}

func (l Limit) Borrow() {
	l.pool <- lang.Placeholder
}

func (l Limit) Return() error {
	select {
	case <-l.pool:
		return nil
	default:
		return ErrLimitReturn
	}
}

func (l Limit) TryBorrow() bool {
	select {
	case l.pool <- lang.Placeholder:
		return true
	default:
		return false
	}
}
