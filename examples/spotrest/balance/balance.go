package main

import (
	"fmt"
	"os"

	"github.com/krakenfx/api-go/pkg/spot"
)

func main() {
	client := spot.NewREST()
	client.BaseURL = os.Getenv("KRAKEN_API_SPOT_REST_URL")
	client.PublicKey = os.Getenv("KRAKEN_API_SPOT_PUBLIC") + "A"
	client.PrivateKey = os.Getenv("KRAKEN_API_SPOT_SECRET")
	fmt.Printf("> Fetching spot balances.\n")
	balances, err := client.Balances()
	if err != nil {
		panic(err)
	}
	if len(balances.Result) == 0 {
		fmt.Printf("No balances found\n")
	}
	for asset, amount := range balances.Result {
		fmt.Printf("%s: %s\n", asset, amount)
	}
}
