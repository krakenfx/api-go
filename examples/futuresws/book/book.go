package main

import (
	"fmt"
	"os"
	"time"

	"github.com/krakenfx/api-go/v2/internal/helper"
	"github.com/krakenfx/api-go/v2/pkg/book"
	"github.com/krakenfx/api-go/v2/pkg/callback"
	"github.com/krakenfx/api-go/v2/pkg/derivatives"
	"github.com/krakenfx/api-go/v2/pkg/kraken"
)

var contract string = "PF_XBTUSD"

func main() {
	client := derivatives.NewWebSocket()
	client.URL = os.Getenv("KRAKEN_API_FUTURES_WS_URL")
	client.REST.BaseURL = os.Getenv("KRAKEN_API_FUTURES_REST_URL")
	bookManager := derivatives.NewBookManager()
	bookManager.OnCreateBook.Recurring(func(e *callback.Event[*book.Book]) {
		b := e.Data
		fmt.Printf("Create book: %s\n", b.Name)
		b.OnUpdated.Recurring(func(e *callback.Event[*book.UpdateOptions]) {
			fmt.Printf("%s: %s\n", b.Name, helper.ToJSON(e.Data))
		})
		b.OnBookCrossed.Recurring(func(e *callback.Event[*book.CrossedResult]) {
			fmt.Printf("%s: %s\\n", b.Name, helper.ToJSON(e.Data))
		})
	})
	client.OnSent.Recurring(func(e *callback.Event[*kraken.WebSocketMessage]) {
		fmt.Printf("Sent: %s\n", e.Data)
	})
	client.OnReceived.Recurring(func(e *callback.Event[*kraken.WebSocketMessage]) {
		if err := bookManager.Update(e); err != nil {
			panic(err)
		}
	})
	client.OnConnected.Recurring(func(e *callback.Event[any]) {
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
