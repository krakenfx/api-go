package derivatives

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/krakenfx/api-go/pkg/kraken"
)

// BookManager manages the lifecycle of a collection of [Book] structs.
type BookManager struct {
	books        map[string]*Book
	mux          sync.RWMutex
	OnCreateBook *kraken.CallbackManager[*Book]
}

// NewBookManager constructs a new [BookManager] struct.
func NewBookManager() *BookManager {
	return &BookManager{
		books:        make(map[string]*Book),
		OnCreateBook: kraken.NewCallbackManager[*Book](),
	}
}

// Update accepts a [kraken.WebSocketMessage] and processes an update.
func (b *BookManager) Update(m *kraken.Event[*kraken.WebSocketMessage]) error {
	event, err := m.Data.Map()
	if err != nil {
		return err
	}
	channel, err := kraken.Traverse[string](event, "feed")
	if err != nil || !slices.Contains([]string{"book_snapshot", "book"}, *channel) {
		return nil
	}
	productID, err := kraken.Traverse[string](event, "product_id")
	if err != nil {
		return nil
	}
	book := b.GetBook(*productID)
	if book == nil {
		book = b.CreateBook(*productID)
	}
	switch *channel {
	case "book_snapshot":
		return b.UpdateSnapshot(book, event)
	case "book":
		return b.UpdateDelta(book, event)
	default:
		return fmt.Errorf("unknown channel: %s", *channel)
	}
}

// UpdateSnapshot processes a snapshot response into the book.
func (b *BookManager) UpdateSnapshot(book *Book, m map[string]any) error {
	timestamp, err := kraken.Traverse[json.Number](m, "timestamp")
	if err != nil {
		return err
	}
	timestampInt, err := timestamp.Int64()
	if err != nil {
		return fmt.Errorf("timestamp: %ws", err)
	}
	timestampTime := time.UnixMilli(timestampInt)
	bids, err := kraken.Traverse[[]any](m, "bids")
	if err != nil {
		return err
	}
	asks, err := kraken.Traverse[[]any](m, "asks")
	if err != nil {
		return err
	}
	sides := map[kraken.BookDirection][]any{
		kraken.Bid: *bids,
		kraken.Ask: *asks,
	}
	for direction, records := range sides {
		for _, record := range records {
			price, err := kraken.Traverse[json.Number](record, "price")
			if err != nil {
				return err
			}
			priceMoney, err := kraken.NewMoneyFromString(price.String())
			if err != nil {
				return fmt.Errorf("price: %w", err)
			}
			quantity, err := kraken.Traverse[json.Number](record, "qty")
			if err != nil {
				return err
			}
			priceQuantity, err := kraken.NewMoneyFromString(quantity.String())
			if err != nil {
				return fmt.Errorf("quantity: %w", err)
			}
			book.Update(&kraken.BookUpdateOptions{
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
func (b *BookManager) UpdateDelta(book *Book, m map[string]any) error {
	timestamp, err := kraken.Traverse[json.Number](m, "timestamp")
	if err != nil {
		return err
	}
	timestampInt, err := timestamp.Int64()
	if err != nil {
		return fmt.Errorf("timestamp: %ws", err)
	}
	timestampTime := time.UnixMilli(timestampInt)
	price, err := kraken.Traverse[json.Number](m, "price")
	if err != nil {
		return err
	}
	priceMoney, err := kraken.NewMoneyFromString(price.String())
	if err != nil {
		return fmt.Errorf("price: %w", err)
	}
	quantity, err := kraken.Traverse[json.Number](m, "qty")
	if err != nil {
		return err
	}
	priceQuantity, err := kraken.NewMoneyFromString(quantity.String())
	if err != nil {
		return fmt.Errorf("quantity: %w", err)
	}
	side, err := kraken.Traverse[string](m, "side")
	if err != nil {
		return err
	}
	var direction kraken.BookDirection
	switch *side {
	case "sell":
		direction = kraken.Ask
	case "buy":
		direction = kraken.Bid
	default:
		return fmt.Errorf("unknown direction: %s", *side)
	}
	book.Update(&kraken.BookUpdateOptions{
		Direction: direction,
		Price:     priceMoney,
		Quantity:  priceQuantity,
		Timestamp: timestampTime,
	})
	return nil
}

// CreateBook constructs a managed [Book] struct.
func (b *BookManager) CreateBook(name string) *Book {
	b.mux.Lock()
	defer b.mux.Unlock()
	nameUpper := strings.ToUpper(name)
	book := NewBook()
	book.Symbol = nameUpper
	book.EnableMaxDepth = false
	b.books[nameUpper] = book
	b.OnCreateBook.Call(book)
	return book
}

// GetBook returns the [Book] struct associated with the given symbol.
func (b *BookManager) GetBook(name string) *Book {
	b.mux.RLock()
	defer b.mux.RUnlock()
	nameUpper := strings.ToUpper(name)
	book, ok := b.books[nameUpper]
	if !ok {
		return nil
	}
	return book
}

// GetBooks returns a list of all managed [Book] structs.
func (b *BookManager) GetBooks() []string {
	b.mux.RLock()
	defer b.mux.RUnlock()
	var books []string
	for k := range b.books {
		books = append(books, k)
	}
	return books
}

// Book wraps a [kraken.Book] struct with a symbol field.
type Book struct {
	Symbol string
	*kraken.Book
}

// NewBook constructs a new [Book] struct with default values.
func NewBook() *Book {
	return &Book{
		Book: kraken.NewBook(),
	}
}
