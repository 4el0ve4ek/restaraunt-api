package order

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/log"
)

func NewProcessor(logger log.Logger, orderRepository orderRepository, dishesGetter dishesGetter) *processor {
	ctx, cancel := context.WithCancel(context.Background())
	return &processor{
		orderRepository: orderRepository,
		dishesGetter:    dishesGetter,

		ctx:    ctx,
		cancel: cancel,
		logger: logger,

		ticker: time.NewTicker(time.Second * 10),
	}
}

type processor struct {
	orderRepository orderRepository
	dishesGetter    dishesGetter

	cancel context.CancelFunc
	ctx    context.Context
	logger log.Logger

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
			err := p.tryTakeOrder()
			if err != nil {
				p.logger.Warn(errors.Wrap(err, "process order failure"))
			}
		}
	}
}

func (p *processor) tryTakeOrder() error {
	waitingOrders, err := p.orderRepository.GetWaitingOrdersID(p.ctx)
	if err != nil {
		return errors.Wrap(err, "get waiting orders")
	}
	if len(waitingOrders) == 0 {
		return nil
	}

	processedOrder, err := p.orderRepository.GetOrderByID(p.ctx, waitingOrders[0])
	if err != nil {
		return errors.Wrap(err, "get order by id")
	}

	err = p.orderRepository.SetProcessingOrderByID(p.ctx, processedOrder.ID)
	if err != nil {
		return errors.Wrap(err, "set status processing")
	}

	dishToQuantity := make(map[int]int, len(processedOrder.Dishes))
	for _, orderDish := range processedOrder.Dishes {
		dishToQuantity[orderDish.DishID] = orderDish.Quantity
	}

	dishes, err := p.dishesGetter.GetAllDishes(p.ctx)
	if err != nil {
		return errors.Wrap(err, "get all dishes")
	}

	needCancel := false
	for _, dish := range dishes {
		if !dish.Available || dish.Quantity < 0 {
			continue
		}
		required := dishToQuantity[dish.ID]
		if required <= dish.Quantity {
			continue
		}
		needCancel = true
	}
	if needCancel {
		return errors.Wrap(p.orderRepository.SetCancelOrderByID(p.ctx, processedOrder.ID), " set cancel order")
	}

	time.Sleep(30 * time.Second)
	err = p.orderRepository.SetSuccessOrderByID(p.ctx, processedOrder.ID)
	if err != nil {
		return errors.Wrap(err, "set success order")
	}

	for dishID, quantity := range dishToQuantity {
		p.orderRepository.ReduceDishQuantity(p.ctx, dishID, quantity)
	}
	return nil
}
