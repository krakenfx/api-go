package main

import (
	"fmt"
	"os"
	"time"

	"github.com/krakenfx/api-go/v2/pkg/callback"
	"github.com/krakenfx/api-go/v2/pkg/kraken"
	"github.com/krakenfx/api-go/v2/pkg/spot"
)

func main() {
	client := spot.NewWebSocket()
	client.URL = os.Getenv("KRAKEN_API_SPOT_WS_URL")
	client.REST.BaseURL = os.Getenv("KRAKEN_API_SPOT_REST_URL")
	client.OnSent.Recurring(func(e *callback.Event[*kraken.WebSocketMessage]) {
		fmt.Printf("Sent: %s\n", e.Data)
	})
	client.OnReceived.Recurring(func(e *callback.Event[*kraken.WebSocketMessage]) {
		fmt.Printf("Received: %s\n", e.Data)
	})
	if err := client.Connect(); err != nil {
		panic(err)
	}
	if err := client.SubTicker([]string{"BTC/USD"}); err != nil {
		panic(err)
	}
	for {
		time.Sleep(time.Second)
	}
}
