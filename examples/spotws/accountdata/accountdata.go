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
	client.URL = os.Getenv("KRAKEN_API_SPOT_WS_AUTH_URL")
	client.REST.BaseURL = os.Getenv("KRAKEN_API_SPOT_REST_URL")
	client.REST.PublicKey = os.Getenv("KRAKEN_API_SPOT_PUBLIC")
	client.REST.PrivateKey = os.Getenv("KRAKEN_API_SPOT_SECRET")
	client.OnSent.Recurring(func(e *kraken.Event[*kraken.WebSocketMessage]) {
		fmt.Printf("Sent: %s\n", e.Data)
	})
	client.OnReceived.Recurring(func(e *kraken.Event[*kraken.WebSocketMessage]) {
		fmt.Printf("Received: %s\n", e.Data)
	})
	client.OnAuthenticated.Recurring(func(e *kraken.Event[string]) {
		err := client.SubBalances()
		if err != nil {
			panic(err)
		}
		err = client.SubExecutions()
		if err != nil {
			panic(err)
		}
	})
	client.OnConnected.Recurring(func(e *kraken.Event[any]) {
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
