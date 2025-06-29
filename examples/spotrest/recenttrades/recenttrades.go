package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/v2/internal/helper"
	"github.com/krakenfx/api-go/v2/pkg/spot"
)

func main() {
	client := spot.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_SPOT_REST_URL")
	resp, err := client.RecentTrades(&spot.RecentTradesRequest{
		Pair:  "BTC/USD",
		Count: 5,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Recent trades: %s", helper.ToJSONIndent(resp))
}
