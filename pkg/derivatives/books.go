package derivatives

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/krakenfx/api-go/v2/internal/helper"
	"github.com/krakenfx/api-go/v2/pkg/book"
	"github.com/krakenfx/api-go/v2/pkg/callback"
	"github.com/krakenfx/api-go/v2/pkg/decimal"
	"github.com/krakenfx/api-go/v2/pkg/kraken"
)

// BookManager manages the lifecycle of a collection of [Book] structs.
type BookManager struct {
	books        map[string]*book.Book
	mux          sync.RWMutex
	OnCreateBook *callback.Manager[*book.Book]
}

// NewBookManager constructs a new [BookManager] struct.
func NewBookManager() *BookManager {
	return &BookManager{
		books:        make(map[string]*book.Book),
		OnCreateBook: callback.NewManager[*book.Book](),
	}
}

// Update accepts a [kraken.WebSocketMessage] and processes an update.
func (bm *BookManager) Update(m *callback.Event[*kraken.WebSocketMessage]) error {
	event, err := m.Data.Map()
	if err != nil {
		return err
	}
	channel, err := helper.Traverse[string](event, "feed")
	if err != nil || !slices.Contains([]string{"book_snapshot", "book"}, *channel) {
		return nil
	}
	productID, err := helper.Traverse[string](event, "product_id")
	if err != nil {
		return nil
	}
	book := bm.GetBook(*productID)
	if book == nil {
		book = bm.CreateBook(*productID)
	}
	switch *channel {
	case "book_snapshot":
		return bm.UpdateSnapshot(book, event)
	case "book":
		return bm.UpdateDelta(book, event)
	default:
		return fmt.Errorf("unknown channel: %s", *channel)
	}
}

// UpdateSnapshot processes a snapshot response into the book.
func (bm *BookManager) UpdateSnapshot(b *book.Book, m map[string]any) error {
	timestamp, err := helper.Traverse[json.Number](m, "timestamp")
	if err != nil {
		return err
	}
	timestampInt, err := timestamp.Int64()
	if err != nil {
		return fmt.Errorf("timestamp: %ws", err)
	}
	timestampTime := time.UnixMilli(timestampInt)
	bids, err := helper.Traverse[[]any](m, "bids")
	if err != nil {
		return err
	}
	asks, err := helper.Traverse[[]any](m, "asks")
	if err != nil {
		return err
	}
	sides := map[book.BookDirection][]any{
		book.Bid: *bids,
		book.Ask: *asks,
	}
	for direction, records := range sides {
		for _, record := range records {
			price, err := helper.Traverse[json.Number](record, "price")
			if err != nil {
				return err
			}
			priceMoney, err := decimal.NewFromString(price.String())
			if err != nil {
				return fmt.Errorf("price: %w", err)
			}
			quantity, err := helper.Traverse[json.Number](record, "qty")
			if err != nil {
				return err
			}
			priceQuantity, err := decimal.NewFromString(quantity.String())
			if err != nil {
				return fmt.Errorf("quantity: %w", err)
			}
			b.Update(&book.UpdateOptions{
				Direction: direction,
				Price:     priceMoney,
				Quantity:  priceQuantity,
				Timestamp: timestampTime,
			})
		}
	}
	return nil
}

// UpdateDelta processes a delta response into the book.
func (bm *BookManager) UpdateDelta(b *book.Book, m map[string]any) error {
	timestamp, err := helper.Traverse[json.Number](m, "timestamp")
	if err != nil {
		return err
	}
	timestampInt, err := timestamp.Int64()
	if err != nil {
		return fmt.Errorf("timestamp: %ws", err)
	}
	timestampTime := time.UnixMilli(timestampInt)
	price, err := helper.Traverse[json.Number](m, "price")
	if err != nil {
		return err
	}
	priceMoney, err := decimal.NewFromString(price.String())
	if err != nil {
		return fmt.Errorf("price: %w", err)
	}
	quantity, err := helper.Traverse[json.Number](m, "qty")
	if err != nil {
		return err
	}
	priceQuantity, err := decimal.NewFromString(quantity.String())
	if err != nil {
		return fmt.Errorf("quantity: %w", err)
	}
	side, err := helper.Traverse[string](m, "side")
	if err != nil {
		return err
	}
	var direction book.BookDirection
	switch *side {
	case "sell":
		direction = book.Ask
	case "buy":
		direction = book.Bid
	default:
		return fmt.Errorf("unknown direction: %s", *side)
	}
	b.Update(&book.UpdateOptions{
		Direction: direction,
		Price:     priceMoney,
		Quantity:  priceQuantity,
		Timestamp: timestampTime,
	})
	return nil
}

// CreateBook constructs a managed [Book] struct.
func (b *BookManager) CreateBook(name string) *book.Book {
	b.mux.Lock()
	defer b.mux.Unlock()
	nameUpper := strings.ToUpper(name)
	book := book.New()
	book.Name = nameUpper
	book.EnableMaxDepth = false
	b.books[nameUpper] = book
	b.OnCreateBook.Call(book)
	return book
}

// GetBook returns the [Book] struct associated with the given symbol.
func (bm *BookManager) GetBook(name string) *book.Book {
	bm.mux.RLock()
	defer bm.mux.RUnlock()
	nameUpper := strings.ToUpper(name)
	book, ok := bm.books[nameUpper]
	if !ok {
		return nil
	}
	return book
}

// GetBooks returns a list of all managed [Book] structs.
func (bm *BookManager) GetBooks() []string {
	bm.mux.RLock()
	defer bm.mux.RUnlock()
	var books []string
	for k := range bm.books {
		books = append(books, k)
	}
	return books
}
