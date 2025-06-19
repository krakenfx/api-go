package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/pkg/decimal"
	"github.com/krakenfx/api-go/pkg/derivatives"
)

func main() {
	client := derivatives.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_FUTURES_REST_URL")
	client.PublicKey = os.Getenv("KRAKEN_API_FUTURES_PUBLIC")
	client.PrivateKey = os.Getenv("KRAKEN_API_FUTURES_SECRET")
	fmt.Printf("> Fetch open orders.\n")
	response, err := client.OpenOrders()
	if err != nil {
		panic(err)
	}
	if len(response.Result.OpenOrders) == 0 {
		fmt.Printf("No orders are open.\n")
	}
	for _, order := range response.Result.OpenOrders {
		size := decimal.NewFromInt64(0)
		if order.FilledSize != nil {
			size = size.Add(order.FilledSize)
		}
		if order.UnfilledSize != nil {
			size = size.Add(order.UnfilledSize)
		}
		var price string
		if order.LimitPrice != nil {
			price = order.LimitPrice.String()
		}
		if order.StopPrice != nil {
			price = order.StopPrice.String()
		}
		fmt.Printf("%s %s %s @ %s %s\n", order.Side, size, order.Symbol, order.OrderType, price)
	}
}
