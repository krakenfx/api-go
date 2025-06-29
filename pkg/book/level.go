package book

import (
	"sort"
	"sync/atomic"
	"time"

	"github.com/krakenfx/api-go/v2/pkg/decimal"
)

// Level contains price level information.
type Level struct {
	Price      *decimal.Decimal `json:"price,omitempty"`
	Quantity   *decimal.Decimal `json:"quantity,omitempty"`
	Timestamp  time.Time        `json:"time,omitempty"`
	Lower      *Level           `json:"-"`
	Higher     *Level           `json:"-"`
	orders     map[string]*Order
	queue      []*Order
	queueDirty atomic.Bool
}

// NewLevel constructs a new [Level] struct with default values.
func NewLevel() *Level {
	return &Level{
		orders: make(map[string]*Order),
	}
}

func (l *Level) update(opts *UpdateOptions) {
	l.Price = opts.Price.Copy()
	if len(opts.ID) == 0 {
		l.orders = make(map[string]*Order)
		l.queue = nil
		l.Quantity = opts.Quantity.Copy()
	} else {
		order, orderExisted := l.orders[opts.ID]
		if orderExisted {
			l.Quantity = l.Quantity.Add(opts.Quantity.Sub(order.Quantity))
			order.Quantity = opts.Quantity.Copy()
			order.Timestamp = opts.Timestamp
		} else if opts.Quantity.Sign() == 1 {
			l.orders[opts.ID] = &Order{
				ID:         opts.ID,
				LimitPrice: opts.Price.Copy(),
				Quantity:   opts.Quantity.Copy(),
				Timestamp:  opts.Timestamp,
				Level:      l,
			}
			if l.Quantity == nil {
				l.Quantity = opts.Quantity.Copy()
			} else {
				l.Quantity = l.Quantity.Add(opts.Quantity)
			}
		}
		if opts.Quantity.Sign() <= 0 {
			delete(l.orders, opts.ID)
			l.Quantity = l.Quantity.Sub(opts.Quantity)
		}
	}
	l.Timestamp = opts.Timestamp
	l.queueDirty.Store(true)
}

// Queue returns a list of orders arranged by time priority.
func (l *Level) Queue() []*Order {
	if !l.queueDirty.Load() {
		return l.queue
	}
	queue := make([]*Order, len(l.orders))
	var i int
	for _, o := range l.orders {
		queue[i] = o
		i++
	}
	sort.Slice(queue, func(i, j int) bool {
		return queue[i].Timestamp.Before(queue[j].Timestamp)
	})
	l.queue = queue
	l.queueDirty.Store(false)
	return l.queue
}

// Order contains information regarding a limit order.
type Order struct {
	ID         string           `json:"id,omitempty"`
	LimitPrice *decimal.Decimal `json:"limitPrice,omitempty"`
	Quantity   *decimal.Decimal `json:"quantity,omitempty"`
	Timestamp  time.Time        `json:"timestamp,omitempty"`
	Level      *Level           `json:"-"`
}
