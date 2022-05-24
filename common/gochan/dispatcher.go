package gochan

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sync/atomic"
)

const (
	dispatcherStatusOpen   int32 = 0
	dispatcherStatusClosed int32 = 1
)

// Dispatcher pool of gochan
type Dispatcher struct {
	status int32
	gcNum  int
	gcs    []*gochan
}

// NewDispatcher new pool
// @gochanNum the number of gochan
// @bufferNum the buffer number of chan in each gochan
func NewDispatcher(gochanNum, bufferNum int) *Dispatcher {
	log.Infof("%d gochans and %d bufferNum chan buffer", gochanNum, bufferNum)
	d := &Dispatcher{
		gcNum:  gochanNum,
		gcs:    make([]*gochan, gochanNum),
		status: dispatcherStatusOpen,
	}
	for index := range d.gcs {
		gc := newGochan(bufferNum)
		gc.setUUID(index)
		d.gcs[index] = gc
		gc.run()
	}
	return d
}

// Dispatch dispatch task referenced by objID
func (d *Dispatcher) Dispatch(objID int, task DispatchItem) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("dispatcher closed")
		}
	}()

	// dispatch at random if objID is less than 0
	if objID < 0 {
		objID = rand.Intn(d.gcNum)
	}
	// dispatching to closed channel is limited
	if atomic.LoadInt32(&d.status) == dispatcherStatusClosed {
		return errors.New("dispatcher closed")
	}

	index := objID % d.gcNum
	d.gcs[index].tasksChan <- task
	return
}

// Close close diapatcher
func (d *Dispatcher) Close() {
	if atomic.LoadInt32(&d.status) == dispatcherStatusClosed {
		return
	}

	atomic.StoreInt32(&d.status, dispatcherStatusClosed)
	// close gochan
	for _, gc := range d.gcs {
		gc.dieChan <- struct{}{}
	}
}
