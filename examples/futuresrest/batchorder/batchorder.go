// Places a maker order and cancels them in the same request with batchorder.
package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/pkg/derivatives"
	"github.com/krakenfx/api-go/pkg/kraken"
)

// Derivative contract.
var contract = "PF_XBTUSD"

// Notional size.
var notionalSize = kraken.NewMoneyFromFloat64(10)

// Side
var side = "buy"

// Offset from market price in percentage.
var priceOffset = kraken.NewMoneyFromFloat64(-0.5)

func main() {
	client := derivatives.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_FUTURES_REST_URL")
	client.PublicKey = os.Getenv("KRAKEN_API_FUTURES_PUBLIC")
	client.PrivateKey = os.Getenv("KRAKEN_API_FUTURES_SECRET")

	asset := derivatives.NewAssetManager()
	if err := asset.Use(client); err != nil {
		panic(err)
	}

	fmt.Printf("> Retrieving %s market price\n", contract)
	ticker, err := client.TickerSymbol(contract)
	if err != nil {
		panic(err)
	}
	price := ticker.Result.Data.Bid.Add(ticker.Result.Data.Ask).Div(kraken.NewMoneyFromInt64(2))
	fmt.Printf("Mid price: %s\n", price)

	limitPrice := price.OffsetPercent(priceOffset)
	limitPrice, err = asset.FormatPrice(contract, limitPrice)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Limit price: %s\n", limitPrice)
	size, err := asset.FormatSize(contract, notionalSize)
	if err != nil {
		panic(err)
	}
	size = size.Div(limitPrice)
	fmt.Printf("Size: %s\n", size)

	clientOrderID := kraken.UUID()
	fmt.Printf("Client order ID: %s\n", clientOrderID)

	fmt.Printf("> Sending batch order request\n")
	response, err := client.BatchOrder(&derivatives.BatchOrderRequest{
		JSON: &derivatives.BatchOrderJson{
			BatchOrder: []*derivatives.BatchOrderInstruction{
				{
					Order:         "send",
					OrderTag:      "test",
					ClientOrderID: clientOrderID,
					Symbol:        "PF_XBTUSD",
					Side:          side,
					Size:          size.String(),
					LimitPrice:    limitPrice.String(),
					OrderType:     "lmt",
				},
				{
					Order:         "cancel",
					ClientOrderID: clientOrderID,
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %s\n", kraken.ToJSONIndent(response))
}
