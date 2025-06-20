package main

import (
	"fmt"
	"os"
	"time"

	"github.com/krakenfx/api-go/pkg/callback"
	"github.com/krakenfx/api-go/pkg/derivatives"
	"github.com/krakenfx/api-go/pkg/kraken"
)

func main() {
	client := derivatives.NewWebSocket()
	client.URL = os.Getenv("KRAKEN_API_FUTURES_WS_URL")
	client.PublicKey = os.Getenv("KRAKEN_API_FUTURES_PUBLIC")
	client.PrivateKey = os.Getenv("KRAKEN_API_FUTURES_SECRET")
	client.OnSent.Recurring(func(e *callback.Event[*kraken.WebSocketMessage]) {
		fmt.Printf("Sent: %s\n", e.Data)
	})
	client.OnReceived.Recurring(func(e *callback.Event[*kraken.WebSocketMessage]) {
		fmt.Printf("Received: %s\n", e.Data)
	})
	client.OnAuthenticated.Recurring(func(e *callback.Event[string]) {
		if err := client.SubOpenOrders(); err != nil {
			panic(err)
		}
		if err := client.SubExecutions(); err != nil {
			panic(err)
		}
	})
	client.OnConnected.Recurring(func(e *callback.Event[any]) {
		go func() {
			if err := client.Authenticate(); err != nil {
				panic(err)
			}
		}()
	})
	if err := client.Connect(); err != nil {
		panic(err)
	}
	for {
		time.Sleep(time.Second)
	}
}
