package kraken

import (
	"fmt"
	"hash/crc32"
	"math"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Order book structure for L2 and L3.
type Book struct {
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

	OnUpdated          *CallbackManager[*BookUpdateOptions]
	OnBookCrossed      *CallbackManager[*BookCrossedResult]
	OnMaxDepthExceeded *CallbackManager[*MaxDepthExceededResult]
	OnChecksummed      *CallbackManager[*ChecksumResult]
	mux                sync.RWMutex
}

// NewBook constructs a new [Book] struct with default values.
func NewBook() *Book {
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
		OnUpdated:          NewCallbackManager[*BookUpdateOptions](),
		OnBookCrossed:      NewCallbackManager[*BookCrossedResult](),
		OnMaxDepthExceeded: NewCallbackManager[*MaxDepthExceededResult](),
		OnChecksummed:      NewCallbackManager[*ChecksumResult](),
	}
}

// Midpoint returns the midpoint of the order book.
func (b *Book) Midpoint() *Money {
	bid, ask := b.bestBid(), b.bestAsk()
	switch {
	case bid != nil && ask != nil:
		return bid.Price.
			Add(ask.Price).
			Mul(NewMoneyFromFloat64(0.5))
	case bid != nil:
		return bid.Price
	case ask != nil:
		return ask.Price
	default:
		return NewMoneyFromInt64(0)
	}
}

// Spread returns the relative difference between the bid-ask price.
func (b *Book) Spread() *Money {
	bid, ask := b.bestBid(), b.BestAsk()
	switch {
	case bid == nil || ask == nil:
		return NewMoneyFromInt64(0)
	default:
		return ask.Price.
			SetDecimals(int64(math.Max(float64(ask.Price.Decimals), float64(DefaultDecimals)))).
			Sub(bid.Price).
			Div(ask.Price).
			Mul(NewMoneyFromInt64(100))
	}
}

// Whether a bid or an ask.
type BookDirection string

const (
	Bid = "bid"
	Ask = "ask"
)

// BookUpdateOptions is used to communicate an update to the [Book].
type BookUpdateOptions struct {
	Direction BookDirection `json:"direction,omitempty"`
	ID        string        `json:"orderid,omitempty"`
	Price     *Money        `json:"price,omitempty"`
	Quantity  *Money        `json:"quantity,omitempty"`
	Timestamp time.Time     `json:"timestamp,omitempty"`
}

// Update routes the [BookUpdateOptions] to the correct side of the book and enforces checks to preserve book integrity.
func (b *Book) Update(input *BookUpdateOptions) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.update(input)
}

func (b *Book) update(opts *BookUpdateOptions) {
	switch opts.Direction {
	case Ask:
		b.Asks.Update(opts)
	case Bid:
		b.Bids.Update(opts)
	}
	if b.NoBookCrossing {
		b.enforceOrder()
	}
	if b.EnableMaxDepth {
		b.enforceDepth()
	}
	b.OnUpdated.Call(opts)
}

// BestBid returns the highest bid price level.
func (b *Book) BestBid() *Level {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.bestBid()
}

func (b *Book) bestBid() *Level {
	return b.Bids.High
}

// BestAsk returns the lowest ask price level.
func (b *Book) BestAsk() *Level {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.bestAsk()
}

func (b *Book) bestAsk() *Level {
	return b.Asks.Low
}

// WorstAsk returns the highest ask price level.
func (b *Book) WorstAsk() *Level {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.worstAsk()
}

func (b *Book) worstAsk() *Level {
	return b.Asks.High
}

// WorstBid returns the lowest bid price level.
func (b *Book) WorstBid() *Level {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.worstBid()
}

func (b *Book) worstBid() *Level {
	return b.Bids.Low
}

// EnforceOrder check whether the bid is greater or equal to the ask and attempts to remove older conflicting price levels.
func (b *Book) EnforceOrder() {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.enforceOrder()
}

type BookCrossedResult struct {
	Bid *Level `json:"bid,omitempty"`
	Ask *Level `json:"ask,omitempty"`
}

func (b *Book) enforceOrder() {
	for bid, ask := b.bestBid(), b.bestAsk(); bid != nil && ask != nil && bid.Price.Cmp(ask.Price) >= 0; bid, ask = b.bestBid(), b.bestAsk() {
		b.OnBookCrossed.Call(&BookCrossedResult{
			Bid: bid,
			Ask: ask,
		})
		var input *BookUpdateOptions
		if bid.Timestamp.After(ask.Timestamp) {
			input = &BookUpdateOptions{
				Direction: Ask,
				Price:     ask.Price,
				Quantity:  NewMoneyFromInt64(0),
				Timestamp: time.Now(),
			}
		} else {
			input = &BookUpdateOptions{
				Direction: Bid,
				Price:     bid.Price,
				Quantity:  NewMoneyFromInt64(0),
				Timestamp: time.Now(),
			}
		}
		b.update(input)
	}
}

// EnforceDepth checks whether the length of each side of the book exceeds the max depth and attempts to removes the worst out-of-depth price levels.
func (b *Book) EnforceDepth() {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.enforceDepth()
}

type MaxDepthExceededResult struct {
	Side         BookDirection `json:"side,omitempty"`
	CurrentDepth int           `json:"currentDepth,omitempty"`
	MaxDepth     int           `json:"maxDepth,omitempty"`
	Worst        *Level        `json:"worst,omitempty"`
}

func (b *Book) enforceDepth() {
	for b.Bids.Levels.Length() > b.MaxDepth {
		b.OnMaxDepthExceeded.Call(&MaxDepthExceededResult{
			Side:         Bid,
			CurrentDepth: b.Bids.Levels.Length(),
			MaxDepth:     b.MaxDepth,
			Worst:        b.worstBid(),
		})
		b.update(&BookUpdateOptions{
			Direction: Bid,
			Price:     b.worstBid().Price,
			Quantity:  NewMoneyFromInt64(0),
			Timestamp: time.Now(),
		})
	}
	for b.Asks.Levels.Length() > b.MaxDepth {
		b.OnMaxDepthExceeded.Call(&MaxDepthExceededResult{
			Side:         Ask,
			CurrentDepth: b.Asks.Levels.Length(),
			MaxDepth:     b.MaxDepth,
			Worst:        b.worstAsk(),
		})
		b.update(&BookUpdateOptions{
			Direction: Ask,
			Price:     b.worstAsk().Price,
			Quantity:  NewMoneyFromInt64(0),
			Timestamp: time.Now(),
		})

	}
}

// ChecksumPart contains information regarding a price level or order.
type ChecksumPart struct {
	Level        *Level `json:"-"`
	Order        *Order `json:"order,omitempty"`
	Price        string `json:"price,omitempty"`
	Quantity     string `json:"quantity,omitempty"`
	Concatenated string `json:"concatenated,omitempty"`
}

// ChecksumResult contains the result of the book checksum validation.
type ChecksumResult struct {
	Level          int             `json:"level,omitempty"`
	ServerChecksum string          `json:"serverChecksum,omitempty"`
	LocalChecksum  string          `json:"localChecksum,omitempty"`
	Match          bool            `json:"match,omitempty"`
	AskParts       []*ChecksumPart `json:"askParts,omitempty"`
	BidParts       []*ChecksumPart `json:"bidParts,omitempty"`
	Asks           string          `json:"asks,omitempty"`
	Bids           string          `json:"bids,omitempty"`
}

// L2Checksum verifies that the L2 book is synchronized with the exchange.
//
// https://docs.kraken.com/api/docs/guides/spot-ws-book-v2
func (b *Book) L2Checksum(checksum string) *ChecksumResult {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.l2Checksum(checksum)
}

func (b *Book) l2Checksum(checksum string) *ChecksumResult {
	result := &ChecksumResult{
		Level:          2,
		ServerChecksum: checksum,
	}
	cursor := b.bestAsk()
	for range 10 {
		if cursor == nil {
			break
		}
		price := strings.TrimLeft(strings.ReplaceAll(cursor.GetPriceString(), ".", ""), "0")
		quantity := strings.TrimLeft(strings.ReplaceAll(cursor.GetQuantityString(), ".", ""), "0")
		concatenated := price + quantity
		result.AskParts = append(result.AskParts, &ChecksumPart{
			Level:        cursor,
			Price:        price,
			Quantity:     quantity,
			Concatenated: concatenated,
		})
		result.Asks += concatenated
		cursor = cursor.Higher
	}
	cursor = b.bestBid()
	for range 10 {
		if cursor == nil {
			break
		}
		price := strings.TrimLeft(strings.ReplaceAll(cursor.GetPriceString(), ".", ""), "0")
		quantity := strings.TrimLeft(strings.ReplaceAll(cursor.GetQuantityString(), ".", ""), "0")
		concatenated := price + quantity
		result.BidParts = append(result.BidParts, &ChecksumPart{
			Level:        cursor,
			Price:        price,
			Quantity:     quantity,
			Concatenated: concatenated,
		})
		result.Bids += concatenated
		cursor = cursor.Lower
	}
	result.LocalChecksum = fmt.Sprint(crc32.Checksum([]byte(result.Asks+result.Bids), crc32.IEEETable))
	if result.LocalChecksum == result.ServerChecksum {
		result.Match = true
	}
	b.OnChecksummed.Call(result)
	return result
}

// L3Checksum verifies that the L3 book is synchronized with the exchange.
//
// https://docs.kraken.com/api/docs/guides/spot-ws-l3-v2
func (b *Book) L3Checksum(ref string) *ChecksumResult {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.l3Checksum(ref)
}

func (b *Book) l3Checksum(checksum string) *ChecksumResult {
	result := &ChecksumResult{
		Level:          3,
		ServerChecksum: checksum,
	}
	cursor := b.bestAsk()
	for range 10 {
		if cursor == nil {
			break
		}
		for _, order := range cursor.Queue() {
			price := strings.TrimLeft(strings.ReplaceAll(order.LimitPrice.String(), ".", ""), "0")
			quantity := strings.TrimLeft(strings.ReplaceAll(order.Quantity.String(), ".", ""), "0")
			concatenated := price + quantity
			result.AskParts = append(result.AskParts, &ChecksumPart{
				Level:        cursor,
				Order:        order,
				Price:        price,
				Quantity:     quantity,
				Concatenated: concatenated,
			})
			result.Asks += concatenated
		}
		cursor = cursor.Higher
	}
	cursor = b.bestBid()
	for range 10 {
		if cursor == nil {
			break
		}
		for _, order := range cursor.Queue() {
			price := strings.TrimLeft(strings.ReplaceAll(order.LimitPrice.String(), ".", ""), "0")
			quantity := strings.TrimLeft(strings.ReplaceAll(order.Quantity.String(), ".", ""), "0")
			concatenated := price + quantity
			result.BidParts = append(result.BidParts, &ChecksumPart{
				Level:        cursor,
				Order:        order,
				Price:        price,
				Quantity:     quantity,
				Concatenated: concatenated,
			})
			result.Bids += concatenated
		}
		cursor = cursor.Lower
	}
	result.LocalChecksum = fmt.Sprint(crc32.Checksum([]byte(result.Asks+result.Bids), crc32.IEEETable))
	if result.LocalChecksum == result.ServerChecksum {
		result.Match = true
	}
	b.OnChecksummed.Call(result)
	return result
}

// Side encompasses the price levels in one side of the book.
type Side struct {
	Direction BookDirection
	High      *Level
	Low       *Level
	Last      *Level
	Levels    *Map[string, Level]
	mux       sync.RWMutex
}

// NewSide constructs a new [Side] with default values.
func NewSide() *Side {
	return &Side{
		Levels: NewMap[string, Level](),
	}
}

// FindAdjacent finds the nearest price level close to the given price.
func (s *Side) FindAdjacent(price *Money) *Level {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.findAdjacent(price)
}

func (s *Side) findAdjacent(price *Money) *Level {
	if s.High == nil || s.Low == nil {
		return nil
	}
	if price.Cmp(s.High.Price) > 0 {
		return s.High
	}
	if price.Cmp(s.Low.Price) < 0 {
		return s.Low
	}
	highDistance := s.High.Price.Sub(price)
	lowDistance := price.Sub(s.Low.Price)
	if highDistance.Cmp(lowDistance) == 1 {
		return s.findAdjacentBelow(price)
	} else {
		return s.findAdjacentAbove(price)
	}
}

// FindAdjacentBelow finds the nearest price level from below the given price.
func (s *Side) FindAdjacentBelow(price *Money) *Level {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.findAdjacentBelow(price)
}

func (s *Side) findAdjacentBelow(price *Money) *Level {
	if s.Low == nil || price.Cmp(s.Low.Price) <= 0 {
		return nil
	}
	nearest := s.Low
	for {
		next := nearest.Higher
		if next == nil {
			break
		}
		nearestDiff := nearest.Price.Sub(price).Abs()
		nextDiff := next.Price.Sub(price).Abs()
		if nearestDiff.Cmp(nextDiff) < 0 || next.Price.Cmp(price) >= 0 {
			break
		}
		nearest = next
	}
	return nearest
}

// FindAdjacentAbove finds the nearest price level above the given price.
func (s *Side) FindAdjacentAbove(price *Money) *Level {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.findAdjacentAbove(price)
}

func (s *Side) findAdjacentAbove(price *Money) *Level {
	if s.High == nil || price.Cmp(s.High.Price) >= 0 {
		return nil
	}
	nearest := s.High
	for {
		next := nearest.Lower
		if next == nil {
			break
		}
		nearestDiff := nearest.Price.Sub(price).Abs()
		nextDiff := next.Price.Sub(price).Abs()
		if nearestDiff.Cmp(nextDiff) < 0 || next.Price.Cmp(price) <= 0 {
			break
		}
		nearest = next
	}
	return nearest
}

// Update interprets a [BookUpdateOptions] message and decides if it should add, update, or delete the price level.
func (s *Side) Update(opts *BookUpdateOptions) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.update(opts)
}

func (s *Side) update(opts *BookUpdateOptions) {
	level, err := s.Levels.Get(opts.Price.String())
	if err != nil && opts.Quantity.Sign() == 1 {
		s.add(opts)
	} else if err == nil {
		level.update(opts)
	}
	if level != nil && level.Quantity.Sign() <= 0 {
		s.delete(level)
	}
}

func (s *Side) add(opts *BookUpdateOptions) {
	level := NewLevel()
	level.Price = opts.Price
	nearest := s.findAdjacent(level.Price)
	if nearest == nil || s.High == nil || level.Price.Cmp(s.High.Price) > 0 {
		s.High = level
	}
	if nearest == nil || s.Low == nil || level.Price.Cmp(s.Low.Price) < 0 {
		s.Low = level
	}
	if nearest != nil {
		if level.Price.Cmp(nearest.Price) > 0 {
			level.Lower = nearest
			level.Higher = nearest.Higher
			nearest.Higher = level
			if level.Higher != nil {
				level.Higher.Lower = level
			}
		} else if level.Price.Cmp(nearest.Price) < 0 {
			level.Higher = nearest
			level.Lower = nearest.Lower
			nearest.Lower = level
			if level.Lower != nil {
				level.Lower.Higher = level
			}
		}
	}
	level.update(opts)
	s.Levels.Set(level.GetPriceString(), level)
}

func (s *Side) delete(level *Level) {
	if s.High == level {
		s.High = level.Lower
	}
	if s.Low == level {
		s.Low = level.Higher
	}
	if level.Lower != nil {
		level.Lower.Higher = level.Higher
	}
	if level.Higher != nil {
		level.Higher.Lower = level.Lower
	}
	s.Levels.Delete(level.GetPriceString())
}

// Level contains price level information.
type Level struct {
	Price      *Money    `json:"price,omitempty"`
	Quantity   *Money    `json:"quantity,omitempty"`
	Timestamp  time.Time `json:"time,omitempty"`
	Lower      *Level    `json:"-"`
	Higher     *Level    `json:"-"`
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

func (l *Level) update(opts *BookUpdateOptions) {
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
		totalQuantity := NewMoneyFromInt64(0)
		for _, order := range l.orders {
			totalQuantity = totalQuantity.
				SetDecimals(int64(math.Max(float64(totalQuantity.Decimals), float64(order.Quantity.Decimals)))).
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
	ID         string    `json:"id,omitempty"`
	LimitPrice *Money    `json:"limitPrice,omitempty"`
	Quantity   *Money    `json:"quantity,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
	Level      *Level    `json:"-"`
}
