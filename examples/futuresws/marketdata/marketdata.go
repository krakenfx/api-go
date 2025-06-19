package main

import (
	"fmt"
	"os"
	"time"

	"github.com/krakenfx/api-go/pkg/callback"
	"github.com/krakenfx/api-go/pkg/derivatives"
	"github.com/krakenfx/api-go/pkg/kraken"
)

// Derivative contract
var contract = "PF_XBTUSD"

func main() {
	client := derivatives.NewWebSocket()
	client.URL = os.Getenv("KRAKEN_API_FUTURES_WS_URL")
	client.OnSent.Recurring(func(e *callback.Event[*kraken.WebSocketMessage]) {
		fmt.Printf("Sent: %s\n", e.Data)
	})
	client.OnReceived.Recurring(func(e *callback.Event[*kraken.WebSocketMessage]) {
		fmt.Printf("Received: %s\n", e.Data)
	})
	client.OnConnected.Recurring(func(e *callback.Event[any]) {
		if err := client.SubTicker(contract); err != nil {
			panic(err)
		}
		if err := client.SubBook(contract); err != nil {
			panic(err)
		}
		if err := client.SubTrade(contract); err != nil {
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
