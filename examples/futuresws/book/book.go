package main

import (
	"fmt"
	"os"
	"time"

	"github.com/krakenfx/api-go/pkg/derivatives"
	"github.com/krakenfx/api-go/pkg/kraken"
)

var contract string = "PF_XBTUSD"

func main() {
	client := derivatives.NewWebSocket()
	client.URL = os.Getenv("KRAKEN_API_FUTURES_WS_URL")
	client.REST.BaseURL = os.Getenv("KRAKEN_API_FUTURES_REST_URL")
	bookManager := derivatives.NewBookManager()
	bookManager.OnCreateBook.Recurring(func(e *kraken.Event[*derivatives.Book]) {
		book := e.Data
		fmt.Printf("Create book: %s\n", book.Symbol)
		book.OnUpdated.Recurring(func(e *kraken.Event[*kraken.BookUpdateOptions]) {
			fmt.Printf("%s: %s\n", book.Symbol, kraken.ToJSON(e.Data))
		})
		book.OnBookCrossed.Recurring(func(e *kraken.Event[*kraken.BookCrossedResult]) {
			fmt.Printf("%s: %s\\n", book.Symbol, kraken.ToJSON(e.Data))
		})
	})
	client.OnSent.Recurring(func(e *kraken.Event[*kraken.WebSocketMessage]) {
		fmt.Printf("Sent: %s\n", e.Data)
	})
	client.OnReceived.Recurring(func(e *kraken.Event[*kraken.WebSocketMessage]) {
		if err := bookManager.Update(e); err != nil {
			panic(err)
		}
	})
	client.OnConnected.Recurring(func(e *kraken.Event[any]) {
		if err := client.SubBook(contract); err != nil {
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
