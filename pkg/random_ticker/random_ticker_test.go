package randomticker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNormalWork(t *testing.T) {
	assert := assert.New(t)

	ticker := New(1*time.Millisecond, 9*time.Millisecond)
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

func TestStop(t *testing.T) {
	assert := assert.New(t)

	ticker := New(2*time.Millisecond, 9*time.Millisecond)
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

	assert.Equal(count, 0, "Ticker must not tick after it's stop")
}

func TestRandomTicks(t *testing.T) {
	assert := assert.New(t)

	ticker := New(1*time.Millisecond, 100*time.Millisecond)
	defer ticker.Stop()

	tickPeriod := 40 * time.Millisecond
	lastTick := time.Now()

	count := 0
	for timeStamp := range ticker.C {
		count++
		if count >= 5 {
			break
		}

		newPeriod := timeStamp.Sub(lastTick)

		assert.NotEqual(newPeriod.Milliseconds(), tickPeriod.Milliseconds(), "periods must be random")

		lastTick = timeStamp
		tickPeriod = newPeriod
	}
}

func TestSameTicks(t *testing.T) {
	assert := assert.New(t)

	ticker := New(2*time.Millisecond, 2*time.Millisecond)
	defer ticker.Stop()

	tickPeriod := 2 * time.Millisecond
	lastTick := time.Now()

	count := 0
	for timeStamp := range ticker.C {
		count++
		if count >= 20 {
			break
		}

		newPeriod := timeStamp.Sub(lastTick)

		assert.Equal(newPeriod.Milliseconds(), tickPeriod.Milliseconds(), "periods must be the same")

		lastTick = timeStamp
		tickPeriod = newPeriod
	}
}
