package main

import (
	"fmt"
	"os"
	"time"

	"github.com/krakenfx/api-go/pkg/kraken"
	"github.com/krakenfx/api-go/pkg/spot"
)

func main() {
	client := spot.NewWebSocket()
	client.URL = os.Getenv("KRAKEN_API_SPOT_WS_URL")
	client.REST.BaseURL = os.Getenv("KRAKEN_API_SPOT_REST_URL")
	bookManager := spot.NewBookManager()
	bookManager.OnCreateBook.Recurring(func(e *kraken.Event[*spot.Book]) {
		book := e.Data
		fmt.Printf("Create book: %s\n", book.Symbol)
		book.OnUpdated.Recurring(func(e *kraken.Event[*kraken.BookUpdateOptions]) {
			fmt.Printf("%s: %s\n", book.Symbol, kraken.ToJSON(e.Data))
		})
		book.OnBookCrossed.Recurring(func(e *kraken.Event[*kraken.BookCrossedResult]) {
			fmt.Printf("%s: %s\\n", book.Symbol, kraken.ToJSON(e.Data))
		})
		book.OnMaxDepthExceeded.Recurring(func(e *kraken.Event[*kraken.MaxDepthExceededResult]) {
			fmt.Printf("%s: %s\n", book.Symbol, kraken.ToJSON(e.Data))
		})
		book.OnChecksummed.Recurring(func(e *kraken.Event[*kraken.ChecksumResult]) {
			if !e.Data.Match {
				fmt.Printf("%s: %s\n", book.Symbol, kraken.ToJSON(e.Data))
			}
		})
	})
	client.OnSent.Recurring(func(e *kraken.Event[*kraken.WebSocketMessage]) {
		fmt.Printf("Sent: %s\n", e.Data)
		if err := bookManager.Update(e); err != nil {
			panic(err)
		}
	})
	client.OnReceived.Recurring(func(e *kraken.Event[*kraken.WebSocketMessage]) {
		if err := bookManager.Update(e); err != nil {
			panic(err)
		}
	})
	client.OnConnected.Recurring(func(e *kraken.Event[any]) {
		if err := client.SubBook([]string{"BTC/USD"}, 10); err != nil {
			panic(err)
		}
	})
	if err := client.Connect(); err != nil {
		panic(err)
	}
	for {
		time.Sleep(time.Second)
	}
}
