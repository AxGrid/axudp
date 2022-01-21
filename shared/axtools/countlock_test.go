package axtools

import (
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestCountLock(t *testing.T) {
	cl := NewCountLock(3)
	ch := make(chan int)
	var ai uint32 = 0
	for i := 0; i < 10; i++ {
		go func(ii int) {
			cl.Lock()
			defer cl.Unlock()
			atomic.AddUint32(&ai, 1)
			time.Sleep(time.Millisecond * 100)
			log.Info().Msg("done")
			ch <- ii
		}(i)
	}

	time.Sleep(time.Millisecond * 10)
	if !assert.Equal(t, ai, uint32(3)) {
		t.Fatal()
	}

	for i := 0; i < 10; i++ {
		log.Debug().Msgf("Done %d", <-ch)
	}
	assert.Equal(t, ai, uint32(10))
}
