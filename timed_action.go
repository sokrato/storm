package storm

import (
	"errors"
	"time"
)

var ErrTimeout = errors.New("timeout")

func TimedCall(fn func(), timeout time.Duration) bool {
	done := make(chan bool, 1)
	go func() {
		fn()
		done <- true
	}()

	select {
	case <-time.After(timeout):
		return true
	case <-done:
		return false
	}
}
