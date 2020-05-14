package broker_test

import (
	"context"
	"sync"
	"testing"

	"my.pkg/race/broker"
	"my.pkg/race/broker/mocks"

	"github.com/golang/mock/gomock"
	"github.com/tj/assert"
)

type tstBroker struct {
	*broker.Broker
	cfunc context.CancelFunc
	ctx   context.Context
	ctrl  *gomock.Controller
}

func getBroker(t *testing.T) *tstBroker {
	ctx, cfunc := context.WithCancel(context.Background())
	ctrl := gomock.NewController(t)
	return &tstBroker{
		Broker: broker.New(ctx),
		cfunc:  cfunc,
		ctx:    ctx,
		ctrl:   ctrl,
	}
}

func TestRace(t *testing.T) {
	broker := getBroker(t)
	defer broker.Finish()
	sub := mocks.NewMockSub(broker.ctrl)
	cCh, dCh := make(chan interface{}, 1), make(chan struct{})
	vals := []interface{}{1, 2, 3}
	wg := sync.WaitGroup{}
	wg.Add(len(vals))
	sub.EXPECT().Done().Times(len(vals)).Return(dCh)
	sub.EXPECT().C().Times(len(vals)).Return(cCh).Do(func() {
		wg.Done()
	})
	k := broker.Subscribe(sub)
	assert.NotZero(t, k)
	for _, v := range vals {
		broker.Send(v)
	}
	wg.Wait()
	// I've tried to send all 3 values, channels should be safe to close now
	close(dCh)
	// channel had buffer of 1, so first value should be present
	assert.Equal(t, vals[0], <-cCh)
	// other values should be skipped due to default select
	// assert.Equal(t, 0, len(cCh))
	close(cCh)
}

func TestNoRace(t *testing.T) {
	broker := getBroker(t)
	defer broker.Finish()
	sub := mocks.NewMockSub(broker.ctrl)
	cCh, dCh := make(chan interface{}, 1), make(chan struct{})
	vals := []interface{}{1, 2, 3}
	wg := sync.WaitGroup{}
	wg.Add(len(vals))
	sub.EXPECT().Done().Times(len(vals)).Return(dCh)
	sub.EXPECT().C().Times(len(vals)).Return(cCh).Do(func() {
		wg.Done()
	})
	k := broker.Subscribe(sub)
	assert.NotZero(t, k)
	for _, v := range vals {
		broker.Send(v)
	}
	wg.Wait()
	// I've tried to send all 3 values, channels should be safe to close now
	close(dCh)
	// channel had buffer of 1, so first value should be present
	assert.Equal(t, vals[0], <-cCh)
	// other values should be skipped due to default select
	assert.Equal(t, 0, len(cCh))
	// add this line, and data race magically vanishes
	broker.Unsubscribe(k)
	close(cCh)
}

func (b *tstBroker) Finish() {
	b.cfunc()
	b.ctrl.Finish()
}
