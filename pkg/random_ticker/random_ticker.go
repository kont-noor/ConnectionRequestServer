package randomticker

import (
	"math/rand"
	"sync"
	"time"
)

type RandomTicker struct {
	C           chan time.Time
	stopChannel chan struct{}
	minDuration time.Duration
	maxDuration time.Duration
	closeOnce   sync.Once
}

func New(min, max time.Duration) *RandomTicker {
	ticker := RandomTicker{
		minDuration: min,
		maxDuration: max,
	}

	ticker.start()

	return &ticker
}

func (t *RandomTicker) start() {
	t.C = make(chan time.Time)
	t.stopChannel = make(chan struct{})

	go func() {
		sleep := t.setRandomDuration()
		for {
			select {
			case <-t.stopChannel:
				close(t.C)
				return
			case <-time.After(sleep):
				sleep = t.setRandomDuration()
				t.C <- time.Now()
			}
		}
	}()
}

func (t *RandomTicker) setRandomDuration() time.Duration {
	return time.Duration(rand.Int63n(int64(t.maxDuration-t.minDuration))) + t.minDuration
}

func (t *RandomTicker) Stop() {
	t.closeOnce.Do(func() {
		close(t.stopChannel)
	})
}
