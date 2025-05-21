package spot

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
	method, err := kraken.Traverse[string](event, "method")
	if err == nil && *method == "subscribe" {
		channel, err := kraken.Traverse[string](event, "params", "channel")
		if err != nil {
			return nil
		}
		if !slices.Contains([]string{"level3", "book"}, *channel) {
			return nil
		}
		symbols, err := kraken.Traverse[[]any](event, "params", "symbol")
		if err != nil {
			return err
		}
		var depthInt int64
		depth, err := kraken.Traverse[json.Number](event, "params", "depth")
		if err != nil {
			depthInt = 10
		} else {
			depthInt, err = depth.Int64()
			if err != nil {
				return fmt.Errorf("depth: %w", err)
			}
		}
		for _, symbol := range *symbols {
			b.CreateBook(fmt.Sprint(symbol), int(depthInt))
		}
		return nil
	}
	channel, err := kraken.Traverse[string](event, "channel")
	if err != nil || !slices.Contains([]string{"level3", "book"}, *channel) {
		return nil
	}
	updates, err := kraken.Traverse[[]any](event, "data")
	if err != nil {
		return err
	}
	for _, update := range *updates {
		bookUpdate, ok := update.(map[string]any)
		if !ok {
			return fmt.Errorf("assert \"%v\" as map[string]any failed", bookUpdate)
		}
		symbol, err := kraken.Traverse[string](bookUpdate, "symbol")
		if err != nil {
			return err
		}
		book := b.GetBook(*symbol)
		if book == nil {
			return fmt.Errorf("%s not found in library (%s)", *symbol, strings.Join(b.GetBooks(), ","))
		}
		switch *channel {
		case "level3":
			if err := b.UpdateL3(book, bookUpdate); err != nil {
				return fmt.Errorf("\"%s\" update l3: %w", *symbol, err)
			}
		case "book":
			if err := b.UpdateL2(book, bookUpdate); err != nil {
				return fmt.Errorf("\"%s\" update l2: %w", *symbol, err)
			}
		}
	}
	return nil
}

// CreateBook constructs a managed [Book] struct.
func (b *BookManager) CreateBook(name string, depth int) *Book {
	b.mux.Lock()
	defer b.mux.Unlock()
	nameUpper := strings.ToUpper(name)
	book := NewBook()
	book.Symbol = nameUpper
	book.EnableMaxDepth = true
	book.MaxDepth = depth
	b.books[nameUpper] = book
	b.OnCreateBook.Call(book)
	return book
}

// GetBook returns the [Book] struct associated with the given symbol.
func (b *BookManager) GetBook(symbol string) *Book {
	b.mux.RLock()
	defer b.mux.RUnlock()
	nameUpper := strings.ToUpper(symbol)
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

// UpdateL2 processes a map into an L2 book and performs a checksum.
func (b *BookManager) UpdateL2(book *Book, m map[string]any) error {
	bids, err := kraken.Traverse[[]any](m, "bids")
	if err != nil {
		return err
	}
	asks, err := kraken.Traverse[[]any](m, "asks")
	if err != nil {
		return err
	}
	var timestamp time.Time
	timestampString, err := kraken.Traverse[string](m, "timestamp")
	if err != nil {
		timestamp = time.Now()
	} else {
		timestamp, err = time.Parse(time.RFC3339, *timestampString)
		if err != nil {
			return fmt.Errorf("timestamp parse: %w", err)
		}
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
				Timestamp: timestamp,
			})
		}
	}
	serverChecksum, err := kraken.Traverse[json.Number](m, "checksum")
	if err != nil {
		return err
	}
	if result := book.L2Checksum(serverChecksum.String()); !result.Match {
		return fmt.Errorf("checksum failed, server \"%s\" versus local \"%s\"", result.ServerChecksum, result.LocalChecksum)
	}
	return nil
}

// UpdateL3 processes a map into an L3 book and performs a checksum.
func (b *BookManager) UpdateL3(book *Book, m map[string]any) error {
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
			id, err := kraken.Traverse[string](record, "order_id")
			if err != nil {
				return err
			}
			price, err := kraken.Traverse[json.Number](record, "limit_price")
			if err != nil {
				return err
			}
			priceMoney, err := kraken.NewMoneyFromString(price.String())
			if err != nil {
				return fmt.Errorf("price: %w", err)
			}
			timestampString, err := kraken.Traverse[string](record, "timestamp")
			if err != nil {
				return err
			}
			timestamp, err := time.Parse(time.RFC3339, *timestampString)
			if err != nil {
				return fmt.Errorf("time parse: %w", err)
			}
			event, _ := kraken.Traverse[string](record, "event")
			quantityMoney := kraken.NewMoneyFromInt64(0)
			if event == nil || *event != "delete" {
				quantity, err := kraken.Traverse[json.Number](record, "order_qty")
				if err != nil {
					return err
				}
				quantityMoney, err = kraken.NewMoneyFromString(quantity.String())
				if err != nil {
					return fmt.Errorf("quantity: %w", err)
				}
			}
			book.Update(&kraken.BookUpdateOptions{
				Direction: direction,
				ID:        *id,
				Price:     priceMoney,
				Quantity:  quantityMoney,
				Timestamp: timestamp,
			})
		}
	}
	serverChecksum, err := kraken.Traverse[json.Number](m, "checksum")
	if err != nil {
		return err
	}
	if result := book.L3Checksum(serverChecksum.String()); !result.Match {
		return fmt.Errorf("checksum failed, server \"%s\" versus local \"%s\"", result.ServerChecksum, result.LocalChecksum)
	}
	return nil
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
