package main

import (
	"fmt"
	"os"
	"time"

	"github.com/krakenfx/api-go/internal/helper"
	"github.com/krakenfx/api-go/pkg/book"
	"github.com/krakenfx/api-go/pkg/callback"
	"github.com/krakenfx/api-go/pkg/kraken"
	"github.com/krakenfx/api-go/pkg/spot"
)

func main() {
	client := spot.NewWebSocket()
	client.URL = os.Getenv("KRAKEN_API_SPOT_WS_AUTH_URL")
	client.REST.BaseURL = os.Getenv("KRAKEN_API_SPOT_REST_URL")
	client.REST.PublicKey = os.Getenv("KRAKEN_API_SPOT_PUBLIC")
	client.REST.PrivateKey = os.Getenv("KRAKEN_API_SPOT_SECRET")
	bookManager := spot.NewBookManager()
	bookManager.OnCreateBook.Recurring(func(e *callback.Event[*book.Book]) {
		b := e.Data
		fmt.Printf("Create book: %s\n", b.Name)
		b.OnUpdated.Recurring(func(e *callback.Event[*book.UpdateOptions]) {
			fmt.Printf("%s: %s\n", b.Name, helper.ToJSON(e.Data))
		})
		b.OnBookCrossed.Recurring(func(e *callback.Event[*book.CrossedResult]) {
			fmt.Printf("%s: %s\\n", b.Name, helper.ToJSON(e.Data))
		})
		b.OnMaxDepthExceeded.Recurring(func(e *callback.Event[*book.MaxDepthExceededResult]) {
			fmt.Printf("%s: %s\n", b.Name, helper.ToJSON(e.Data))
		})
		b.OnChecksummed.Recurring(func(e *callback.Event[*book.ChecksumResult]) {
			if !e.Data.Match {
				fmt.Printf("%s: %s\n", b.Name, helper.ToJSON(e.Data))
			}
		})
	})
	client.OnSent.Recurring(func(e *callback.Event[*kraken.WebSocketMessage]) {
		fmt.Printf("Sent: %s\n", e.Data)
		if err := bookManager.Update(e); err != nil {
			panic(err)
		}
	})
	client.OnReceived.Recurring(func(e *callback.Event[*kraken.WebSocketMessage]) {
		if err := bookManager.Update(e); err != nil {
			panic(err)
		}
	})
	client.OnAuthenticated.Recurring(func(e *callback.Event[string]) {
		if err := client.SubL3([]string{"BTC/USD"}, 10); err != nil {
			panic(err)
		}
	})
	client.OnConnected.Recurring(func(e *callback.Event[any]) {
		if err := client.Authenticate(); err != nil {
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
