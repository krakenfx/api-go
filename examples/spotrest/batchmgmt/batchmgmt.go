package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/pkg/kraken"
	"github.com/krakenfx/api-go/pkg/spot"
)

// Spot market to trade from.
var symbol = "BTC/USD"

// Direction of the order.
var direction = "buy"

// Percentage away from market price.
var priceOffsets = []float64{-0.5, -0.25}

// Notional size of each order.
var notionalSize = kraken.NewMoneyFromFloat64(5.0)

func main() {
	client := spot.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_SPOT_REST_URL")
	client.PublicKey = os.Getenv("KRAKEN_API_SPOT_PUBLIC")
	client.PrivateKey = os.Getenv("KRAKEN_API_SPOT_SECRET")

	assets := spot.NewAssetManager()
	if err := assets.Use(client); err != nil {
		panic(err)
	}

	fmt.Printf("> Fetch %s market price\n", symbol)
	ticker, err := client.TickerSingle(symbol)
	if err != nil {
		panic(err)
	}
	var marketPrice *kraken.Money
	switch direction {
	case "buy":
		marketPrice = ticker.Ask[0]
	case "sell":
		marketPrice = ticker.Bid[0]
	default:
		panic("unknown direction")
	}

	fmt.Printf("Market price: %s\n", marketPrice)
	var orderRequests []*spot.OrderRequest
	for _, priceOffset := range priceOffsets {
		limitPrice := marketPrice.OffsetPercent(kraken.NewMoneyFromFloat64(priceOffset))
		limitPrice, err := assets.FormatPrice(symbol, limitPrice)
		if err != nil {
			panic(err)
		}
		volume, err := assets.FormatSize(symbol, notionalSize)
		if err != nil {
			panic(err)
		}
		volume = volume.Div(limitPrice)
		order := &spot.OrderRequest{
			OrderType: "limit",
			Price:     limitPrice.String(),
			Type:      direction,
			Volume:    volume.String(),
		}
		orderRequests = append(orderRequests, order)
		fmt.Printf("%s\n", kraken.ToJSON(order))
	}
	fmt.Printf("> Sending batch order.\n")
	resp, err := client.AddBatch(&spot.AddBatchRequest{
		Orders: orderRequests,
		Pair:   "BTC/USD",
	})
	if err != nil {
		panic(err)
	}
	for _, order := range resp.Result.Orders {
		fmt.Printf("%s - %s\n", order.ID, order.Descr.Order)
	}
	for _, order := range resp.Result.Orders {
		fmt.Printf("> Cancel %s\n", order.ID)
		if _, err := client.CancelOrder(&spot.CancelOrderRequest{
			TxID: order.ID,
		}); err != nil {
			panic(err)
		}
	}
}
