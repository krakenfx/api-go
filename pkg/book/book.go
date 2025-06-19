package book

import (
	"math"
	"time"

	"github.com/krakenfx/api-go/pkg/callback"
	"github.com/krakenfx/api-go/pkg/decimal"
)

// Order book structure for L2 and L3.
type Book struct {
	// Name of the book.
	Name string

	// Configuration

	// Number of price levels for checksum.
	MaxDepth int `json:"maxDepth,omitempty"`

	// Whether to eliminate book crossings for each update.
	NoBookCrossing bool

	// Whether to trim price levels outside of the max depth.
	// Must be disable for whole books.
	EnableMaxDepth bool

	// Sides

	Bids *Side
	Asks *Side

	// Events

	OnUpdated          *callback.Manager[*UpdateOptions]
	OnBookCrossed      *callback.Manager[*CrossedResult]
	OnMaxDepthExceeded *callback.Manager[*MaxDepthExceededResult]
	OnChecksummed      *callback.Manager[*ChecksumResult]
}

// New constructs a new [Book] struct with default values.
func New() *Book {
	bids := NewSide()
	bids.Direction = Bid
	asks := NewSide()
	asks.Direction = Ask
	return &Book{
		MaxDepth:           1e10,
		NoBookCrossing:     true,
		EnableMaxDepth:     true,
		Bids:               bids,
		Asks:               asks,
		OnUpdated:          callback.NewManager[*UpdateOptions](),
		OnBookCrossed:      callback.NewManager[*CrossedResult](),
		OnMaxDepthExceeded: callback.NewManager[*MaxDepthExceededResult](),
		OnChecksummed:      callback.NewManager[*ChecksumResult](),
	}
}

// Midpoint returns the midpoint of the order book.
func (b *Book) Midpoint() *decimal.Decimal {
	bid, ask := b.BestBid(), b.BestAsk()
	switch {
	case bid != nil && ask != nil:
		return bid.Price.
			Add(ask.Price).
			Mul(decimal.NewFromFloat64(0.5))
	case bid != nil:
		return bid.Price
	case ask != nil:
		return ask.Price
	default:
		return decimal.NewFromInt64(0)
	}
}

// Spread returns the relative difference between the bid-ask price.
func (b *Book) Spread() *decimal.Decimal {
	bid, ask := b.BestBid(), b.BestAsk()
	switch {
	case bid == nil || ask == nil:
		return decimal.NewFromInt64(0)
	default:
		return ask.Price.
			SetScale(int64(math.Max(float64(ask.Price.GetScale()), float64(decimal.DefaultScale)))).
			Sub(bid.Price).
			Div(ask.Price).
			Mul(decimal.NewFromInt64(100))
	}
}

// Whether a bid or an ask.
type BookDirection string

const (
	Bid = "bid"
	Ask = "ask"
)

// UpdateOptions is used to communicate an update to the [Book].
type UpdateOptions struct {
	Direction BookDirection    `json:"direction,omitempty"`
	ID        string           `json:"orderid,omitempty"`
	Price     *decimal.Decimal `json:"price,omitempty"`
	Quantity  *decimal.Decimal `json:"quantity,omitempty"`
	Timestamp time.Time        `json:"timestamp,omitempty"`
	Silent    bool             `json:"silent,omitempty"`
}

// Update routes the [UpdateOptions] to the correct side of the book and enforces checks to preserve book integrity.
func (b *Book) Update(opts *UpdateOptions) {
	switch opts.Direction {
	case Ask:
		b.Asks.Update(opts)
	case Bid:
		b.Bids.Update(opts)
	}
	if b.NoBookCrossing {
		b.EnforceOrder()
	}
	if b.EnableMaxDepth {
		b.EnforceDepth()
	}
	if !opts.Silent {
		b.OnUpdated.Call(opts)
	}
}

// BestBid returns the highest bid price level.
func (b *Book) BestBid() *Level {
	return b.Bids.High
}

// BestAsk returns the lowest ask price level.
func (b *Book) BestAsk() *Level {
	return b.Asks.Low
}

// WorstAsk returns the highest ask price level.
func (b *Book) WorstAsk() *Level {
	return b.Asks.High
}

// WorstBid returns the lowest bid price level.
func (b *Book) WorstBid() *Level {
	return b.Bids.Low
}

// EnforceOrder check whether the bid is greater or equal to the ask and attempts to remove older conflicting price levels.
func (b *Book) EnforceOrder() {
	for bid, ask := b.BestBid(), b.BestAsk(); bid != nil && ask != nil && bid.Price.Cmp(ask.Price) >= 0; bid, ask = b.BestBid(), b.BestAsk() {
		b.OnBookCrossed.Call(&CrossedResult{
			Bid: bid,
			Ask: ask,
		})
		var input *UpdateOptions
		if bid.Timestamp.After(ask.Timestamp) {
			input = &UpdateOptions{
				Direction: Ask,
				Price:     ask.Price,
				Quantity:  decimal.NewFromInt64(0),
				Timestamp: time.Now(),
			}
		} else {
			input = &UpdateOptions{
				Direction: Bid,
				Price:     bid.Price,
				Quantity:  decimal.NewFromInt64(0),
				Timestamp: time.Now(),
			}
		}
		b.Update(input)
	}
}

type CrossedResult struct {
	Bid *Level `json:"bid,omitempty"`
	Ask *Level `json:"ask,omitempty"`
}
type MaxDepthExceededResult struct {
	Side         BookDirection `json:"side,omitempty"`
	CurrentDepth int           `json:"currentDepth,omitempty"`
	MaxDepth     int           `json:"maxDepth,omitempty"`
	Worst        *Level        `json:"worst,omitempty"`
}

func (b *Book) EnforceDepth() {
	for len(b.Bids.Levels) > b.MaxDepth {
		b.OnMaxDepthExceeded.Call(&MaxDepthExceededResult{
			Side:         Bid,
			CurrentDepth: len(b.Bids.Levels),
			MaxDepth:     b.MaxDepth,
			Worst:        b.WorstBid(),
		})
		b.Update(&UpdateOptions{
			Direction: Bid,
			Price:     b.WorstBid().Price,
			Quantity:  decimal.NewFromInt64(0),
			Timestamp: time.Now(),
		})
	}
	for len(b.Asks.Levels) > b.MaxDepth {
		b.OnMaxDepthExceeded.Call(&MaxDepthExceededResult{
			Side:         Ask,
			CurrentDepth: len(b.Asks.Levels),
			MaxDepth:     b.MaxDepth,
			Worst:        b.WorstAsk(),
		})
		b.Update(&UpdateOptions{
			Direction: Ask,
			Price:     b.WorstAsk().Price,
			Quantity:  decimal.NewFromInt64(0),
			Timestamp: time.Now(),
		})

	}
}
