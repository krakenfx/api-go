package book

import (
	"math"
	"sort"
	"sync"
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
	mux        sync.RWMutex
}

// NewLevel constructs a new [Level] struct with default values.
func NewLevel() *Level {
	return &Level{
		orders: make(map[string]*Order),
	}
}

func (l *Level) update(opts *UpdateOptions) {
	if opts.ID == "" {
		l.orders = make(map[string]*Order)
		l.queue = nil
		l.Quantity = opts.Quantity
	} else {
		order, ok := l.orders[opts.ID]
		if !ok && opts.Quantity.Sign() == 1 {
			order = &Order{}
			order.ID = opts.ID
			order.LimitPrice = l.Price
			order.Level = l
			l.orders[opts.ID] = order
		}
		if order != nil {
			order.Quantity = opts.Quantity
			order.Timestamp = opts.Timestamp
		}
		if opts.Quantity.Sign() <= 0 {
			delete(l.orders, opts.ID)
		}
		totalQuantity := decimal.NewFromInt64(0)
		for _, order := range l.orders {
			totalQuantity = totalQuantity.
				SetScale(int64(math.Max(float64(totalQuantity.GetScale()), float64(order.Quantity.GetScale())))).
				Add(order.Quantity)
		}
		l.Quantity = totalQuantity
	}
	l.Timestamp = opts.Timestamp
	l.queueDirty.Store(true)
}

// Queue returns a list of orders arranged by time priority.
func (l *Level) Queue() []*Order {
	l.mux.Lock()
	defer l.mux.Unlock()
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

// GetPriceString returns the level's price.
func (l *Level) GetPriceString() string {
	l.mux.RLock()
	defer l.mux.RUnlock()
	return l.Price.String()
}

// GetQuantityString returns the level's total quantity.
func (l *Level) GetQuantityString() string {
	l.mux.RLock()
	defer l.mux.RUnlock()
	return l.Quantity.String()
}

// Order contains information regarding a limit order.
type Order struct {
	ID         string           `json:"id,omitempty"`
	LimitPrice *decimal.Decimal `json:"limitPrice,omitempty"`
	Quantity   *decimal.Decimal `json:"quantity,omitempty"`
	Timestamp  time.Time        `json:"timestamp,omitempty"`
	Level      *Level           `json:"-"`
}
