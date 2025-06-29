package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/v2/internal/helper"
	"github.com/krakenfx/api-go/v2/pkg/decimal"
	"github.com/krakenfx/api-go/v2/pkg/derivatives"
)

// Derivative contract.
var contract = "PI_XBTUSD"

// Direction of the order.
var side = "buy"

// Size of the order in quote unit.
var notionalSize = decimal.NewFromFloat64(5)

func main() {
	client := derivatives.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_FUTURES_REST_URL")
	client.PublicKey = os.Getenv("KRAKEN_API_FUTURES_PUBLIC")
	client.PrivateKey = os.Getenv("KRAKEN_API_FUTURES_SECRET")
	asset := derivatives.NewNormalizer()
	if err := asset.Use(client); err != nil {
		panic(err)
	}

	fmt.Printf("> Fetch %s market price\n", contract)
	ticker, err := client.TickerSymbol(contract)
	if err != nil {
		panic(err)
	}
	price := ticker.Result.Data.Bid.Add(ticker.Result.Data.Ask).Div(decimal.NewFromInt64(2))
	fmt.Printf("Mid price: %s\n", price)
	limitPrice := price.OffsetPercent(decimal.NewFromFloat64(-0.05))
	limitPrice, err = asset.FormatPrice(contract, limitPrice)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Limit price: %s\n", limitPrice)
	size, err := asset.FormatSize(contract, notionalSize)
	if err != nil {
		panic(err)
	}

	fmt.Printf("> Submit send order request\n")
	cliOrdID := helper.UUID()
	response, err := client.SendOrder(&derivatives.OrderRequest{
		OrderType:     "lmt",
		Symbol:        contract,
		Side:          side,
		LimitPrice:    limitPrice.String(),
		ClientOrderID: cliOrdID,
		Size:          size.String(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %s\n", helper.ToJSON(response))

	fmt.Printf("> Set limit price to -10 ticks away from the original limit price\n")
	limitPrice = limitPrice.OffsetTicks(decimal.NewFromInt64(-10))
	fmt.Printf("Limit price: %s\n", limitPrice)

	fmt.Printf("> Submit edit order request\n")
	editOrderResponse, err := client.EditOrder(&derivatives.OrderRequest{
		ClientOrderID: cliOrdID,
		LimitPrice:    limitPrice.String(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %s\n", helper.ToJSON(editOrderResponse))

	fmt.Printf("> Sending cancel order request\n")
	cancelOrderResponse, err := client.CancelOrder(&derivatives.CancelOrderRequest{
		ClientOrderID: cliOrdID,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %s\n", helper.ToJSON(cancelOrderResponse))
}
