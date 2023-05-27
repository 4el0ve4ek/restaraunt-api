package order

import (
	"context"
	"time"
)

func NewProcessor(orderRepository orderRepository) *processor {
	ctx, cancel := context.WithCancel(context.Background())
	return &processor{
		ctx:    ctx,
		cancel: cancel,

		ticker: time.NewTicker(time.Second * 10),
	}
}

type processor struct {
	orderRepository orderRepository

	cancel context.CancelFunc
	ctx    context.Context

	ticker *time.Ticker
}

func (p *processor) Close() error {
	p.cancel()
	p.ticker.Stop()
	return nil
}

func (p *processor) Run() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-p.ticker.C:
			p.tryTakeOrder()
		}
	}
}

func (p *processor) tryTakeOrder() {
	// select from db where status = waiting
	// get first
	// if not enough quantity cancel otherwise minus quantity

	// time.Sleep(30 * time.Second)
	// update status
}
