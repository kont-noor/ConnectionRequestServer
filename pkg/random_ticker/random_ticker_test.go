package randomticker_test

import (
	randomticker "connection_request_server/pkg/random_ticker"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	assert := assert.New(t)

	ticker := randomticker.New(1*time.Millisecond, 9*time.Millisecond)
	defer ticker.Stop()

	count := 0
	for range ticker.C {
		count++
		if count >= 3 {
			break
		}
	}

	assert.Equal(count, 3, "Must receive 3 ticks")
}

func Test2(t *testing.T) {
	assert := assert.New(t)

	ticker := randomticker.New(2*time.Millisecond, 9*time.Millisecond)
	defer ticker.Stop()

	count := 0
	done := make(chan struct{})

	go func() {
		for range ticker.C {
			count++
			if count > 10 {
				break
			}
		}
		close(done)
	}()

	ticker.Stop()

	select {
	case <-done:
	case <-time.After(50 * time.Millisecond):
	}

	assert.LessOrEqual(count, 0, "Ticker must not tick after it's stop")
}
