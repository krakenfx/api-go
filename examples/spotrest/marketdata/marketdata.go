package main

import (
	"fmt"
	"os"
	"time"

	"github.com/krakenfx/api-go/pkg/kraken"
	"github.com/krakenfx/api-go/pkg/spot"
)

func main() {
	client := spot.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_SPOT_REST_URL")
	assetPair, err := client.AssetPairs(&spot.AssetPairsRequest{
		Pair: "BTC/USD",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Asset pair: %s\n", kraken.ToJSON(assetPair))
	ticker, err := client.TickerSingle("BTC/USD")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Ticker: %s\n", kraken.ToJSON(ticker))
	depth, err := client.OrderBook(&spot.OrderBookRequest{
		Pair: "BTC/USD",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Depth: %s\n", kraken.ToJSONIndent(depth))
	ohlc, err := client.OHLC(&spot.OHLCRequest{
		Pair:  "BTC/USD",
		Since: int(time.Now().Add(-5 * time.Minute).Unix()),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("OHLC: %s\n", kraken.ToJSON(ohlc))
}
