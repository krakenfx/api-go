package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/v2/pkg/decimal"
	"github.com/krakenfx/api-go/v2/pkg/derivatives"
)

func main() {
	client := derivatives.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_FUTURES_REST_URL")
	resp, err := client.Instruments()
	if err != nil {
		panic(err)
	}
	for _, instrument := range resp.Result.Instruments {
		fmt.Printf("Symbol: %s", instrument.Symbol)
		fmt.Printf(", Lot decimals: %s", instrument.ContractValueTradePrecision)
		fmt.Printf(", Tick size: %s", instrument.TickSize)
		if instrument.MaxRelativeFundingRate != nil {
			fmt.Printf(", Max funding rate: %s%%", instrument.MaxRelativeFundingRate.Mul(decimal.NewFromInt64(100)))
		}
		fmt.Printf("\n")
	}
}
