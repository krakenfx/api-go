package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/pkg/derivatives"
	"github.com/krakenfx/api-go/pkg/kraken"
)

// Derivative contract.
var contract = "PI_XBTUSD"

// Direction of the order.
var side = "buy"

// Size of the order in quote unit.
var notionalSize = kraken.NewMoneyFromFloat64(5)

func main() {
	client := derivatives.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_FUTURES_REST_URL")
	client.PublicKey = os.Getenv("KRAKEN_API_FUTURES_PUBLIC")
	client.PrivateKey = os.Getenv("KRAKEN_API_FUTURES_SECRET")
	asset := derivatives.NewAssetManager()
	if err := asset.Use(client); err != nil {
		panic(err)
	}

	fmt.Printf("> Fetch %s market price\n", contract)
	ticker, err := client.TickerSymbol(contract)
	if err != nil {
		panic(err)
	}
	price := ticker.Ticker.Bid.Add(ticker.Ticker.Ask).Div(kraken.NewMoneyFromInt64(2))
	fmt.Printf("Mid price: %s\n", price)
	limitPrice := price.OffsetPercent(kraken.NewMoneyFromFloat64(-0.05))
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
	cliOrdID := kraken.UUID()
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
	fmt.Printf("Response: %s\n", kraken.ToJSON(response))

	fmt.Printf("> Set limit price to -10 ticks away from the original limit price\n")
	limitPrice = limitPrice.OffsetTicks(kraken.NewMoneyFromInt64(-10))
	fmt.Printf("Limit price: %s\n", limitPrice)

	fmt.Printf("> Submit edit order request\n")
	editOrderResponse, err := client.EditOrder(&derivatives.OrderRequest{
		ClientOrderID: cliOrdID,
		LimitPrice:    limitPrice.String(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %s\n", kraken.ToJSON(editOrderResponse))

	fmt.Printf("> Sending cancel order request\n")
	cancelOrderResponse, err := client.CancelOrder(&derivatives.CancelOrderRequest{
		ClientOrderID: cliOrdID,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %s\n", kraken.ToJSON(cancelOrderResponse))
}
